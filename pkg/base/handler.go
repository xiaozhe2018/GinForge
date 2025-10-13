package base

import (
	"goweb/pkg/logger"
	"goweb/pkg/response"

	"github.com/gin-gonic/gin"
)

// BaseHandler 基础处理器类
type BaseHandler struct {
	logger logger.Logger
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler(logger logger.Logger) *BaseHandler {
	return &BaseHandler{
		logger: logger,
	}
}

// SetLogger 设置日志器
func (h *BaseHandler) SetLogger(logger logger.Logger) {
	h.logger = logger
}

// GetLogger 获取日志器
func (h *BaseHandler) GetLogger() logger.Logger {
	return h.logger
}

// LogInfo 记录信息日志
func (h *BaseHandler) LogInfo(msg string, fields ...interface{}) {
	if h.logger != nil {
		h.logger.Info(msg, fields...)
	}
}

// LogError 记录错误日志
func (h *BaseHandler) LogError(msg string, err error, fields ...interface{}) {
	if h.logger != nil {
		allFields := append([]interface{}{"error", err}, fields...)
		h.logger.Error(msg, allFields...)
	}
}

// LogWarn 记录警告日志
func (h *BaseHandler) LogWarn(msg string, fields ...interface{}) {
	if h.logger != nil {
		h.logger.Warn(msg, fields...)
	}
}

// LogDebug 记录调试日志
func (h *BaseHandler) LogDebug(msg string, fields ...interface{}) {
	if h.logger != nil {
		h.logger.Debug(msg, fields...)
	}
}

// Success 成功响应
func (h *BaseHandler) Success(c *gin.Context, data interface{}) {
	response.Success(c, data)
}

// Error 错误响应
func (h *BaseHandler) Error(c *gin.Context, code int, message string) {
	response.Error(c, code, message)
}

// BadRequest 400错误
func (h *BaseHandler) BadRequest(c *gin.Context, message string) {
	response.BadRequest(c, message)
}

// Unauthorized 401错误
func (h *BaseHandler) Unauthorized(c *gin.Context, message string) {
	response.Unauthorized(c, message)
}

// Forbidden 403错误
func (h *BaseHandler) Forbidden(c *gin.Context, message string) {
	response.Forbidden(c, message)
}

// NotFound 404错误
func (h *BaseHandler) NotFound(c *gin.Context, message string) {
	response.NotFound(c, message)
}

// InternalError 500错误
func (h *BaseHandler) InternalError(c *gin.Context, message string) {
	response.InternalError(c, message)
}

// GetUserID 获取用户ID
func (h *BaseHandler) GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(string)
	}
	return ""
}

// GetUsername 获取用户名
func (h *BaseHandler) GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// GetTraceID 获取追踪ID
func (h *BaseHandler) GetTraceID(c *gin.Context) string {
	return c.GetString("request_id")
}
