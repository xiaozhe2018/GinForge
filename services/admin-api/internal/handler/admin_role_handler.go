package handler

import (
	"strconv"

	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminRoleHandler 管理后台角色处理器
type AdminRoleHandler struct {
	roleService *service.RoleService
	logger      logger.Logger
}

// NewAdminRoleHandler 创建管理后台角色处理器实例
func NewAdminRoleHandler(roleService *service.RoleService) *AdminRoleHandler {
	return &AdminRoleHandler{
		roleService: roleService,
	}
}

// SetLogger 设置日志器
func (h *AdminRoleHandler) SetLogger(logger logger.Logger) {
	h.logger = logger
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新的管理员角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.AdminRoleCreateRequest true "创建角色请求"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/roles [post]
func (h *AdminRoleHandler) CreateRole(c *gin.Context) {
	var req model.AdminRoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind create role request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层创建角色
	if err := h.roleService.CreateRole(&req); err != nil {
		h.logger.Error("create role error", err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "创建成功"})
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新管理员角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param request body model.AdminRoleUpdateRequest true "更新角色请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/roles/{id} [put]
func (h *AdminRoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req model.AdminRoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind update role request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层更新角色
	if err := h.roleService.UpdateRole(id, &req); err != nil {
		h.logger.Error("update role error", err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// GetRoles 获取角色列表
// @Summary 获取角色列表
// @Description 获取管理员角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param status query int false "角色状态"
// @Success 200 {object} response.Response{data=model.AdminRoleListResponse} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/roles [get]
func (h *AdminRoleHandler) GetRoles(c *gin.Context) {
	var req model.AdminRoleListRequest
	// 使用ShouldBindQuery，如果参数不存在不会报错
	_ = c.ShouldBindQuery(&req)

	// 设置默认值（更科学的做法：参数缺失时提供默认值）
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用服务层获取角色列表
	result, err := h.roleService.GetRoles(&req)
	if err != nil {
		h.logger.Error("get roles error", err)
		response.InternalError(c, "获取角色列表失败")
		return
	}

	response.Success(c, result)
}

// GetRole 获取角色详情
// @Summary 获取角色详情
// @Description 获取指定管理员角色的详细信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=model.AdminRole} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "角色不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/roles/{id} [get]
func (h *AdminRoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	// 调用服务层获取角色详情
	role, err := h.roleService.GetRole(id)
	if err != nil {
		h.logger.Error("get role error", err)
		response.InternalError(c, "获取角色详情失败")
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除指定的管理员角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/roles/{id} [delete]
func (h *AdminRoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	// 调用服务层删除角色
	if err := h.roleService.DeleteRole(id); err != nil {
		h.logger.Error("delete role error", err)
		response.InternalError(c, "删除角色失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
