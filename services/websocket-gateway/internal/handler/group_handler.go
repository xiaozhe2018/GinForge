package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"goweb/pkg/logger"
	"goweb/pkg/websocket"
	"goweb/pkg/websocket/group"
)

// GroupHandler 分组处理器
type GroupHandler struct {
	manager  *websocket.Manager
	groupMgr *group.GroupManager
	logger   logger.Logger
}

// NewGroupHandler 创建分组处理器
func NewGroupHandler(manager *websocket.Manager, groupMgr *group.GroupManager, logger logger.Logger) *GroupHandler {
	return &GroupHandler{
		manager:  manager,
		groupMgr: groupMgr,
		logger:   logger,
	}
}

// GetGroupMembers 获取分组成员
func (h *GroupHandler) GetGroupMembers(c *gin.Context) {
	groupID := c.Param("group_id")
	
	// 获取分组成员
	members, err := h.groupMgr.GetGroupMembers(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get group members",
			"error":   err.Error(),
		})
		return
	}
	
	// 返回分组成员
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"group_id": groupID,
		"members":  members,
		"count":    len(members),
	})
}

// JoinGroup 加入分组
func (h *GroupHandler) JoinGroup(c *gin.Context) {
	groupID := c.Param("group_id")
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
	
	// 客户端信息
	clientInfo := map[string]interface{}{
		"user_id":   client.UserID,
		"user_name": client.UserName,
	}
	
	// 加入分组
	if err := h.groupMgr.JoinGroup(groupID, clientID, clientInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to join group",
			"error":   err.Error(),
		})
		return
	}
	
	// 返回成功
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "joined group successfully",
	})
}

// LeaveGroup 离开分组
func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	groupID := c.Param("group_id")
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
	
	// 离开分组
	if err := h.groupMgr.LeaveGroup(groupID, clientID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to leave group",
			"error":   err.Error(),
		})
		return
	}
	
	// 返回成功
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "left group successfully",
	})
}

// GetClientGroups 获取客户端分组
func (h *GroupHandler) GetClientGroups(c *gin.Context) {
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
	
	// 获取客户端分组
	groups, err := h.groupMgr.GetClientGroups(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get client groups",
			"error":   err.Error(),
		})
		return
	}
	
	// 返回客户端分组
	c.JSON(http.StatusOK, gin.H{
		"code":     0,
		"message":  "success",
		"client_id": clientID,
		"groups":    groups,
		"count":     len(groups),
	})
}

// SetGroupMetadata 设置分组元数据
func (h *GroupHandler) SetGroupMetadata(c *gin.Context) {
	groupID := c.Param("group_id")
	
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
	
	// 设置分组元数据
	if err := h.groupMgr.SetGroupMetadata(groupID, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to set group metadata",
			"error":   err.Error(),
		})
		return
	}
	
	// 返回成功
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "group metadata set successfully",
	})
}

// GetGroupMetadata 获取分组元数据
func (h *GroupHandler) GetGroupMetadata(c *gin.Context) {
	groupID := c.Param("group_id")
	
	// 获取分组元数据
	metadata, err := h.groupMgr.GetGroupMetadata(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get group metadata",
			"error":   err.Error(),
		})
		return
	}
	
	// 返回分组元数据
	c.JSON(http.StatusOK, gin.H{
		"code":     0,
		"message":  "success",
		"group_id": groupID,
		"metadata": metadata,
	})
}

// GetAllGroups 获取所有分组
func (h *GroupHandler) GetAllGroups(c *gin.Context) {
	// 获取所有分组
	groups, err := h.groupMgr.GetAllGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "failed to get all groups",
			"error":   err.Error(),
		})
		return
	}
	
	// 返回所有分组
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"groups":  groups,
		"count":   len(groups),
	})
}
