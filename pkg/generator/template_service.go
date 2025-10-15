package generator

import (
	"bytes"
	"text/template"
)

// renderServiceTemplate 渲染 Service 模板
func renderServiceTemplate(data *TemplateData) (string, error) {
	tmpl := `package service

import (
	"errors"
	"gorm.io/gorm"
	
	"goweb/pkg/logger"
	"goweb/services/{{.Module}}-api/internal/model"
	"goweb/services/{{.Module}}-api/internal/repository"
)

// {{.ModelName}}Service {{.Title}} Service
type {{.ModelName}}Service struct {
	repo   *repository.{{.ModelName}}Repository
	logger logger.Logger
}

// New{{.ModelName}}Service 创建 Service 实例
func New{{.ModelName}}Service(repo *repository.{{.ModelName}}Repository, logger logger.Logger) *{{.ModelName}}Service {
	return &{{.ModelName}}Service{
		repo:   repo,
		logger: logger,
	}
}

// Create 创建{{.Title}}
func (s *{{.ModelName}}Service) Create(req *model.{{.ModelName}}CreateRequest) (*model.{{.ModelName}}, error) {
	{{.ModelNameCamel}} := &model.{{.ModelName}}{
{{- range .Fields}}
{{- if and .FormVisible (not .AutoIncrement) (not .IsPrimaryKey)}}
		{{toPascalCase .Name}}: req.{{toPascalCase .Name}},
{{- end}}
{{- end}}
	}
	
	if err := s.repo.Create({{.ModelNameCamel}}); err != nil {
		s.logger.Error("创建{{.Title}}失败", err)
		return nil, errors.New("创建{{.Title}}失败")
	}
	
	return {{.ModelNameCamel}}, nil
}

// GetByID 根据 ID 获取{{.Title}}
func (s *{{.ModelName}}Service) GetByID(id uint64) (*model.{{.ModelName}}, error) {
	{{.ModelNameCamel}}, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("{{.Title}}不存在")
		}
		s.logger.Error("获取{{.Title}}失败", err, "id", id)
		return nil, errors.New("获取{{.Title}}失败")
	}
	
	return {{.ModelNameCamel}}, nil
}

// Update 更新{{.Title}}
func (s *{{.ModelName}}Service) Update(id uint64, req *model.{{.ModelName}}UpdateRequest) error {
	// 检查{{.Title}}是否存在
	{{.ModelNameCamel}}, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("{{.Title}}不存在")
		}
		return errors.New("获取{{.Title}}失败")
	}
	
	// 更新字段
{{- range .Fields}}
{{- if and .FormVisible (not .IsPrimaryKey) (not .AutoIncrement)}}
	{{- if .Nullable}}
	if req.{{toPascalCase .Name}} != nil {
		{{$.ModelNameCamel}}.{{toPascalCase .Name}} = *req.{{toPascalCase .Name}}
	}
	{{- else}}
	{{$.ModelNameCamel}}.{{toPascalCase .Name}} = req.{{toPascalCase .Name}}
	{{- end}}
{{- end}}
{{- end}}
	
	if err := s.repo.Update({{.ModelNameCamel}}); err != nil {
		s.logger.Error("更新{{.Title}}失败", err, "id", id)
		return errors.New("更新{{.Title}}失败")
	}
	
	return nil
}

// Delete 删除{{.Title}}
func (s *{{.ModelName}}Service) Delete(id uint64) error {
	// 检查{{.Title}}是否存在
	exists, err := s.repo.Exists(id)
	if err != nil {
		s.logger.Error("检查{{.Title}}是否存在失败", err, "id", id)
		return errors.New("检查{{.Title}}是否存在失败")
	}
	
	if !exists {
		return errors.New("{{.Title}}不存在")
	}
	
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("删除{{.Title}}失败", err, "id", id)
		return errors.New("删除{{.Title}}失败")
	}
	
	return nil
}

// List 获取{{.Title}}列表
func (s *{{.ModelName}}Service) List(req *model.{{.ModelName}}ListRequest) ([]*model.{{.ModelName}}, int64, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	
	list, total, err := s.repo.List(req)
	if err != nil {
		s.logger.Error("获取{{.Title}}列表失败", err)
		return nil, 0, errors.New("获取{{.Title}}列表失败")
	}
	
	return list, total, nil
}
`

	funcMap := template.FuncMap{
		"toPascalCase": toPascalCase,
		"toSnakeCase":  toSnakeCase,
	}

	t, err := template.New("service").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
