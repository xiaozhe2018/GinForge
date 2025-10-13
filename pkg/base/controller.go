package base

import (
	"goweb/pkg/logger"
	"goweb/pkg/model"

	"github.com/gin-gonic/gin"
)

// BaseController 基础控制器类
type BaseController struct {
	logger logger.Logger
}

// NewBaseController 创建基础控制器
func NewBaseController(logger logger.Logger) *BaseController {
	return &BaseController{
		logger: logger,
	}
}

// SetLogger 设置日志器
func (c *BaseController) SetLogger(logger logger.Logger) {
	c.logger = logger
}

// GetLogger 获取日志器
func (c *BaseController) GetLogger() logger.Logger {
	return c.logger
}

// LogInfo 记录信息日志
func (c *BaseController) LogInfo(msg string, fields ...interface{}) {
	if c.logger != nil {
		c.logger.Info(msg, fields...)
	}
}

// LogError 记录错误日志
func (c *BaseController) LogError(msg string, err error, fields ...interface{}) {
	if c.logger != nil {
		allFields := append([]interface{}{"error", err}, fields...)
		c.logger.Error(msg, allFields...)
	}
}

// LogWarn 记录警告日志
func (c *BaseController) LogWarn(msg string, fields ...interface{}) {
	if c.logger != nil {
		c.logger.Warn(msg, fields...)
	}
}

// LogDebug 记录调试日志
func (c *BaseController) LogDebug(msg string, fields ...interface{}) {
	if c.logger != nil {
		c.logger.Debug(msg, fields...)
	}
}

// GetUserID 获取用户ID
func (c *BaseController) GetUserID(ctx *gin.Context) string {
	if userID, exists := ctx.Get("user_id"); exists {
		return userID.(string)
	}
	return ""
}

// GetUsername 获取用户名
func (c *BaseController) GetUsername(ctx *gin.Context) string {
	if username, exists := ctx.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// GetTraceID 获取追踪ID
func (c *BaseController) GetTraceID(ctx *gin.Context) string {
	return ctx.GetString("request_id")
}

// GetPagination 获取分页参数
func (c *BaseController) GetPagination(ctx *gin.Context) *model.Pagination {
	var pagination model.Pagination
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		// 使用默认分页参数
		return model.NewPagination(1, 10)
	}
	return &pagination
}

// GetClientIP 获取客户端IP
func (c *BaseController) GetClientIP(ctx *gin.Context) string {
	return ctx.ClientIP()
}

// GetUserAgent 获取用户代理
func (c *BaseController) GetUserAgent(ctx *gin.Context) string {
	return ctx.GetHeader("User-Agent")
}

// GetRequestID 获取请求ID
func (c *BaseController) GetRequestID(ctx *gin.Context) string {
	return ctx.GetString("request_id")
}
