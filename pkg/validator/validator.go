package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator 验证器
type Validator struct {
	validator *validator.Validate
}

// NewValidator 创建验证器
func NewValidator() *Validator {
	v := validator.New()

	// 注册自定义验证规则
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("idcard", validateIDCard)
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("username", validateUsername)

	// 注册自定义标签名函数
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validator: v,
	}
}

// Validate 验证结构体
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// ValidateVar 验证单个变量
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validator.Var(field, tag)
}

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) != 11 {
		return false
	}

	// 简单的手机号验证
	return strings.HasPrefix(phone, "1") && len(phone) == 11
}

// validateIDCard 验证身份证号
func validateIDCard(fl validator.FieldLevel) bool {
	idcard := fl.Field().String()
	if len(idcard) != 18 {
		return false
	}

	// 简单的身份证号验证
	return len(idcard) == 18
}

// validatePassword 验证密码
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 6 {
		return false
	}

	// 密码必须包含字母和数字
	hasLetter := false
	hasDigit := false

	for _, char := range password {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	// 用户名只能包含字母、数字、下划线
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}

	return true
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// Error 实现 error 接口
func (ve ValidationError) Error() string {
	return ve.Message
}

// ValidationErrors 验证错误列表
type ValidationErrors []ValidationError

// Error 实现 error 接口
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// GetValidationErrors 获取验证错误
func (v *Validator) GetValidationErrors(err error) ValidationErrors {
	var validationErrors ValidationErrors

	if err == nil {
		return validationErrors
	}

	// 检查是否为验证错误
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErr {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   err.Param(),
				Message: getValidationMessage(err),
			})
		}
	}

	return validationErrors
}

// getValidationMessage 获取验证错误消息
func getValidationMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return field + " is required"
	case "min":
		return field + " must be at least " + param + " characters"
	case "max":
		return field + " must be at most " + param + " characters"
	case "email":
		return field + " must be a valid email address"
	case "len":
		return field + " must be exactly " + param + " characters"
	case "gte":
		return field + " must be greater than or equal to " + param
	case "lte":
		return field + " must be less than or equal to " + param
	case "gt":
		return field + " must be greater than " + param
	case "lt":
		return field + " must be less than " + param
	case "oneof":
		return field + " must be one of: " + param
	case "numeric":
		return field + " must be numeric"
	case "alpha":
		return field + " must contain only letters"
	case "alphanum":
		return field + " must contain only letters and numbers"
	case "url":
		return field + " must be a valid URL"
	case "datetime":
		return field + " must be a valid datetime"
	case "unique":
		return field + " must be unique"
	case "phone":
		return field + " must be a valid phone number"
	case "idcard":
		return field + " must be a valid ID card number"
	case "password":
		return field + " must be a valid password (at least 6 characters, contains letters and numbers)"
	case "username":
		return field + " must be a valid username (3-20 characters, letters, numbers, underscore only)"
	default:
		return field + " is invalid"
	}
}

// ValidateStruct 验证结构体
func ValidateStruct(obj interface{}) error {
	v := NewValidator()
	return v.Validate(obj)
}

// ValidateVar 验证单个变量
func ValidateVar(field interface{}, tag string) error {
	v := NewValidator()
	return v.ValidateVar(field, tag)
}

// GetValidationErrors 获取验证错误
func GetValidationErrors(err error) ValidationErrors {
	v := NewValidator()
	return v.GetValidationErrors(err)
}

// IsValid 检查是否有效
func IsValid(obj interface{}) bool {
	return ValidateStruct(obj) == nil
}

// IsValidVar 检查变量是否有效
func IsValidVar(field interface{}, tag string) bool {
	return ValidateVar(field, tag) == nil
}

// ValidateAndGetErrors 验证并获取错误
func ValidateAndGetErrors(obj interface{}) ValidationErrors {
	err := ValidateStruct(obj)
	return GetValidationErrors(err)
}

// ValidateAndGetFirstError 验证并获取第一个错误
func ValidateAndGetFirstError(obj interface{}) error {
	errors := ValidateAndGetErrors(obj)
	if len(errors) == 0 {
		return nil
	}
	return errors[0]
}

// ValidateAndGetFirstErrorMessage 验证并获取第一个错误消息
func ValidateAndGetFirstErrorMessage(obj interface{}) string {
	err := ValidateAndGetFirstError(obj)
	if err == nil {
		return ""
	}
	return err.Error()
}

// ValidateAndGetErrorMessages 验证并获取所有错误消息
func ValidateAndGetErrorMessages(obj interface{}) []string {
	errors := ValidateAndGetErrors(obj)
	messages := make([]string, len(errors))
	for i, err := range errors {
		messages[i] = err.Message
	}
	return messages
}

// ValidateAndGetErrorMap 验证并获取错误映射
func ValidateAndGetErrorMap(obj interface{}) map[string]string {
	errors := ValidateAndGetErrors(obj)
	errorMap := make(map[string]string)
	for _, err := range errors {
		errorMap[err.Field] = err.Message
	}
	return errorMap
}

// ValidateAndGetErrorSummary 验证并获取错误摘要
func ValidateAndGetErrorSummary(obj interface{}) string {
	errors := ValidateAndGetErrors(obj)
	if len(errors) == 0 {
		return ""
	}

	var messages []string
	for _, err := range errors {
		messages = append(messages, err.Message)
	}

	return fmt.Sprintf("Validation failed: %s", strings.Join(messages, "; "))
}
