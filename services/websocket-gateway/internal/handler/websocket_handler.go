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
	if userName == nil {
		userName = userID
	}

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
func (h *WebSocketHandler) GetStats(c *gin.Context) {
	stats := h.manager.GetStats()
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// GetOnlineUsers 获取在线用户列表
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

