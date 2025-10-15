package handler

import (
	"strconv"

	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminPermissionHandler 管理后台权限处理器
type AdminPermissionHandler struct {
	permissionService *service.PermissionService
	logger            logger.Logger
}

// NewAdminPermissionHandler 创建管理后台权限处理器实例
func NewAdminPermissionHandler(permissionService *service.PermissionService) *AdminPermissionHandler {
	return &AdminPermissionHandler{
		permissionService: permissionService,
	}
}

// SetLogger 设置日志器
func (h *AdminPermissionHandler) SetLogger(logger logger.Logger) {
	h.logger = logger
}

// CreatePermission 创建权限
// @Summary 创建权限
// @Description 创建新的管理员权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.AdminPermissionCreateRequest true "创建权限请求"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/permissions [post]
func (h *AdminPermissionHandler) CreatePermission(c *gin.Context) {
	var req model.AdminPermissionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind create permission request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层创建权限
	if err := h.permissionService.CreatePermission(&req); err != nil {
		h.logger.Error("create permission error", err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "创建成功"})
}

// UpdatePermission 更新权限
// @Summary 更新权限
// @Description 更新管理员权限信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "权限ID"
// @Param request body model.AdminPermissionUpdateRequest true "更新权限请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/permissions/{id} [put]
func (h *AdminPermissionHandler) UpdatePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	var req model.AdminPermissionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind update permission request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层更新权限
	if err := h.permissionService.UpdatePermission(id, &req); err != nil {
		h.logger.Error("update permission error", err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// GetPermissions 获取权限列表
// @Summary 获取权限列表
// @Description 获取管理员权限列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param type query string false "权限类型"
// @Success 200 {object} response.Response{data=model.AdminPermissionListResponse} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/permissions [get]
func (h *AdminPermissionHandler) GetPermissions(c *gin.Context) {
	var req model.AdminPermissionListRequest
	// 使用ShouldBindQuery，如果参数不存在不会报错
	_ = c.ShouldBindQuery(&req)

	// 设置默认值（更科学的做法：参数缺失时提供默认值）
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用服务层获取权限列表
	result, err := h.permissionService.GetPermissions(&req)
	if err != nil {
		h.logger.Error("get permissions error", err)
		response.InternalError(c, "获取权限列表失败")
		return
	}

	response.Success(c, result)
}

// GetPermission 获取权限详情
// @Summary 获取权限详情
// @Description 获取指定管理员权限的详细信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "权限ID"
// @Success 200 {object} response.Response{data=model.AdminPermission} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "权限不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/permissions/{id} [get]
func (h *AdminPermissionHandler) GetPermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	// 调用服务层获取权限详情
	permission, err := h.permissionService.GetPermission(id)
	if err != nil {
		h.logger.Error("get permission error", err)
		response.InternalError(c, "获取权限详情失败")
		return
	}

	response.Success(c, permission)
}

// DeletePermission 删除权限
// @Summary 删除权限
// @Description 删除指定的管理员权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "权限ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/permissions/{id} [delete]
func (h *AdminPermissionHandler) DeletePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	// 调用服务层删除权限
	if err := h.permissionService.DeletePermission(id); err != nil {
		h.logger.Error("delete permission error", err)
		response.InternalError(c, "删除权限失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// UpdatePermissionStatus 更新权限状态
// @Summary 更新权限状态
// @Description 更新权限的启用/禁用状态
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "权限ID"
// @Param status query int true "状态：0=禁用，1=启用"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/permissions/{id}/status [put]
func (h *AdminPermissionHandler) UpdatePermissionStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	statusStr := c.Query("status")
	status, err := strconv.ParseInt(statusStr, 10, 8)
	if err != nil || (status != 0 && status != 1) {
		response.BadRequest(c, "无效的状态值，必须是0或1")
		return
	}

	// 调用服务层更新权限状态
	if err := h.permissionService.UpdatePermissionStatus(id, int8(status)); err != nil {
		h.logger.Error("update permission status error", err)
		response.InternalError(c, "更新权限状态失败")
		return
	}

	response.Success(c, gin.H{"message": "状态更新成功"})
}
