package router

import (
	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/middleware"
	pkgRedis "goweb/pkg/redis"
	"goweb/pkg/websocket"
	"goweb/pkg/websocket/group"
	"goweb/pkg/websocket/session"
	"goweb/services/websocket-gateway/internal/handler"
	"goweb/services/websocket-gateway/internal/service"
)

func NewRouter(cfg *config.Config, log logger.Logger, wsManager *websocket.Manager, sessionMgr *session.SessionManager, groupMgr *group.GroupManager, redisClient *pkgRedis.Client) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestID())
	r.Use(middleware.AccessLogger(log))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-Id"},
		ExposeHeaders:    []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Redis PubSub 服务（如果 Redis 可用）
	if redisClient != nil {
		pubsubService := service.NewPubSubService(wsManager, redisClient, log)
		ctx := context.Background()
		if err := pubsubService.Start(ctx); err != nil {
			log.Error("failed to start pubsub service", err)
		} else {
			log.Info("Redis PubSub service started")
		}
	} else {
		log.Warn("Redis is disabled, WebSocket pubsub will not work")
	}

	// 创建 WebSocket 处理器
	wsHandler := handler.NewWebSocketHandler(wsManager, log)

	// 创建会话处理器
	sessionHandler := handler.NewSessionHandler(wsManager, sessionMgr, log)

	// 创建分组处理器
	groupHandler := handler.NewGroupHandler(wsManager, groupMgr, log)

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "GinForge WebSocket Gateway",
			"version": "0.1.0",
			"status":  "running",
			"endpoints": gin.H{
				"health":       "/healthz",
				"ws":           "/ws (需要 JWT 认证)",
				"stats":        "/ws/stats",
				"online_users": "/ws/online-users (需要 JWT 认证)",
			},
		})
	})

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "websocket-gateway",
		})
	})

	// WebSocket 路由
	wsGroup := r.Group("/ws")
	{
		// 需要 JWT 认证的路由
		wsAuth := wsGroup.Group("")
		wsAuth.Use(middleware.WebSocketAuth(cfg.GetString("jwt.secret")))
		{
			// WebSocket 连接端点
			wsAuth.GET("", wsHandler.HandleConnection)

			// 在线用户列表
			wsAuth.GET("/online-users", wsHandler.GetOnlineUsers)
		}

		// 无需认证的路由（监控、健康检查）
		wsGroup.GET("/stats", wsHandler.GetStats)

		// 会话管理路由
		sessionGroup := wsGroup.Group("/session")
		sessionGroup.Use(middleware.JWTAuth(cfg.GetString("jwt.secret")))
		{
			sessionGroup.GET("/client/:client_id", sessionHandler.GetSessionData)
			sessionGroup.POST("/client/:client_id", sessionHandler.SetSessionData)
			sessionGroup.DELETE("/client/:client_id/:key", sessionHandler.DeleteSessionData)
			sessionGroup.GET("/user/:user_id", sessionHandler.GetUserSessions)
		}

		// 分组管理路由
		groupRoutes := wsGroup.Group("/group")
		groupRoutes.Use(middleware.JWTAuth(cfg.GetString("jwt.secret")))
		{
			groupRoutes.GET("/:group_id/members", groupHandler.GetGroupMembers)
			groupRoutes.POST("/:group_id/join/:client_id", groupHandler.JoinGroup)
			groupRoutes.DELETE("/:group_id/leave/:client_id", groupHandler.LeaveGroup)
			groupRoutes.GET("/client/:client_id", groupHandler.GetClientGroups)
			groupRoutes.GET("/:group_id/metadata", groupHandler.GetGroupMetadata)
			groupRoutes.POST("/:group_id/metadata", groupHandler.SetGroupMetadata)
			groupRoutes.GET("", groupHandler.GetAllGroups)
		}

		// 消息发送路由
		msgGroup := wsGroup.Group("/message")
		msgGroup.Use(middleware.JWTAuth(cfg.GetString("jwt.secret")))
		{
			// 发送通知给用户
			msgGroup.POST("/notification/user/:user_id", func(c *gin.Context) {
				userID := c.Param("user_id")

				var req struct {
					Title string `json:"title" binding:"required"`
					Body  string `json:"body" binding:"required"`
					Icon  string `json:"icon"`
					Link  string `json:"link"`
				}

				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}

				// 创建通知消息
				notification := websocket.NewNotificationMessage(req.Title, req.Body)
				if req.Icon != "" {
					notification.Icon = req.Icon
				}
				if req.Link != "" {
					notification.Link = req.Link
				}

				// 发送通知
				msg := websocket.NewMessage(websocket.MessageTypeNotification, notification)
				if err := wsManager.SendToUser(userID, msg); err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, gin.H{"success": true})
			})

			// 广播系统消息
			msgGroup.POST("/system/broadcast", func(c *gin.Context) {
				var req struct {
					Message string `json:"message" binding:"required"`
					Level   string `json:"level"`
					Code    int    `json:"code"`
				}

				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}

				// 创建系统消息
				sysMsg := &websocket.SystemMessage{
					Message: req.Message,
					Level:   req.Level,
					Code:    req.Code,
				}

				// 广播系统消息
				msg := websocket.NewMessage(websocket.MessageTypeSystem, sysMsg)
				wsManager.Broadcast(msg)

				c.JSON(200, gin.H{"success": true})
			})

			// 发送数据更新消息到分组
			msgGroup.POST("/data-update/group/:group_id", func(c *gin.Context) {
				groupID := c.Param("group_id")

				var req struct {
					Entity string                 `json:"entity" binding:"required"`
					Action string                 `json:"action" binding:"required"`
					ID     string                 `json:"id"`
					Data   map[string]interface{} `json:"data"`
				}

				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}

				// 创建数据更新消息
				updateMsg := &websocket.DataUpdateMessage{
					Entity: req.Entity,
					Action: req.Action,
					ID:     req.ID,
					Data:   req.Data,
				}

				// 发送数据更新消息
				msg := websocket.NewMessage(websocket.MessageTypeDataUpdate, updateMsg)
				if err := wsManager.BroadcastToRoom(groupID, msg); err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}

				c.JSON(200, gin.H{"success": true})
			})
		}
	}

	return r
}
