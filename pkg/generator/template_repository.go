package generator

import (
	"bytes"
	"strings"
	"text/template"
)

// renderRepositoryTemplate 渲染 Repository 模板
func renderRepositoryTemplate(data *TemplateData) (string, error) {
	tmpl := `package repository

import (
	"gorm.io/gorm"
	"goweb/services/{{.Module}}-api/internal/model"
)

// {{.ModelName}}Repository {{.Title}} Repository
type {{.ModelName}}Repository struct {
	db *gorm.DB
}

// New{{.ModelName}}Repository 创建 Repository 实例
func New{{.ModelName}}Repository(db *gorm.DB) *{{.ModelName}}Repository {
	return &{{.ModelName}}Repository{
		db: db,
	}
}

// Create 创建{{.Title}}
func (r *{{.ModelName}}Repository) Create({{.ModelNameCamel}} *model.{{.ModelName}}) error {
	return r.db.Create({{.ModelNameCamel}}).Error
}

// GetByID 根据 ID 获取{{.Title}}
func (r *{{.ModelName}}Repository) GetByID(id uint64) (*model.{{.ModelName}}, error) {
	var {{.ModelNameCamel}} model.{{.ModelName}}
	err := r.db.First(&{{.ModelNameCamel}}, id).Error
	if err != nil {
		return nil, err
	}
	return &{{.ModelNameCamel}}, nil
}

// Update 更新{{.Title}}
func (r *{{.ModelName}}Repository) Update({{.ModelNameCamel}} *model.{{.ModelName}}) error {
	return r.db.Save({{.ModelNameCamel}}).Error
}

// Delete 删除{{.Title}}
func (r *{{.ModelName}}Repository) Delete(id uint64) error {
{{- if .HasSoftDelete}}
	return r.db.Delete(&model.{{.ModelName}}{}, id).Error
{{- else}}
	return r.db.Unscoped().Delete(&model.{{.ModelName}}{}, id).Error
{{- end}}
}

// List 获取{{.Title}}列表
func (r *{{.ModelName}}Repository) List(req *model.{{.ModelName}}ListRequest) ([]*model.{{.ModelName}}, int64, error) {
	var list []*model.{{.ModelName}}
	var total int64
	
	db := r.db.Model(&model.{{.ModelName}}{})
	
{{- if .HasSearch}}
	// 搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		db = db.Where("{{getSearchCondition .Fields}}", keyword{{getSearchParams .Fields}})
	}
{{- end}}
	
	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
{{- if .HasSort}}
	// 排序
	if req.SortBy != "" {
		order := req.SortBy
		if req.SortOrder == "desc" {
			order += " DESC"
		}
		db = db.Order(order)
	} else {
		db = db.Order("{{getPrimaryKeyName .Fields}} DESC")
	}
{{- else}}
	db = db.Order("{{getPrimaryKeyName .Fields}} DESC")
{{- end}}
	
{{- if .HasPagination}}
	// 分页
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		db = db.Offset(offset).Limit(req.PageSize)
	}
{{- end}}
	
	err := db.Find(&list).Error
	return list, total, err
}

// Exists 检查{{.Title}}是否存在
func (r *{{.ModelName}}Repository) Exists(id uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.{{.ModelName}}{}).Where("{{getPrimaryKeyName .Fields}} = ?", id).Count(&count).Error
	return count > 0, err
}

{{- if .HasSoftDelete}}

// Restore 恢复已删除的{{.Title}}
func (r *{{.ModelName}}Repository) Restore(id uint64) error {
	return r.db.Model(&model.{{.ModelName}}{}).Unscoped().Where("{{getPrimaryKeyName .Fields}} = ?", id).Update("deleted_at", nil).Error
}

// ForceDelete 永久删除{{.Title}}
func (r *{{.ModelName}}Repository) ForceDelete(id uint64) error {
	return r.db.Unscoped().Delete(&model.{{.ModelName}}{}, id).Error
}
{{- end}}
`

	funcMap := template.FuncMap{
		"getPrimaryKeyName":  getPrimaryKeyName,
		"getSearchCondition": getSearchCondition,
		"getSearchParams":    getSearchParams,
	}

	t, err := template.New("repository").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getPrimaryKeyName 获取主键名称
func getPrimaryKeyName(fields []FieldConfig) string {
	for _, field := range fields {
		if field.IsPrimaryKey {
			return field.Name
		}
	}
	return "id"
}

// getSearchCondition 获取搜索条件
func getSearchCondition(fields []FieldConfig) string {
	searchFields := getSearchableFields(fields)
	if len(searchFields) == 0 {
		return "1=1"
	}

	conditions := []string{}
	for _, field := range searchFields {
		conditions = append(conditions, field.Name+" LIKE ?")
	}

	return strings.Join(conditions, " OR ")
}

// getSearchParams 获取搜索参数
func getSearchParams(fields []FieldConfig) string {
	searchFields := getSearchableFields(fields)
	if len(searchFields) <= 1 {
		return ""
	}

	params := []string{}
	for i := 1; i < len(searchFields); i++ {
		params = append(params, ", keyword")
	}

	return strings.Join(params, "")
}
