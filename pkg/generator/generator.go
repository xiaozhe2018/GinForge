package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"goweb/pkg/config"
	"goweb/pkg/db"
)

// Generator 代码生成器
type Generator struct {
	db     *gorm.DB
	config *config.Config
}

// New 创建生成器实例
func New() (*Generator, error) {
	// 加载配置
	cfg := config.New()

	// 初始化数据库
	database, err := db.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化数据库失败: %w", err)
	}

	return &Generator{
		db:     database,
		config: cfg,
	}, nil
}

// ListTables 列出所有表
func (g *Generator) ListTables() ([]string, error) {
	var tables []string

	// 获取数据库名
	dbName := g.config.GetString("database.database")

	err := g.db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = ?", dbName).
		Scan(&tables).Error

	if err != nil {
		return nil, err
	}

	return tables, nil
}

// GenerateConfigFromTable 从数据库表生成配置
func (g *Generator) GenerateConfigFromTable(tableName, moduleName string) (*CRUDConfig, error) {
	// 读取表结构
	tableInfo, err := g.getTableInfo(tableName)
	if err != nil {
		return nil, err
	}

	// 生成模型名（去除前缀，转 PascalCase）
	modelName := g.tableNameToModelName(tableName)
	modelNameCamel := toCamelCase(modelName)
	resourceName := toPlural(toSnakeCase(modelName))

	config := &CRUDConfig{
		Table:          tableName,
		Module:         moduleName,
		ModelName:      modelName,
		ModelNameCamel: modelNameCamel,
		ResourceName:   resourceName,
		Fields:         []FieldConfig{},
		Features: Features{
			SoftDelete:  g.hasColumn(tableInfo, "deleted_at"),
			Timestamps:  g.hasColumn(tableInfo, "created_at") && g.hasColumn(tableInfo, "updated_at"),
			Pagination:  true,
			Search:      true,
			Sort:        true,
			Export:      false,
			Import:      false,
			BatchDelete: false,
		},
		Frontend: FrontendConfig{
			Title:      g.generateTitle(modelName),
			Icon:       "Document",
			ShowInMenu: true,
		},
	}

	// 转换字段
	for _, col := range tableInfo.Columns {
		field := g.columnToField(col)
		config.Fields = append(config.Fields, field)
	}

	return config, nil
}

// getTableInfo 获取表信息
func (g *Generator) getTableInfo(tableName string) (*TableInfo, error) {
	dbName := g.config.GetString("database.database")

	var columns []ColumnInfo
	err := g.db.Raw(`
		SELECT 
			COLUMN_NAME as name,
			COLUMN_TYPE as type,
			IS_NULLABLE as nullable,
			COLUMN_KEY as 'key',
			COLUMN_DEFAULT as 'default',
			EXTRA as extra,
			COLUMN_COMMENT as comment
		FROM information_schema.COLUMNS
		WHERE table_schema = ? AND table_name = ?
		ORDER BY ORDINAL_POSITION
	`, dbName, tableName).Scan(&columns).Error

	if err != nil {
		return nil, err
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("表 %s 不存在或没有列", tableName)
	}

	return &TableInfo{
		Name:    tableName,
		Columns: columns,
	}, nil
}

// columnToField 将数据库列转换为字段配置
func (g *Generator) columnToField(col ColumnInfo) FieldConfig {
	nullable := col.Nullable == "YES"
	goType := GetGoType(col.Type, nullable)
	tsType := GetTSType(goType)
	formType := GetFormType(col.Name)

	// 生成标签
	label := g.generateLabel(col.Name, col.Comment)

	// 生成验证规则
	validations := g.generateValidations(col)

	// 判断是否在列表和表单中显示
	listVisible := !g.isHiddenInList(col.Name)
	formVisible := !g.isHiddenInForm(col.Name)

	// 判断是否可搜索和排序
	searchable := g.isSearchable(col.Name, col.Type)
	sortable := true

	field := FieldConfig{
		Name:          col.Name,
		Type:          col.Type,
		GoType:        goType,
		TSType:        tsType,
		Nullable:      nullable,
		IsPrimaryKey:  col.Key == "PRI",
		AutoIncrement: strings.Contains(col.Extra, "auto_increment"),
		DefaultValue:  fmt.Sprintf("%v", col.Default),
		Comment:       col.Comment,
		Validations:   validations,
		Label:         label,
		FormType:      formType,
		ListVisible:   listVisible,
		FormVisible:   formVisible,
		Searchable:    searchable,
		Sortable:      sortable,
	}

	return field
}

// hasColumn 检查表是否有某个字段
func (g *Generator) hasColumn(tableInfo *TableInfo, colName string) bool {
	for _, col := range tableInfo.Columns {
		if col.Name == colName {
			return true
		}
	}
	return false
}

