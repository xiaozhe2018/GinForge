package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// PasswordComplexity 密码复杂度要求
type PasswordComplexity struct {
	RequireUppercase bool // 需要大写字母
	RequireLowercase bool // 需要小写字母
	RequireNumbers   bool // 需要数字
	RequireSymbols   bool // 需要特殊字符
}

// ParseComplexity 从字符串数组解析密码复杂度要求
func ParseComplexity(requirements []string) PasswordComplexity {
	complexity := PasswordComplexity{}
	
	for _, req := range requirements {
		switch strings.ToLower(req) {
		case "uppercase":
			complexity.RequireUppercase = true
		case "lowercase":
			complexity.RequireLowercase = true
		case "numbers":
			complexity.RequireNumbers = true
		case "symbols", "special":
			complexity.RequireSymbols = true
		}
	}
	
	return complexity
}

// ValidatePassword 验证密码
func ValidatePassword(password string, minLength int, complexity PasswordComplexity) error {
	// 检查长度
	if len(password) < minLength {
		return fmt.Errorf("密码长度至少为 %d 位", minLength)
	}
	
	// 检查是否包含大写字母
	if complexity.RequireUppercase {
		hasUpper := false
		for _, ch := range password {
			if unicode.IsUpper(ch) {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			return errors.New("密码必须包含大写字母")
		}
	}
	
	// 检查是否包含小写字母
	if complexity.RequireLowercase {
		hasLower := false
		for _, ch := range password {
			if unicode.IsLower(ch) {
				hasLower = true
				break
			}
		}
		if !hasLower {
			return errors.New("密码必须包含小写字母")
		}
	}
	
	// 检查是否包含数字
	if complexity.RequireNumbers {
		hasNumber := false
		for _, ch := range password {
			if unicode.IsDigit(ch) {
				hasNumber = true
				break
			}
		}
		if !hasNumber {
			return errors.New("密码必须包含数字")
		}
	}
	
	// 检查是否包含特殊字符
	if complexity.RequireSymbols {
		symbolPattern := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
		if !symbolPattern.MatchString(password) {
			return errors.New("密码必须包含特殊字符")
		}
	}
	
	return nil
}

// GetPasswordRequirements 获取密码要求描述
func GetPasswordRequirements(minLength int, complexity PasswordComplexity) string {
	requirements := []string{
		fmt.Sprintf("至少 %d 位", minLength),
	}
	
	if complexity.RequireLowercase {
		requirements = append(requirements, "包含小写字母")
	}
	if complexity.RequireUppercase {
		requirements = append(requirements, "包含大写字母")
	}
	if complexity.RequireNumbers {
		requirements = append(requirements, "包含数字")
	}
	if complexity.RequireSymbols {
		requirements = append(requirements, "包含特殊字符")
	}
	
	return strings.Join(requirements, "、")
}

