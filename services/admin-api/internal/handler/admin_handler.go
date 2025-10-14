package handler

import (
	"github.com/gin-gonic/gin"

	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/admin-api/internal/service"
)

// AdminHandler 管理后台处理器 (已废弃，使用新的处理器)
type AdminHandler struct {
	userService *service.UserService
	logger      logger.Logger
}

func NewAdminHandler(userService *service.UserService) *AdminHandler {
	return &AdminHandler{
		userService: userService,
	}
}

// SetLogger 设置日志器
func (h *AdminHandler) SetLogger(logger logger.Logger) {
	h.logger = logger
}

// GetUsers 获取用户列表 (已废弃)
func (h *AdminHandler) GetUsers(c *gin.Context) {
	response.Success(c, gin.H{
		"list":  []interface{}{},
		"total": 0,
	})
}

// UpdateUserStatus 更新用户状态 (已废弃)
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	response.Success(c, gin.H{"message": "更新成功"})
}