// tableNameToModelName 表名转模型名
func (g *Generator) tableNameToModelName(tableName string) string {
	// 去除常见前缀
	prefixes := []string{"admin_", "user_", "sys_", "tb_", "t_"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(tableName, prefix) {
			tableName = strings.TrimPrefix(tableName, prefix)
			break
		}
	}

	// 转 PascalCase
	return toPascalCase(tableName)
}

// generateLabel 生成标签
func (g *Generator) generateLabel(fieldName, comment string) string {
	if comment != "" {
		return comment
	}

	// 默认标签映射
	labelMap := map[string]string{
		"id":          "ID",
		"name":        "名称",
		"title":       "标题",
		"content":     "内容",
		"description": "描述",
		"status":      "状态",
		"sort":        "排序",
		"created_at":  "创建时间",
		"updated_at":  "更新时间",
		"deleted_at":  "删除时间",
	}

	if label, ok := labelMap[fieldName]; ok {
		return label
	}

	return toPascalCase(fieldName)
}

// generateTitle 生成页面标题
func (g *Generator) generateTitle(modelName string) string {
	titleMap := map[string]string{
		"Article":    "文章管理",
		"User":       "用户管理",
		"Category":   "分类管理",
		"Tag":        "标签管理",
		"Comment":    "评论管理",
		"File":       "文件管理",
		"Config":     "配置管理",
		"Log":        "日志管理",
		"Role":       "角色管理",
		"Permission": "权限管理",
		"Menu":       "菜单管理",
	}

	if title, ok := titleMap[modelName]; ok {
		return title
	}

	return modelName + "管理"
}

// generateValidations 生成验证规则
func (g *Generator) generateValidations(col ColumnInfo) []string {
	validations := []string{}

	// 非空验证
	if col.Nullable == "NO" && col.Key != "PRI" && !strings.Contains(col.Extra, "auto_increment") {
		validations = append(validations, "required")
	}

	// 根据字段名添加特定验证
	switch col.Name {
	case "email":
		validations = append(validations, "email")
	case "phone":
		validations = append(validations, "len:11")
	case "password":
		validations = append(validations, "min:6")
	case "url":
		validations = append(validations, "url")
	}

	// 根据类型添加验证
	if strings.Contains(col.Type, "varchar") {
		// 提取长度
		length := extractLength(col.Type)
		if length > 0 {
			validations = append(validations, fmt.Sprintf("max:%d", length))
		}
	}

	return validations
}

// isHiddenInList 是否在列表中隐藏
func (g *Generator) isHiddenInList(fieldName string) bool {
	hiddenFields := []string{"password", "deleted_at", "content", "description"}
	for _, hidden := range hiddenFields {
		if fieldName == hidden {
			return true
		}
	}
	return false
}

// isHiddenInForm 是否在表单中隐藏
func (g *Generator) isHiddenInForm(fieldName string) bool {
	hiddenFields := []string{"id", "created_at", "updated_at", "deleted_at"}
	for _, hidden := range hiddenFields {
		if fieldName == hidden {
			return true
		}
	}
	return false
}

// isSearchable 是否可搜索
func (g *Generator) isSearchable(fieldName, fieldType string) bool {
	// 字符串类型的字段可搜索
	if strings.Contains(fieldType, "char") || strings.Contains(fieldType, "text") {
		return true
	}

	// 特定字段可搜索
	searchableFields := []string{"id", "name", "title", "email", "phone"}
	for _, field := range searchableFields {
		if fieldName == field {
			return true
		}
	}

	return false
}

// SaveConfigToFile 保存配置到文件
func (g *Generator) SaveConfigToFile(config *CRUDConfig, outputDir string) (string, error) {
	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	// 生成文件路径
	filename := toSnakeCase(config.ModelName) + ".yaml"
	filepath := filepath.Join(outputDir, filename)

	// 序列化为 YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	// 写入文件
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", err
	}

	return filepath, nil
}

// LoadConfigFromFile 从文件加载配置
func LoadConfigFromFile(filepath string) (*CRUDConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config CRUDConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// 生成 camelCase 名称
	config.ModelNameCamel = toCamelCase(config.ModelName)

	return &config, nil
}

// extractLength 提取类型长度
func extractLength(typeStr string) int {
	// varchar(255) -> 255
	if start := strings.Index(typeStr, "("); start != -1 {
		if end := strings.Index(typeStr, ")"); end != -1 {
			lengthStr := typeStr[start+1 : end]
			var length int
			fmt.Sscanf(lengthStr, "%d", &length)
			return length
		}
	}
	return 0
}
