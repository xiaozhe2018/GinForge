package generator

import (
	"regexp"
	"strings"
	"unicode"
)

// toPascalCase 转换为 PascalCase
func toPascalCase(s string) string {
	// 先转为 words
	words := splitWords(s)

	result := ""
	for _, word := range words {
		if len(word) > 0 {
			result += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return result
}

// toCamelCase 转换为 camelCase
func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) == 0 {
		return ""
	}

	return strings.ToLower(pascal[:1]) + pascal[1:]
}

// toSnakeCase 转换为 snake_case
func toSnakeCase(s string) string {
	// 处理连续大写字母
	re := regexp.MustCompile("([A-Z]+)([A-Z][a-z])")
	s = re.ReplaceAllString(s, "${1}_${2}")

	// 在大写字母前添加下划线
	re = regexp.MustCompile("([a-z0-9])([A-Z])")
	s = re.ReplaceAllString(s, "${1}_${2}")

	// 替换空格、横线等为下划线
	re = regexp.MustCompile("[\\s-]+")
	s = re.ReplaceAllString(s, "_")

	// 转小写
	return strings.ToLower(s)
}

// toKebabCase 转换为 kebab-case
func toKebabCase(s string) string {
	snake := toSnakeCase(s)
	return strings.ReplaceAll(snake, "_", "-")
}

// splitWords 分割单词
func splitWords(s string) []string {
	// 处理 snake_case 和 kebab-case
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")

	// 处理 PascalCase 和 camelCase
	var words []string
	var currentWord strings.Builder

	for i, r := range s {
		if unicode.IsSpace(r) {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			continue
		}

		if i > 0 && unicode.IsUpper(r) {
			prevRune := rune(s[i-1])
			if !unicode.IsUpper(prevRune) {
				// 遇到大写字母，且前一个不是大写，分词
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		}

		currentWord.WriteRune(r)
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

// toPlural 转复数形式（简单实现）
func toPlural(s string) string {
	// 特殊情况
	irregulars := map[string]string{
		"person": "people",
		"child":  "children",
		"man":    "men",
		"woman":  "women",
		"tooth":  "teeth",
		"foot":   "feet",
		"mouse":  "mice",
		"goose":  "geese",
	}

	if plural, ok := irregulars[s]; ok {
		return plural
	}

	// 以 s, x, z, ch, sh 结尾
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") ||
		strings.HasSuffix(s, "z") || strings.HasSuffix(s, "ch") ||
		strings.HasSuffix(s, "sh") {
		return s + "es"
	}

	// 以辅音字母 + y 结尾
	if len(s) > 1 && s[len(s)-1] == 'y' {
		prevChar := s[len(s)-2]
		if !isVowel(prevChar) {
			return s[:len(s)-1] + "ies"
		}
	}

	// 以 f 或 fe 结尾
	if strings.HasSuffix(s, "f") {
		return s[:len(s)-1] + "ves"
	}
	if strings.HasSuffix(s, "fe") {
		return s[:len(s)-2] + "ves"
	}

	// 默认加 s
	return s + "s"
}

// isVowel 是否元音字母
func isVowel(c byte) bool {
	vowels := "aeiouAEIOU"
	return strings.ContainsRune(vowels, rune(c))
}

// ucFirst 首字母大写
func ucFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// lcFirst 首字母小写
func lcFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// indent 缩进
func indent(s string, level int) string {
	indentStr := strings.Repeat("\t", level)
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = indentStr + line
		}
	}

	return strings.Join(lines, "\n")
}

// wrapQuotes 添加引号
func wrapQuotes(s string) string {
	return `"` + s + `"`
}

// joinStrings 连接字符串
func joinStrings(items []string, sep string) string {
	return strings.Join(items, sep)
}

// contains 检查切片是否包含元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// unique 去重
func unique(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// filter 过滤
func filter(slice []string, fn func(string) bool) []string {
	result := []string{}

	for _, item := range slice {
		if fn(item) {
			result = append(result, item)
		}
	}

	return result
}

// mapStrings 映射
func mapStrings(slice []string, fn func(string) string) []string {
	result := make([]string, len(slice))

	for i, item := range slice {
		result[i] = fn(item)
	}

	return result
}

// getPrimaryKeyField 获取主键字段
func getPrimaryKeyField(fields []FieldConfig) *FieldConfig {
	for i := range fields {
		if fields[i].IsPrimaryKey {
			return &fields[i]
		}
	}
	return nil
}

// getVisibleFields 获取可见字段
func getVisibleFields(fields []FieldConfig, context string) []FieldConfig {
	result := []FieldConfig{}

	for _, field := range fields {
		if context == "list" && field.ListVisible {
			result = append(result, field)
		} else if context == "form" && field.FormVisible {
			result = append(result, field)
		}
	}

	return result
}

// getSearchableFields 获取可搜索字段
func getSearchableFields(fields []FieldConfig) []FieldConfig {
	result := []FieldConfig{}

	for _, field := range fields {
		if field.Searchable {
			result = append(result, field)
		}
	}

	return result
}

// hasFeature 检查是否有某个特性
func hasFeature(features Features, feature string) bool {
	switch feature {
	case "soft_delete":
		return features.SoftDelete
	case "timestamps":
		return features.Timestamps
	case "pagination":
		return features.Pagination
	case "search":
		return features.Search
	case "sort":
		return features.Sort
	case "export":
		return features.Export
	case "import":
		return features.Import
	case "batch_delete":
		return features.BatchDelete
	default:
		return false
	}
}

// generateImports 生成导入语句
func generateImports(config *CRUDConfig) []string {
	imports := []string{
		"github.com/gin-gonic/gin",
		"goweb/pkg/response",
	}

	// 检查是否需要 time 包
	needTime := false
	for _, field := range config.Fields {
		if field.GoType == "time.Time" || field.GoType == "*time.Time" {
			needTime = true
			break
		}
	}

	if needTime {
		imports = append(imports, "time")
	}

	// 检查是否需要 errors 包
	imports = append(imports, "errors")

	// 检查是否需要 fmt 包
	imports = append(imports, "fmt")

	return unique(imports)
}

// generateValidationTag 生成验证标签
func generateValidationTag(field FieldConfig) string {
	if len(field.Validations) == 0 {
		return ""
	}

	return `binding:"` + strings.Join(field.Validations, ",") + `"`
}

// generateGormTag 生成 GORM 标签
func generateGormTag(field FieldConfig) string {
	tags := []string{}

	// 字段名
	tags = append(tags, "column:"+field.Name)

	// 类型
	tags = append(tags, "type:"+field.Type)

	// 主键
	if field.IsPrimaryKey {
		tags = append(tags, "primaryKey")
	}

	// 自增
	if field.AutoIncrement {
		tags = append(tags, "autoIncrement")
	}

	// 非空
	if !field.Nullable {
		tags = append(tags, "not null")
	}

	// 默认值
	if field.DefaultValue != "" && field.DefaultValue != "<nil>" {
		tags = append(tags, "default:"+field.DefaultValue)
	}

	return `gorm:"` + strings.Join(tags, ";") + `"`
}

// generateJSONTag 生成 JSON 标签
func generateJSONTag(field FieldConfig) string {
	jsonName := toSnakeCase(field.Name)

	// 可为空的字段添加 omitempty
	if field.Nullable {
		return `json:"` + jsonName + `,omitempty"`
	}

	return `json:"` + jsonName + `"`
}
