package handler

import (
	"strconv"

	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminMenuHandler 管理后台菜单处理器
type AdminMenuHandler struct {
	menuService *service.MenuService
	logger      logger.Logger
}

// NewAdminMenuHandler 创建管理后台菜单处理器实例
func NewAdminMenuHandler(menuService *service.MenuService) *AdminMenuHandler {
	return &AdminMenuHandler{
		menuService: menuService,
	}
}

// SetLogger 设置日志器
func (h *AdminMenuHandler) SetLogger(logger logger.Logger) {
	h.logger = logger
}

// CreateMenu 创建菜单
// @Summary 创建菜单
// @Description 创建新的管理员菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.AdminMenuCreateRequest true "创建菜单请求"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/menus [post]
func (h *AdminMenuHandler) CreateMenu(c *gin.Context) {
	var req model.AdminMenuCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind create menu request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层创建菜单
	if err := h.menuService.CreateMenu(&req); err != nil {
		h.logger.Error("create menu error", err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "创建成功"})
}

// UpdateMenu 更新菜单
// @Summary 更新菜单
// @Description 更新管理员菜单信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Param request body model.AdminMenuUpdateRequest true "更新菜单请求"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/menus/{id} [put]
func (h *AdminMenuHandler) UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	var req model.AdminMenuUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("bind update menu request error", err)
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 调用服务层更新菜单
	if err := h.menuService.UpdateMenu(id, &req); err != nil {
		h.logger.Error("update menu error", err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "更新成功"})
}

// GetMenus 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取管理员菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param type query string false "菜单类型"
// @Param status query int false "菜单状态"
// @Param parent_id query int false "父菜单ID"
// @Success 200 {object} response.Response{data=model.AdminMenuListResponse} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/menus [get]
func (h *AdminMenuHandler) GetMenus(c *gin.Context) {
	var req model.AdminMenuListRequest
	// 使用ShouldBindQuery，如果参数不存在不会报错
	_ = c.ShouldBindQuery(&req)

	// 设置默认值（更科学的做法：参数缺失时提供默认值）
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用服务层获取菜单列表
	result, err := h.menuService.GetMenus(&req)
	if err != nil {
		h.logger.Error("get menus error", err)
		response.InternalError(c, "获取菜单列表失败")
		return
	}

	response.Success(c, result)
}

// GetMenuTree 获取菜单树
// @Summary 获取菜单树
// @Description 获取管理员菜单树形结构
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.AdminMenuTreeResponse} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/menus/tree [get]
func (h *AdminMenuHandler) GetMenuTree(c *gin.Context) {
	// 调用服务层获取菜单树
	result, err := h.menuService.GetMenuTree()
	if err != nil {
		h.logger.Error("get menu tree error", err)
		response.InternalError(c, "获取菜单树失败")
		return
	}

	response.Success(c, result)
}

// GetMenu 获取菜单详情
// @Summary 获取菜单详情
// @Description 获取指定管理员菜单的详细信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response{data=model.AdminMenu} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "菜单不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/menus/{id} [get]
func (h *AdminMenuHandler) GetMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	// 调用服务层获取菜单详情
	menu, err := h.menuService.GetMenu(id)
	if err != nil {
		h.logger.Error("get menu error", err)
		response.InternalError(c, "获取菜单详情失败")
		return
	}

	response.Success(c, menu)
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Description 删除指定的管理员菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/menus/{id} [delete]
func (h *AdminMenuHandler) DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	// 调用服务层删除菜单
	if err := h.menuService.DeleteMenu(id); err != nil {
		h.logger.Error("delete menu error", err)
		response.InternalError(c, "删除菜单失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
