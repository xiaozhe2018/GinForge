package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
		TraceID: c.GetString("request_id"),
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		TraceID: c.GetString("request_id"),
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// InternalError 500错误
func InternalError(c *gin.Context, message string) {
	Error(c, 500, message)
}
