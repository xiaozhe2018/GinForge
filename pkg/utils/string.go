package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 检查字符串是否不为空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// IsBlank 检查字符串是否为空白
func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNotBlank 检查字符串是否不为空白
func IsNotBlank(s string) bool {
	return !IsBlank(s)
}

// Trim 去除字符串两端空白
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// TrimLeft 去除字符串左端空白
func TrimLeft(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

// TrimRight 去除字符串右端空白
func TrimRight(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// ToLower 转换为小写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper 转换为大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToTitle 转换为标题格式
func ToTitle(s string) string {
	return strings.Title(s)
}

// ToCamelCase 转换为驼峰格式
func ToCamelCase(s string) string {
	if IsEmpty(s) {
		return s
	}

	words := strings.Fields(strings.ReplaceAll(s, "_", " "))
	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += strings.Title(words[i])
	}

	return result
}

// ToSnakeCase 转换为蛇形格式
func ToSnakeCase(s string) string {
	if IsEmpty(s) {
		return s
	}

	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}

	return string(result)
}

// ToKebabCase 转换为短横线格式
func ToKebabCase(s string) string {
	return strings.ReplaceAll(ToSnakeCase(s), "_", "-")
}

// ToPascalCase 转换为帕斯卡格式
func ToPascalCase(s string) string {
	if IsEmpty(s) {
		return s
	}

	words := strings.Fields(strings.ReplaceAll(s, "_", " "))
	if len(words) == 0 {
		return s
	}

	var result string
	for _, word := range words {
		result += strings.Title(word)
	}

	return result
}

// Contains 检查字符串是否包含子字符串
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ContainsIgnoreCase 忽略大小写检查字符串是否包含子字符串
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// StartsWith 检查字符串是否以指定前缀开始
func StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// EndsWith 检查字符串是否以指定后缀结束
func EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Replace 替换字符串
func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// ReplaceFirst 替换第一个匹配的字符串
func ReplaceFirst(s, old, new string) string {
	return strings.Replace(s, old, new, 1)
}

// ReplaceLast 替换最后一个匹配的字符串
func ReplaceLast(s, old, new string) string {
	lastIndex := strings.LastIndex(s, old)
	if lastIndex == -1 {
		return s
	}
	return s[:lastIndex] + new + s[lastIndex+len(old):]
}

// Split 分割字符串
func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

// Join 连接字符串
func Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Repeat 重复字符串
func Repeat(s string, count int) string {
	return strings.Repeat(s, count)
}

// PadLeft 左填充字符串
func PadLeft(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(pad, length-len(s)) + s
}

// PadRight 右填充字符串
func PadRight(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(pad, length-len(s))
}

// PadCenter 居中填充字符串
func PadCenter(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}

	padLen := length - len(s)
	leftPad := padLen / 2
	rightPad := padLen - leftPad

	return strings.Repeat(pad, leftPad) + s + strings.Repeat(pad, rightPad)
}

// Truncate 截断字符串
func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}

// TruncateWithSuffix 带后缀截断字符串
func TruncateWithSuffix(s string, length int, suffix string) string {
	if len(s) <= length {
		return s
	}
	return s[:length-len(suffix)] + suffix
}

// IsEmail 检查是否为邮箱格式
func IsEmail(s string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsPhone 检查是否为手机号格式
func IsPhone(s string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsURL 检查是否为URL格式
func IsURL(s string) bool {
	pattern := `^https?://[^\s/$.?#].[^\s]*$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsNumeric 检查是否为数字
func IsNumeric(s string) bool {
	pattern := `^\d+$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsAlpha 检查是否为字母
func IsAlpha(s string) bool {
	pattern := `^[a-zA-Z]+$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsAlphaNumeric 检查是否为字母数字
func IsAlphaNumeric(s string) bool {
	pattern := `^[a-zA-Z0-9]+$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// IsUUID 检查是否为UUID格式
func IsUUID(s string) bool {
	pattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// RemoveSpaces 移除所有空格
func RemoveSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

// RemoveSpecialChars 移除特殊字符
func RemoveSpecialChars(s string) string {
	pattern := `[^a-zA-Z0-9\s]`
	reg := regexp.MustCompile(pattern)
	return reg.ReplaceAllString(s, "")
}

// NormalizeSpace 标准化空白字符
func NormalizeSpace(s string) string {
	// 将多个空白字符替换为单个空格
	pattern := `\s+`
	reg := regexp.MustCompile(pattern)
	return strings.TrimSpace(reg.ReplaceAllString(s, " "))
}

// MaskEmail 掩码邮箱
func MaskEmail(email string) string {
	if !IsEmail(email) {
		return email
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 2 {
		return email
	}

	masked := username[:1] + strings.Repeat("*", len(username)-2) + username[len(username)-1:]
	return masked + "@" + domain
}

// MaskPhone 掩码手机号
func MaskPhone(phone string) string {
	if !IsPhone(phone) {
		return phone
	}

	if len(phone) != 11 {
		return phone
	}

	return phone[:3] + "****" + phone[7:]
}
