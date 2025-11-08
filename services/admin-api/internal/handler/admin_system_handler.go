package handler

import (
	"goweb/pkg/base"
	"goweb/pkg/logger"
	"goweb/pkg/notification"
	"goweb/pkg/response"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/service"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminSystemHandler 系统管理处理器
type AdminSystemHandler struct {
	*base.BaseHandler
	systemService *service.AdminSystemService
	notifyService *notification.Service
	logger        logger.Logger
}

// NewAdminSystemHandler 创建系统管理处理器
func NewAdminSystemHandler(systemService *service.AdminSystemService, notifyService *notification.Service, log logger.Logger) *AdminSystemHandler {
	return &AdminSystemHandler{
		BaseHandler:   base.NewBaseHandler(log),
		systemService: systemService,
		notifyService: notifyService,
		logger:        log,
	}
}

// GetSystemBasicInfo 获取系统基本信息（公开接口，不需要认证）
// @Summary 获取系统基本信息
// @Description 获取系统名称、版本、描述等基本信息，用于登录页面等公开场景
// @Tags 系统管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/v1/admin/system/basic-info [get]
func (h *AdminSystemHandler) GetSystemBasicInfo(c *gin.Context) {
	info := h.systemService.GetSystemBasicInfo(c.Request.Context())
	response.Success(c, info)
}

// GetSystemInfo 获取系统信息
// @Summary 获取系统信息
// @Description 获取系统运行状态、资源使用情况等信息
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=model.AdminSystemInfo}
// @Router /api/v1/admin/system/info [get]
func (h *AdminSystemHandler) GetSystemInfo(c *gin.Context) {
	info := model.AdminSystemInfo{
		OnlineUsers: h.systemService.GetOnlineUserCount(),
		CPUUsage:    h.systemService.GetCPUUsage(),
		MemoryUsage: h.systemService.GetMemoryUsage(),
		DiskUsage:   h.systemService.GetDiskUsage(),
		NetworkIn:   h.systemService.GetNetworkIn(),
		NetworkOut:  h.systemService.GetNetworkOut(),
		Uptime:      h.systemService.GetUptime(),
		Version:     h.systemService.GetVersion(),
		Environment: h.systemService.GetEnvironment(),
	}

	response.Success(c, info)
}

// GetConfigList 获取配置列表
// @Summary 获取系统配置列表
// @Description 获取系统配置列表，支持分页和筛选
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param group query string false "配置分组"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=model.AdminSystemConfigListResponse}
// @Router /api/v1/admin/system/configs [get]
func (h *AdminSystemHandler) GetConfigList(c *gin.Context) {
	var req model.AdminSystemConfigListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	list, total, err := h.systemService.GetConfigList(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to get config list", "error", err)
		response.InternalError(c, "获取配置列表失败")
		return
	}

	response.Success(c, model.AdminSystemConfigListResponse{
		List:  list,
		Total: total,
	})
}

// GetConfig 获取单个配置
// @Summary 获取单个系统配置
// @Description 根据配置键获取系统配置
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param key path string true "配置键"
// @Success 200 {object} response.Response{data=model.AdminSystemConfig}
// @Router /api/v1/admin/system/configs/{key} [get]
func (h *AdminSystemHandler) GetConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "配置键不能为空")
		return
	}

	config, err := h.systemService.GetConfig(c.Request.Context(), key)
	if err != nil {
		h.logger.Error("Failed to get config", "error", err, "key", key)
		response.InternalError(c, "获取配置失败")
		return
	}

	if config.ID == 0 {
		response.NotFound(c, "配置不存在")
		return
	}

	response.Success(c, config)
}

// UpdateConfig 更新配置
// @Summary 更新系统配置
// @Description 更新系统配置值
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param key path string true "配置键"
// @Param request body model.UpdateSystemConfigParams true "更新请求"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/system/configs/{key} [put]
func (h *AdminSystemHandler) UpdateConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "配置键不能为空")
		return
	}

	var req struct {
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := h.systemService.UpdateConfig(c.Request.Context(), key, req.Value)
	if err != nil {
		h.logger.Error("Failed to update config", "error", err, "key", key)
		response.InternalError(c, "更新配置失败")
		return
	}

	response.Success(c, nil)
}

