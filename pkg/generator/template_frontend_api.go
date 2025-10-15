package generator

import (
	"bytes"
	"text/template"
)

// renderFrontendAPITemplate 渲染前端 API 模板
func renderFrontendAPITemplate(data *TemplateData) (string, error) {
	tmpl := `import request from '@/utils/request'

// ========== 类型定义 ==========

/**
 * {{.Title}}
 */
export interface {{.ModelName}} {
{{- range .Fields}}
{{- if .ListVisible}}
  {{toSnakeCase .Name}}: {{.TSType}} // {{.Comment}}
{{- end}}
{{- end}}
}

/**
 * {{.Title}}列表请求参数
 */
export interface {{.ModelName}}ListParams {
  page?: number
  page_size?: number
{{- if .HasSearch}}
  keyword?: string
{{- end}}
{{- if .HasSort}}
  sort_by?: string
  sort_order?: 'asc' | 'desc'
{{- end}}
}

/**
 * {{.Title}}列表响应
 */
export interface {{.ModelName}}ListResponse {
  list: {{.ModelName}}[]
  total: number
  page: number
  page_size: number
}

/**
 * 创建{{.Title}}请求参数
 */
export interface {{.ModelName}}CreateParams {
{{- range .Fields}}
{{- if and .FormVisible (not .AutoIncrement) (not .IsPrimaryKey)}}
  {{toSnakeCase .Name}}{{if .Nullable}}?{{end}}: {{.TSType}} // {{.Comment}}
{{- end}}
{{- end}}
}

/**
 * 更新{{.Title}}请求参数
 */
export interface {{.ModelName}}UpdateParams {
{{- range .Fields}}
{{- if and .FormVisible (not .IsPrimaryKey) (not .AutoIncrement)}}
  {{toSnakeCase .Name}}?: {{.TSType}} // {{.Comment}}
{{- end}}
{{- end}}
}

// ========== API 方法 ==========

/**
 * 获取{{.Title}}列表
 */
export const get{{.ModelName}}List = (params?: {{.ModelName}}ListParams) => {
  return request.get<{{.ModelName}}ListResponse>('/api/v1/{{.Module}}/{{.ResourceName}}', { params })
}

/**
 * 获取{{.Title}}详情
 */
export const get{{.ModelName}} = (id: number) => {
  return request.get<{{.ModelName}}>(` + "`" + `/api/v1/{{.Module}}/{{.ResourceName}}/${id}` + "`" + `)
}

/**
 * 创建{{.Title}}
 */
export const create{{.ModelName}} = (data: {{.ModelName}}CreateParams) => {
  return request.post<{{.ModelName}}>('/api/v1/{{.Module}}/{{.ResourceName}}', data)
}

/**
 * 更新{{.Title}}
 */
export const update{{.ModelName}} = (id: number, data: {{.ModelName}}UpdateParams) => {
  return request.put(` + "`" + `/api/v1/{{.Module}}/{{.ResourceName}}/${id}` + "`" + `, data)
}

/**
 * 删除{{.Title}}
 */
export const delete{{.ModelName}} = (id: number) => {
  return request.delete(` + "`" + `/api/v1/{{.Module}}/{{.ResourceName}}/${id}` + "`" + `)
}
`

	funcMap := template.FuncMap{
		"toSnakeCase": toSnakeCase,
	}

	t, err := template.New("frontend-api").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
