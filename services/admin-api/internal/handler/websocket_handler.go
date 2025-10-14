package handler

import (
	"goweb/pkg/logger"
	"goweb/pkg/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WebSocketHandler WebSocket 处理器
type WebSocketHandler struct {
	manager *websocket.Manager
	logger  logger.Logger
}

// NewWebSocketHandler 创建 WebSocket 处理器
func NewWebSocketHandler(manager *websocket.Manager, log logger.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		manager: manager,
		logger:  log,
	}
}

// HandleConnection 处理 WebSocket 连接
// @Summary WebSocket 连接
// @Description 建立 WebSocket 实时连接
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param token query string true "JWT Token"
// @Success 101 "Switching Protocols"
// @Router /ws [get]
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	// 从上下文获取用户信息（已由中间件验证）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}
	
	userName, _ := c.Get("username")
	
	// 升级 HTTP 连接为 WebSocket
	conn, err := websocket.DefaultUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("failed to upgrade websocket", err,
			"user_id", userID,
			"remote_addr", c.ClientIP())
		return
	}
	
	// 创建客户端
	clientID := uuid.New().String()
	client := websocket.NewClient(
		clientID,
		userID.(string),
		userName.(string),
		conn,
		h.manager,
		h.logger,
	)
	
	// 注册客户端
	h.manager.Register(client)
	
	h.logger.Info("websocket connection established",
		"client_id", clientID,
		"user_id", userID,
		"user_name", userName,
		"remote_addr", c.ClientIP())
	
	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()
}

// GetStats 获取 WebSocket 统计信息
// @Summary WebSocket 统计
// @Description 获取 WebSocket 连接统计信息
// @Tags WebSocket
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /ws/stats [get]
// @Security BearerAuth
func (h *WebSocketHandler) GetStats(c *gin.Context) {
	stats := h.manager.GetStats()
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// BroadcastNotification 广播通知给所有用户
// @Summary 广播通知
// @Description 向所有在线用户广播通知
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param request body NotificationRequest true "通知内容"
// @Success 200 {object} response.Response
// @Router /ws/broadcast [post]
// @Security BearerAuth
func (h *WebSocketHandler) BroadcastNotification(c *gin.Context) {
	var req NotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	
	notification := &websocket.NotificationMessage{
		Title: req.Title,
		Body:  req.Body,
		Icon:  req.Icon,
		Link:  req.Link,
	}
	
	msg := websocket.NewMessage(websocket.MessageTypeNotification, notification)
	h.manager.Broadcast(msg)
	
	h.logger.Info("notification broadcasted",
		"title", req.Title,
		"operator", c.GetString("username"))
	
	c.JSON(200, gin.H{
		"code":    0,
		"message": "广播成功",
	})
}

// SendNotificationToUser 发送通知给指定用户
// @Summary 发送通知给用户
// @Description 向指定用户发送通知
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param user_id path string true "用户ID"
// @Param request body NotificationRequest true "通知内容"
// @Success 200 {object} response.Response
// @Router /ws/users/{user_id}/notification [post]
// @Security BearerAuth
func (h *WebSocketHandler) SendNotificationToUser(c *gin.Context) {
	userID := c.Param("user_id")
	
	var req NotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	
	notification := &websocket.NotificationMessage{
		Title: req.Title,
		Body:  req.Body,
		Icon:  req.Icon,
		Link:  req.Link,
	}
	
	err := h.manager.SendNotification(userID, notification)
	if err != nil {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "用户不在线",
		})
		return
	}
	
	h.logger.Info("notification sent to user",
		"user_id", userID,
		"title", req.Title,
		"operator", c.GetString("username"))
	
	c.JSON(200, gin.H{
		"code":    0,
		"message": "发送成功",
	})
}

// GetOnlineUsers 获取在线用户列表
// @Summary 在线用户列表
// @Description 获取当前在线的用户列表
// @Tags WebSocket
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /ws/online-users [get]
// @Security BearerAuth
func (h *WebSocketHandler) GetOnlineUsers(c *gin.Context) {
	users := h.manager.GetOnlineUsers()
	
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total": len(users),
			"users": users,
		},
	})
}

// NotificationRequest 通知请求
type NotificationRequest struct {
	Title string `json:"title" binding:"required"` // 标题
	Body  string `json:"body" binding:"required"`  // 内容
	Icon  string `json:"icon"`                     // 图标
	Link  string `json:"link"`                     // 链接
}

