package router

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/middleware"
)

func NewRouter(cfg *config.Config, log logger.Logger) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestID())
	r.Use(middleware.AccessLogger(log))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-Id"},
		ExposeHeaders:    []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "GinForge Gateway",
			"version": "0.1.0",
			"status":  "running",
			"endpoints": gin.H{
				"health":   "/healthz",
				"api":      "/api/v1",
			},
		})
	})

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "gateway"})
	})

	// API路由
	api := r.Group("/api/v1")
	{
		// 用户端API代理
		user := api.Group("/user")
		{
			user.Any("/*path", func(c *gin.Context) {
				proxyToService(c, cfg.GetString("external_services.user_api_url"), "/api/v1/user")
			})
		}

		// 商户端API代理
		merchant := api.Group("/merchant")
		{
			merchant.Any("/*path", func(c *gin.Context) {
				proxyToService(c, cfg.GetString("external_services.merchant_api_url"), "/api/v1/merchant")
			})
		}

		// 管理后台API代理
		admin := api.Group("/admin")
		{
			admin.Any("/*path", func(c *gin.Context) {
				proxyToService(c, cfg.GetString("external_services.admin_api_url"), "/api/v1/admin")
			})
		}

		// 文件服务API代理
		files := api.Group("/files")
		{
			files.Any("/*path", func(c *gin.Context) {
				proxyToService(c, cfg.GetString("external_services.file_api_url"), "/api/v1/files")
			})
		}
	}

	return r
}

// proxyToService 代理到后端服务
func proxyToService(c *gin.Context, targetAddr, pathPrefix string) {
	if targetAddr == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service not configured"})
		return
	}

	path := c.Param("path")
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	targetURL := targetAddr + pathPrefix + path
	
	// 创建HTTP客户端
	client := &http.Client{Timeout: 30 * time.Second}

	// 读取请求体
	var body io.Reader
	if c.Request.Body != nil {
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		body = bytes.NewReader(bodyBytes)
	}

	// 创建代理请求
	req, _ := http.NewRequest(c.Request.Method, targetURL, body)

	// 复制请求头
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to proxy request"})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 复制响应
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
