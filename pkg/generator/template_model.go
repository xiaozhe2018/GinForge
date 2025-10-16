package generator

import (
	"bytes"
	"strings"
	"text/template"
)

// renderModelTemplate 渲染 Model 模板
func renderModelTemplate(data *TemplateData) (string, error) {
	tmpl := `package model

import (
{{- if needTimeImport .Fields}}
	"time"
{{- end}}
)

// {{.ModelName}} {{.Title}}
type {{.ModelName}} struct {
{{- range .Fields}}
	{{toPascalCase .Name}} {{.GoType}} ` + "`" + `{{generateGormTag .}} {{generateJSONTag .}}` + "`" + ` // {{.Comment}}
{{- end}}
}

// TableName 指定表名
func ({{.ModelNameCamel}} *{{.ModelName}}) TableName() string {
	return "{{.Table}}"
}

// {{.ModelName}}ListRequest {{.Title}}列表请求
type {{.ModelName}}ListRequest struct {
	Page     int    ` + "`" + `form:"page" binding:"omitempty,min=1"` + "`" + `
	PageSize int    ` + "`" + `form:"page_size" binding:"omitempty,min=1,max=100"` + "`" + `
{{- if .HasSearch}}
	Keyword  string ` + "`" + `form:"keyword"` + "`" + `
{{- end}}
{{- if .HasSort}}
	SortBy   string ` + "`" + `form:"sort_by"` + "`" + `
	SortOrder string ` + "`" + `form:"sort_order" binding:"omitempty,oneof=asc desc"` + "`" + `
{{- end}}
}

// {{.ModelName}}CreateRequest 创建{{.Title}}请求
type {{.ModelName}}CreateRequest struct {
{{- range .Fields}}
{{- if and .FormVisible (not .AutoIncrement)}}
	{{toPascalCase .Name}} {{.GoType}} ` + "`" + `json:"{{toSnakeCase .Name}}"{{if .Validations}} binding:"{{joinValidations .Validations}}"{{end}}` + "`" + ` // {{.Comment}}
{{- end}}
{{- end}}
}

// {{.ModelName}}UpdateRequest 更新{{.Title}}请求
type {{.ModelName}}UpdateRequest struct {
{{- range .Fields}}
{{- if and .FormVisible (not .IsPrimaryKey) (not .AutoIncrement)}}
	{{toPascalCase .Name}} {{.GoType}} ` + "`" + `json:"{{toSnakeCase .Name}}"{{if .Validations}} binding:"{{joinValidations .Validations}}"{{end}}` + "`" + ` // {{.Comment}}
{{- end}}
{{- end}}
}

// {{.ModelName}}Response {{.Title}}响应
type {{.ModelName}}Response struct {
{{- range .Fields}}
{{- if .ListVisible}}
	{{toPascalCase .Name}} {{.GoType}} ` + "`" + `json:"{{toSnakeCase .Name}}"` + "`" + ` // {{.Comment}}
{{- end}}
{{- end}}
}

// To{{.ModelName}}Response 转换为响应对象
func To{{.ModelName}}Response({{.ModelNameCamel}} *{{.ModelName}}) *{{.ModelName}}Response {
	if {{.ModelNameCamel}} == nil {
		return nil
	}
	
	return &{{.ModelName}}Response{
{{- range .Fields}}
{{- if .ListVisible}}
		{{toPascalCase .Name}}: {{$.ModelNameCamel}}.{{toPascalCase .Name}},
{{- end}}
{{- end}}
	}
}

// To{{.ModelName}}ResponseList 批量转换为响应对象
func To{{.ModelName}}ResponseList(list []*{{.ModelName}}) []*{{.ModelName}}Response {
	result := make([]*{{.ModelName}}Response, 0, len(list))
	for _, item := range list {
		result = append(result, To{{.ModelName}}Response(item))
	}
	return result
}
`

	funcMap := template.FuncMap{
		"toPascalCase":    toPascalCase,
		"toSnakeCase":     toSnakeCase,
		"stripPointer":    stripPointer,
		"generateGormTag": generateGormTag,
		"generateJSONTag": generateJSONTag,
		"joinValidations": joinValidations,
		"needTimeImport":  needTimeImport,
	}

	t, err := template.New("model").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// 辅助函数
func joinValidations(validations []string) string {
	return strings.Join(validations, ",")
}

func needTimeImport(fields []FieldConfig) bool {
	for _, field := range fields {
		if strings.Contains(field.GoType, "time.Time") {
			return true
		}
	}
	return false
}
