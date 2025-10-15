package generator

import "time"

// CRUDConfig CRUD 生成配置
type CRUDConfig struct {
	// 基础信息
	Table          string `yaml:"table"`         // 数据库表名
	Module         string `yaml:"module"`        // 模块名称（admin/user/file）
	ModelName      string `yaml:"model_name"`    // 模型名称（PascalCase）
	ModelNameCamel string `yaml:"-"`             // 模型名称（camelCase）
	ResourceName   string `yaml:"resource_name"` // 资源名称（复数形式，用于 URL）

	// 字段列表
	Fields []FieldConfig `yaml:"fields"`

	// 生成选项
	Options GenerateOptions `yaml:"options"`

	// 功能特性
	Features Features `yaml:"features"`

	// 前端配置
	Frontend FrontendConfig `yaml:"frontend"`
}

// FieldConfig 字段配置
type FieldConfig struct {
	// 数据库字段
	Name   string `yaml:"name"`    // 字段名（snake_case）
	Type   string `yaml:"type"`    // 数据库类型
	GoType string `yaml:"go_type"` // Go 类型
	TSType string `yaml:"ts_type"` // TypeScript 类型

	// 字段属性
	Nullable      bool   `yaml:"nullable"`       // 是否可为空
	IsPrimaryKey  bool   `yaml:"is_primary_key"` // 是否主键
	AutoIncrement bool   `yaml:"auto_increment"` // 是否自增
	DefaultValue  string `yaml:"default_value"`  // 默认值
	Comment       string `yaml:"comment"`        // 字段注释

	// 验证规则
	Validations []string `yaml:"validations"` // 验证规则（required, email, min:6 等）

	// UI 配置
	Label       string `yaml:"label"`        // 显示标签
	FormType    string `yaml:"form_type"`    // 表单类型（input/textarea/select/date/switch）
	ListVisible bool   `yaml:"list_visible"` // 列表中是否显示
	FormVisible bool   `yaml:"form_visible"` // 表单中是否显示
	Searchable  bool   `yaml:"searchable"`   // 是否可搜索
	Sortable    bool   `yaml:"sortable"`     // 是否可排序

	// 关联关系
	Relation *Relation `yaml:"relation,omitempty"` // 关联关系
}

// Relation 关联关系
type Relation struct {
	Type         string `yaml:"type"`          // belongs_to, has_many, has_one
	Model        string `yaml:"model"`         // 关联模型
	ForeignKey   string `yaml:"foreign_key"`   // 外键
	DisplayField string `yaml:"display_field"` // 显示字段
}

// Features 功能特性
type Features struct {
	SoftDelete  bool `yaml:"soft_delete"`  // 软删除
	Timestamps  bool `yaml:"timestamps"`   // 时间戳
	Pagination  bool `yaml:"pagination"`   // 分页
	Search      bool `yaml:"search"`       // 搜索
	Sort        bool `yaml:"sort"`         // 排序
	Export      bool `yaml:"export"`       // 导出
	Import      bool `yaml:"import"`       // 导入
	BatchDelete bool `yaml:"batch_delete"` // 批量删除
}

// FrontendConfig 前端配置
type FrontendConfig struct {
	Title      string `yaml:"title"`        // 页面标题
	Icon       string `yaml:"icon"`         // 菜单图标
	ShowInMenu bool   `yaml:"show_in_menu"` // 是否显示在菜单
	MenuParent string `yaml:"menu_parent"`  // 父菜单
}

// GenerateOptions 生成选项
type GenerateOptions struct {
	OutputDir    string `yaml:"output_dir"`    // 输出目录
	WithFrontend bool   `yaml:"with_frontend"` // 生成前端代码
	Force        bool   `yaml:"force"`         // 强制覆盖
	DryRun       bool   `yaml:"dry_run"`       // 预览模式
	Verbose      bool   `yaml:"verbose"`       // 详细输出
}

// GenerateResult 生成结果
type GenerateResult struct {
	Files  []FileResult `json:"files"`
	Errors []string     `json:"errors"`
}

// FileResult 文件生成结果
type FileResult struct {
	Path    string `json:"path"`
	Created bool   `json:"created"`
	Skipped bool   `json:"skipped"`
	Error   string `json:"error,omitempty"`
}

// TableInfo 数据库表信息
type TableInfo struct {
	Name    string
	Comment string
	Columns []ColumnInfo
}

