package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"goweb/pkg/logger"
	"goweb/pkg/websocket"
	"goweb/pkg/websocket/session"
)

// SessionHandler 会话处理器
type SessionHandler struct {
	manager      *websocket.Manager
	sessionMgr   *session.SessionManager
	logger       logger.Logger
}

// NewSessionHandler 创建会话处理器
func NewSessionHandler(manager *websocket.Manager, sessionMgr *session.SessionManager, logger logger.Logger) *SessionHandler {
	return &SessionHandler{
		manager:      manager,
		sessionMgr:   sessionMgr,
		logger:       logger,
	}
}

// GetSessionData 获取会话数据
func (h *SessionHandler) GetSessionData(c *gin.Context) {
	clientID := c.Param("client_id")
	
	// 验证客户端是否存在
	client, exists := h.manager.GetClient(clientID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "client not found",
		})
		return
	}
	
	// 获取会话
	sess, err := h.sessionMgr.GetSession(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get session",
			"error":   err.Error(),
		})
		return
	}
	
	// 返回会话数据
	c.JSON(http.StatusOK, gin.H{
		"code":     0,
		"message":  "success",
		"client_id": clientID,
		"user_id":   client.UserID,
		"data":      sess.GetAll(),
	})
}

// SetSessionData 设置会话数据
func (h *SessionHandler) SetSessionData(c *gin.Context) {
	clientID := c.Param("client_id")
	
	// 验证客户端是否存在
	_, exists := h.manager.GetClient(clientID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "client not found",
		})
		return
	}
	
	// 解析请求数据
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "invalid request data",
			"error":   err.Error(),
		})
		return
	}
	
	// 更新会话数据
	if err := h.sessionMgr.UpdateSessionData(clientID, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to update session data",
			"error":   err.Error(),
		})
		return
	}
	
	// 返回成功
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "session data updated successfully",
	})
}

// DeleteSessionData 删除会话数据
func (h *SessionHandler) DeleteSessionData(c *gin.Context) {
	clientID := c.Param("client_id")
	key := c.Param("key")
	
	// 验证客户端是否存在
	_, exists := h.manager.GetClient(clientID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "client not found",
		})
		return
	}
	
	// 删除会话数据
	if err := h.sessionMgr.SetSessionData(clientID, key, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to delete session data",
			"error":   err.Error(),
		})
		return
	}
	
	// 返回成功
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "session data deleted successfully",
	})
}

// GetUserSessions 获取用户所有会话
func (h *SessionHandler) GetUserSessions(c *gin.Context) {
	userID := c.Param("user_id")
	
	// 获取用户会话
	sessions, err := h.sessionMgr.GetUserSessions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get user sessions",
			"error":   err.Error(),
		})
		return
	}
	
	// 构建响应数据
	result := make([]map[string]interface{}, 0, len(sessions))
	for _, sess := range sessions {
		result = append(result, map[string]interface{}{
			"session_id": sess.ID,
			"user_id":    sess.UserID,
			"created_at": sess.CreatedAt,
			"updated_at": sess.UpdatedAt,
			"data":       sess.GetAll(),
		})
	}
	
	// 返回会话数据
	c.JSON(http.StatusOK, gin.H{
		"code":     0,
		"message":  "success",
		"user_id":  userID,
		"sessions": result,
		"count":    len(result),
	})
}