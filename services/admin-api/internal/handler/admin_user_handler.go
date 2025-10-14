package handler

import (
	"strconv"

	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminUserHandler 管理后台用户处理器
type AdminUserHandler struct {
	userService *service.UserService
	logger      logger.Logger
}

// NewAdminUserHandler 创建管理后台用户处理器实例
func NewAdminUserHandler(userService *service.UserService) *AdminUserHandler {
	return &AdminUserHandler{
		userService: userService,
	}
}

// SetLogger 设置日志器
func (h *AdminUserHandler) SetLogger(logger logger.Logger) {
	h.logger = logger
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新的管理员用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.AdminUserCreateRequest true "创建用户请求"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/users [post]
func (h *AdminUserHandler) CreateUser(c *gin.Context) {
	var req model.AdminUserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind create user request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层创建用户
	if err := h.userService.CreateUser(&req); err != nil {
		h.logger.Error("create user error", err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "创建成功"})
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新管理员用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body model.AdminUserUpdateRequest true "更新用户请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/users/{id} [put]
func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req model.AdminUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind update user request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层更新用户
	if err := h.userService.UpdateUser(id, &req); err != nil {
		h.logger.Error("update user error", err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 获取管理员用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param status query int false "用户状态"
// @Param role_id query int false "角色ID"
// @Success 200 {object} response.Response{data=model.AdminUserListResponse} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/users [get]
func (h *AdminUserHandler) GetUsers(c *gin.Context) {
	var req model.AdminUserListRequest
	// 使用ShouldBindQuery，如果参数不存在不会报错
	_ = c.ShouldBindQuery(&req)

	// 设置默认值（更科学的做法：参数缺失时提供默认值）
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用服务层获取用户列表
	result, err := h.userService.GetUsers(&req)
	if err != nil {
		h.logger.Error("get users error", err)
		response.InternalError(c, "获取用户列表失败")
		return
	}

	response.Success(c, result)
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 获取指定管理员用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=model.AdminUser} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/users/{id} [get]
func (h *AdminUserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	// 调用服务层获取用户详情
	user, err := h.userService.GetUser(id)
	if err != nil {
		h.logger.Error("get user error", err)
		response.InternalError(c, "获取用户详情失败")
		return
	}

	response.Success(c, user)
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新管理员用户的状态
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param status query int true "用户状态" Enums(0, 1)
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/users/{id}/status [put]
func (h *AdminUserHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	statusStr := c.Query("status")
	status, err := strconv.ParseInt(statusStr, 10, 8)
	if err != nil {
		response.BadRequest(c, "无效的状态值")
		return
	}

	// 调用服务层更新用户状态
	if err := h.userService.UpdateUserStatus(id, int8(status)); err != nil {
		h.logger.Error("update user status error", err)
		response.InternalError(c, "更新用户状态失败")
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定的管理员用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/users/{id} [delete]
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	// 调用服务层删除用户
	if err := h.userService.DeleteUser(id); err != nil {
		h.logger.Error("delete user error", err)
		response.InternalError(c, "删除用户失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
