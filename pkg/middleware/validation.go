package middleware

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Validator 验证器
var Validator = validator.New()

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
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

// BindAndValidate 绑定并验证请求体
func BindAndValidate(c *gin.Context, obj interface{}) error {
	// 绑定JSON
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}

	// 验证
	if err := ValidateStruct(obj); err != nil {
		return err
	}

	return nil
}

// ValidateStruct 验证结构体
func ValidateStruct(obj interface{}) error {
	err := Validator.Struct(obj)
	if err == nil {
		return nil
	}

	var validationErrors ValidationErrors
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   err.Field(),
			Tag:     err.Tag(),
			Value:   err.Param(),
			Message: getValidationMessage(err),
		})
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
	default:
		return field + " is invalid"
	}
}

// ValidationMiddleware 验证中间件
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// ValidateQuery 验证查询参数
func ValidateQuery(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindQuery(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid query parameters",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		if err := ValidateStruct(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Validation failed",
				"errors":  err,
			})
			c.Abort()
			return
		}

		c.Set("validated_query", obj)
		c.Next()
	}
}

// ValidateJSON 验证JSON请求体
func ValidateJSON(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := BindAndValidate(c, obj); err != nil {
			if validationErrors, ok := err.(ValidationErrors); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "Validation failed",
					"errors":  validationErrors,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "Invalid request body",
					"error":   err.Error(),
				})
			}
			c.Abort()
			return
		}

		c.Set("validated_body", obj)
		c.Next()
	}
}

// ValidateForm 验证表单数据
func ValidateForm(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBind(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid form data",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		if err := ValidateStruct(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Validation failed",
				"errors":  err,
			})
			c.Abort()
			return
		}

		c.Set("validated_form", obj)
		c.Next()
	}
}

// ValidateParams 验证路径参数
func ValidateParams(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建结构体实例
		val := reflect.New(reflect.TypeOf(obj).Elem()).Interface()

		// 绑定路径参数
		if err := c.ShouldBindUri(val); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid path parameters",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 验证
		if err := ValidateStruct(val); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Validation failed",
				"errors":  err,
			})
			c.Abort()
			return
		}

		c.Set("validated_params", val)
		c.Next()
	}
}