// TestEmailConfig 测试邮件配置
// @Summary 测试邮件配置
// @Description 测试SMTP邮件发送配置是否正确
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email body string true "测试邮箱地址"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/system/email/test [post]
func (h *AdminSystemHandler) TestEmailConfig(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取当前邮件配置
	smtpHost, err := h.systemService.GetConfigValue(c.Request.Context(), "email.smtp_host")
	if err != nil {
		response.InternalError(c, "获取邮件配置失败")
		return
	}

	// 如果未配置SMTP，返回错误
	if smtpHost == "" {
		response.BadRequest(c, "邮件服务未配置")
		return
	}

	// 发送测试邮件
	err = h.systemService.SendTestEmail(c.Request.Context(), req.Email)
	if err != nil {
		h.logger.Error("Failed to send test email", "error", err, "email", req.Email)
		response.InternalError(c, "发送测试邮件失败: "+err.Error())
		return
	}

	response.Success(c, "测试邮件发送成功")
}

// TestCacheConnection 测试缓存连接
// @Summary 测试缓存连接
// @Description 测试Redis缓存连接是否正常
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/admin/system/cache/test [post]
func (h *AdminSystemHandler) TestCacheConnection(c *gin.Context) {
	err := h.systemService.TestCacheConnection(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to test cache connection", "error", err)
		response.InternalError(c, "缓存连接测试失败: "+err.Error())
		return
	}

	response.Success(c, "缓存连接正常")
}

// ClearCache 清空缓存
// @Summary 清空缓存
// @Description 清空系统缓存
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/admin/system/cache/clear [post]
func (h *AdminSystemHandler) ClearCache(c *gin.Context) {
	err := h.systemService.ClearCache(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to clear cache", "error", err)
		response.InternalError(c, "清空缓存失败: "+err.Error())
		return
	}

	response.Success(c, "缓存已清空")
}

// GetLogList 获取日志列表
// @Summary 获取系统日志列表
// @Description 获取系统操作日志列表，支持分页和筛选
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param user_id query int false "用户ID"
// @Param username query string false "用户名"
// @Param method query string false "请求方法"
// @Param path query string false "请求路径"
// @Param ip query string false "IP地址"
// @Param status_code query int false "状态码"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.Response{data=model.AdminOperationLogListResponse}
// @Router /api/v1/admin/system/logs [get]
func (h *AdminSystemHandler) GetLogList(c *gin.Context) {
	var req model.AdminOperationLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	list, total, err := h.systemService.GetLogList(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to get log list", "error", err)
		response.InternalError(c, "获取日志列表失败")
		return
	}

	response.Success(c, model.AdminOperationLogListResponse{
		List:  list,
		Total: total,
	})
}

// ClearLogs 清空日志
// @Summary 清空系统日志
// @Description 清空系统操作日志
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/admin/system/logs/clear [post]
func (h *AdminSystemHandler) ClearLogs(c *gin.Context) {
	err := h.systemService.ClearLogs(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to clear logs", "error", err)
		response.InternalError(c, "清空日志失败")
		return
	}

	response.Success(c, "日志已清空")
}

// GetRecentLoginUsers 获取最近登录的用户
// @Summary 获取最近登录的用户
// @Description 获取最近成功登录的用户列表，用于仪表盘展示
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} response.Response{data=[]model.RecentLoginUser}
// @Router /api/v1/admin/system/recent-login-users [get]
func (h *AdminSystemHandler) GetRecentLoginUsers(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	users, err := h.systemService.GetRecentLoginUsers(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get recent login users", "error", err)
		response.InternalError(c, "获取最近登录用户失败")
		return
	}

	response.Success(c, users)
}

// GetRuntimeInfo 获取运行时信息
// @Summary 获取系统运行时信息
// @Description 获取Go运行时信息，包括内存使用、协程数等
// @Tags 系统管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/admin/system/runtime [get]
func (h *AdminSystemHandler) GetRuntimeInfo(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	info := map[string]interface{}{
		"goroutines":     runtime.NumGoroutine(),
		"memory_alloc":   memStats.Alloc,
		"memory_sys":     memStats.Sys,
		"memory_heap":    memStats.HeapAlloc,
		"memory_stack":   memStats.StackInuse,
		"gc_cycles":      memStats.NumGC,
		"gc_pause_total": memStats.PauseTotalNs,
		"num_cpu":        runtime.NumCPU(),
		"go_version":     runtime.Version(),
		"os":             runtime.GOOS,
		"arch":           runtime.GOARCH,
	}

	response.Success(c, info)
}

// HealthCheck 健康检查
// @Summary 系统健康检查
// @Description 检查系统各组件健康状态
// @Tags 系统管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/admin/system/health [get]
func (h *AdminSystemHandler) HealthCheck(c *gin.Context) {
	health := h.systemService.CheckHealth(c.Request.Context())

	if health["status"] == "healthy" {
		c.JSON(http.StatusOK, health)
	} else {
		c.JSON(http.StatusServiceUnavailable, health)
	}
}
