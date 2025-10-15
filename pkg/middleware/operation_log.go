package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OperationLog 操作日志结构体
type OperationLog struct {
	UserID       *uint64     `json:"user_id"`
	Username     *string     `json:"username"`
	Method       string      `json:"method"`
	Path         string      `json:"path"`
	IP           *string     `json:"ip"`
	UserAgent    *string     `json:"user_agent"`
	RequestData  *string     `json:"request_data"`
	ResponseData *string     `json:"response_data"`
	StatusCode   int         `json:"status_code"`
	Duration     int         `json:"duration"`
	CreatedAt    time.Time   `json:"created_at"`
}

// OperationLogWriter 操作日志响应写入器
type OperationLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写Write方法，保存响应内容
func (w *OperationLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// OperationLogger 操作日志中间件
func OperationLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要记录的路径
		if shouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 记录开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 替换响应写入器
		blw := &OperationLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer([]byte{}),
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		duration := int(time.Since(startTime).Milliseconds())

		// 获取用户信息
		var userID *uint64
		var username *string
		if user, exists := c.Get("user"); exists {
			if u, ok := user.(map[string]interface{}); ok {
				if id, ok := u["id"].(uint64); ok {
					userID = &id
				}
				if name, ok := u["username"].(string); ok {
					username = &name
				}
			}
		}

		// 获取IP地址
		ip := c.ClientIP()

		// 获取用户代理
		userAgent := c.Request.UserAgent()

		// 处理请求数据
		var requestData *string
		if len(requestBody) > 0 {
			// 敏感字段过滤
			filteredReqBody := filterSensitiveData(string(requestBody))
			requestData = &filteredReqBody
		}

		// 处理响应数据
		var responseData *string
		if blw.body.Len() > 0 {
			// 敏感字段过滤
			filteredRespBody := filterSensitiveData(blw.body.String())
			responseData = &filteredRespBody
		}

		// 创建操作日志
		operationLog := OperationLog{
			UserID:       userID,
			Username:     username,
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			IP:           &ip,
			UserAgent:    &userAgent,
			RequestData:  requestData,
			ResponseData: responseData,
			StatusCode:   c.Writer.Status(),
			Duration:     duration,
			CreatedAt:    time.Now(),
		}

		// 异步保存日志
		go saveOperationLog(db, operationLog)
	}
}

// shouldSkipPath 判断是否跳过记录
func shouldSkipPath(path string) bool {
	// 跳过静态资源、健康检查等路径
	skipPaths := []string{
		"/healthz",
		"/metrics",
		"/favicon.ico",
		"/static",
		"/uploads",
	}

	for _, p := range skipPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	return false
}

// filterSensitiveData 过滤敏感数据
func filterSensitiveData(data string) string {
	// 尝试解析JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		// 不是JSON或解析失败，直接返回
		return data
	}

	// 敏感字段列表
	sensitiveFields := []string{
		"password", "pwd", "secret", "token", "access_token", "refresh_token",
		"credit_card", "card_number", "cvv", "ssn", "private_key",
	}

	// 替换敏感字段
	for _, field := range sensitiveFields {
		if _, exists := jsonData[field]; exists {
			jsonData[field] = "******"
		}
	}

	// 重新序列化
	filteredData, err := json.Marshal(jsonData)
	if err != nil {
		return data
	}

	return string(filteredData)
}

// saveOperationLog 保存操作日志
func saveOperationLog(db *gorm.DB, log OperationLog) {
	if db != nil {
		db.Table("admin_operation_logs").Create(&log)
	}
}