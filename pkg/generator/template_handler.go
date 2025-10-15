package generator

import (
	"bytes"
	"text/template"
)

// renderHandlerTemplate 渲染 Handler 模板
func renderHandlerTemplate(data *TemplateData) (string, error) {
	tmpl := `package handler

import (
	"strconv"
	
	"github.com/gin-gonic/gin"
	"goweb/pkg/logger"
	"goweb/pkg/response"
	"goweb/services/{{.Module}}-api/internal/model"
	"goweb/services/{{.Module}}-api/internal/service"
)

// {{.ModelName}}Handler {{.Title}} Handler
type {{.ModelName}}Handler struct {
	service *service.{{.ModelName}}Service
	logger  logger.Logger
}

// New{{.ModelName}}Handler 创建 Handler 实例
func New{{.ModelName}}Handler(service *service.{{.ModelName}}Service, logger logger.Logger) *{{.ModelName}}Handler {
	return &{{.ModelName}}Handler{
		service: service,
		logger:  logger,
	}
}

// List 获取{{.Title}}列表
// @Summary 获取{{.Title}}列表
// @Description 获取{{.Title}}列表（支持分页、搜索、排序）
// @Tags {{.Title}}
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
{{- if .HasSearch}}
// @Param keyword query string false "搜索关键词"
{{- end}}
{{- if .HasSort}}
// @Param sort_by query string false "排序字段"
// @Param sort_order query string false "排序方式(asc/desc)"
{{- end}}
// @Success 200 {object} response.Response
// @Router /api/v1/{{.Module}}/{{.ResourceName}} [get]
func (h *{{.ModelName}}Handler) List(c *gin.Context) {
	var req model.{{.ModelName}}ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	
	list, total, err := h.service.List(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	
	response.SuccessWithData(c, gin.H{
		"list":  model.To{{.ModelName}}ResponseList(list),
		"total": total,
		"page":  req.Page,
		"page_size": req.PageSize,
	})
}

// Get 获取{{.Title}}详情
// @Summary 获取{{.Title}}详情
// @Description 根据 ID 获取{{.Title}}详情
// @Tags {{.Title}}
// @Accept json
// @Produce json
// @Param id path int true "{{.Title}} ID"
// @Success 200 {object} response.Response
// @Router /api/v1/{{.Module}}/{{.ResourceName}}/{id} [get]
func (h *{{.ModelName}}Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "ID 格式错误")
		return
	}
	
	{{.ModelNameCamel}}, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, 404, err.Error())
		return
	}
	
	response.Success(c, model.To{{.ModelName}}Response({{.ModelNameCamel}}))
}

// Create 创建{{.Title}}
// @Summary 创建{{.Title}}
// @Description 创建新的{{.Title}}
// @Tags {{.Title}}
// @Accept json
// @Produce json
// @Param body body model.{{.ModelName}}CreateRequest true "{{.Title}}信息"
// @Success 200 {object} response.Response
// @Router /api/v1/{{.Module}}/{{.ResourceName}} [post]
func (h *{{.ModelName}}Handler) Create(c *gin.Context) {
	var req model.{{.ModelName}}CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}
	
	{{.ModelNameCamel}}, err := h.service.Create(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	
	response.Success(c, model.To{{.ModelName}}Response({{.ModelNameCamel}}))
}

// Update 更新{{.Title}}
// @Summary 更新{{.Title}}
// @Description 更新{{.Title}}信息
// @Tags {{.Title}}
// @Accept json
// @Produce json
// @Param id path int true "{{.Title}} ID"
// @Param body body model.{{.ModelName}}UpdateRequest true "{{.Title}}信息"
// @Success 200 {object} response.Response
// @Router /api/v1/{{.Module}}/{{.ResourceName}}/{id} [put]
func (h *{{.ModelName}}Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, 400, "ID 格式错误")
		return
	}
	
	var req model.{{.ModelName}}UpdateRequest
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

// Delete 删除{{.Title}}
// @Summary 删除{{.Title}}
// @Description 删除{{.Title}}
// @Tags {{.Title}}
// @Accept json
// @Produce json
// @Param id path int true "{{.Title}} ID"
// @Success 200 {object} response.Response
// @Router /api/v1/{{.Module}}/{{.ResourceName}}/{id} [delete]
func (h *{{.ModelName}}Handler) Delete(c *gin.Context) {
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
`

	funcMap := template.FuncMap{
		"toPascalCase": toPascalCase,
	}

	t, err := template.New("handler").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