// ColumnInfo 数据库列信息
type ColumnInfo struct {
	Name     string
	Type     string
	Nullable string // YES or NO
	Key      string // PRI, UNI, MUL
	Default  interface{}
	Extra    string // auto_increment
	Comment  string
}

// TemplateData 模板数据
type TemplateData struct {
	// 基础信息
	Table          string
	Module         string
	ModelName      string
	ModelNameCamel string
	ResourceName   string
	PackageName    string

	// 字段信息
	Fields     []FieldConfig
	PrimaryKey *FieldConfig

	// 功能特性
	HasSoftDelete bool
	HasTimestamps bool
	HasPagination bool
	HasSearch     bool
	HasSort       bool

	// 前端信息
	Title string
	Icon  string

	// 导入包
	Imports []string

	// 时间戳
	GeneratedAt string
}

// 类型映射
var (
	// MySQL 类型到 Go 类型的映射
	MySQLToGoType = map[string]string{
		"tinyint":    "int8",
		"smallint":   "int16",
		"mediumint":  "int32",
		"int":        "int",
		"integer":    "int",
		"bigint":     "int64",
		"float":      "float32",
		"double":     "float64",
		"decimal":    "float64",
		"char":       "string",
		"varchar":    "string",
		"text":       "string",
		"tinytext":   "string",
		"mediumtext": "string",
		"longtext":   "string",
		"date":       "time.Time",
		"datetime":   "time.Time",
		"timestamp":  "time.Time",
		"time":       "time.Time",
		"year":       "int",
		"json":       "string",
		"blob":       "[]byte",
		"tinyblob":   "[]byte",
		"mediumblob": "[]byte",
		"longblob":   "[]byte",
		"enum":       "string",
		"set":        "string",
	}

	// Go 类型到 TypeScript 类型的映射
	GoToTSType = map[string]string{
		"int":       "number",
		"int8":      "number",
		"int16":     "number",
		"int32":     "number",
		"int64":     "number",
		"uint":      "number",
		"uint8":     "number",
		"uint16":    "number",
		"uint32":    "number",
		"uint64":    "number",
		"float32":   "number",
		"float64":   "number",
		"string":    "string",
		"bool":      "boolean",
		"time.Time": "string",
		"[]byte":    "string",
	}

	// 字段名到表单类型的映射
	FieldNameToFormType = map[string]string{
		"password":    "password",
		"email":       "email",
		"phone":       "tel",
		"url":         "url",
		"avatar":      "upload",
		"image":       "upload",
		"file":        "upload",
		"content":     "editor",
		"description": "textarea",
		"remark":      "textarea",
		"status":      "switch",
		"is_":         "switch", // 前缀匹配
		"type":        "select",
		"category":    "select",
		"date":        "date",
		"time":        "datetime",
		"created_at":  "datetime",
		"updated_at":  "datetime",
	}
)

// GetGoType 根据 MySQL 类型获取 Go 类型
func GetGoType(mysqlType string, nullable bool) string {
	// 提取基础类型（去除长度等）
	for k, v := range MySQLToGoType {
		if len(mysqlType) >= len(k) && mysqlType[:len(k)] == k {
			goType := v

			// 如果可为空，使用指针类型
			if nullable && goType != "string" && goType != "[]byte" {
				return "*" + goType
			}

			return goType
		}
	}

	return "string"
}

// GetTSType 根据 Go 类型获取 TypeScript 类型
func GetTSType(goType string) string {
	// 去除指针
	goType = stripPointer(goType)

	if tsType, ok := GoToTSType[goType]; ok {
		return tsType
	}

	return "any"
}

// GetFormType 根据字段名获取表单类型
func GetFormType(fieldName string) string {
	fieldNameLower := toLower(fieldName)

	// 精确匹配
	if formType, ok := FieldNameToFormType[fieldNameLower]; ok {
		return formType
	}

	// 前缀匹配
	for prefix, formType := range FieldNameToFormType {
		if len(fieldNameLower) > len(prefix) && fieldNameLower[:len(prefix)] == prefix {
			return formType
		}
	}

	return "input"
}

// stripPointer 去除指针符号
func stripPointer(s string) string {
	if len(s) > 0 && s[0] == '*' {
		return s[1:]
	}
	return s
}

// toLower 转小写
func toLower(s string) string {
	return string([]rune(s))
}

// Now 获取当前时间
func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
