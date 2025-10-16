package handler

import (
	"strconv"
	
	"github.com/gin-gonic/gin"
	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/service"
)

// ArticlesHandler Articles管理 Handler
type ArticlesHandler struct {
	service *service.ArticlesService
	logger  logger.Logger
}

// NewArticlesHandler 创建 Handler 实例
func NewArticlesHandler(service *service.ArticlesService, logger logger.Logger) *ArticlesHandler {
	return &ArticlesHandler{
		service: service,
		logger:  logger,
	}
}

// List 获取Articles管理列表
// @Summary 获取Articles管理列表
// @Description 获取Articles管理列表（支持分页、搜索、排序）
// @Tags Articles管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param sort_by query string false "排序字段"
// @Param sort_order query string false "排序方式(asc/desc)"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/articleses [get]
func (h *ArticlesHandler) List(c *gin.Context) {
	var req model.ArticlesListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	
	list, total, err := h.service.List(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	
	response.Success(c, gin.H{
		"list":  model.ToArticlesResponseList(list),
		"total": total,
		"page":  req.Page,
		"page_size": req.PageSize,
	})
}

// Get 获取Articles管理详情
// @Summary 获取Articles管理详情
// @Description 根据 ID 获取Articles管理详情
// @Tags Articles管理
// @Accept json
// @Produce json
// @Param id path int true "Articles管理 ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/articleses/{id} [get]
func (h *ArticlesHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "ID 格式错误")
		return
	}
	
	articles, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, 404, err.Error())
		return
	}
	
	response.Success(c, model.ToArticlesResponse(articles))
}

// Create 创建Articles管理
// @Summary 创建Articles管理
// @Description 创建新的Articles管理
// @Tags Articles管理
// @Accept json
// @Produce json
// @Param body body model.ArticlesCreateRequest true "Articles管理信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/articleses [post]
func (h *ArticlesHandler) Create(c *gin.Context) {
	var req model.ArticlesCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	
	articles, err := h.service.Create(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	
	response.Success(c, model.ToArticlesResponse(articles))
}

// Update 更新Articles管理
// @Summary 更新Articles管理
// @Description 更新Articles管理信息
// @Tags Articles管理
// @Accept json
// @Produce json
// @Param id path int true "Articles管理 ID"
// @Param body body model.ArticlesUpdateRequest true "Articles管理信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/articleses/{id} [put]
func (h *ArticlesHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "ID 格式错误")
		return
	}
	
	var req model.ArticlesUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	
	if err := h.service.Update(id, &req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	
	response.Success(c, nil)
}

// Delete 删除Articles管理
// @Summary 删除Articles管理
// @Description 删除Articles管理
// @Tags Articles管理
// @Accept json
// @Produce json
// @Param id path int true "Articles管理 ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/articleses/{id} [delete]
func (h *ArticlesHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "ID 格式错误")
		return
	}
	
	if err := h.service.Delete(id); err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	
	response.Success(c, nil)
}
