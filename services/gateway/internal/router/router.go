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
				"user":     "/api/v1/user",
				"merchant": "/api/v1/merchant",
				"product":  "/api/v1/product",
				"order":    "/api/v1/order",
				"admin":    "/api/v1/admin",
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
				target := cfg.GetString("external_services.user_api_url")
				if target == "" {
					target = "http://localhost:8081"
				}
				proxyToService(c, target, "/api/v1/user")
			})
		}

		// 商户端API代理
		merchant := api.Group("/merchant")
		{
			merchant.Any("/*path", func(c *gin.Context) {
				target := cfg.GetString("external_services.merchant_api_url")
				if target == "" {
					target = "http://localhost:8082"
				}
				proxyToService(c, target, "/api/v1/merchant")
			})
		}

		// 商品API代理
		product := api.Group("/product")
		{
			product.Any("/*path", func(c *gin.Context) {
				target := cfg.GetString("external_services.merchant_api_url")
				if target == "" {
					target = "http://localhost:8082"
				}
				proxyToService(c, target, "/api/v1/product")
			})
		}

		// 订单API代理
		order := api.Group("/order")
		{
			order.Any("/*path", func(c *gin.Context) {
				target := cfg.GetString("external_services.merchant_api_url")
				if target == "" {
					target = "http://localhost:8082"
				}
				proxyToService(c, target, "/api/v1/order")
			})
		}

		// 管理后台API代理
		admin := api.Group("/admin")
		{
			admin.Any("/*path", func(c *gin.Context) {
				target := cfg.GetString("external_services.admin_api_url")
				if target == "" {
					target = "http://localhost:8083"
				}
				proxyToService(c, target, "/api/v1/admin")
			})
		}
	}

	return r
}

// proxyToService 代理到后端服务
func proxyToService(c *gin.Context, targetAddr, pathPrefix string) {
	path := c.Param("path")
	if path == "" {
		path = "/"
	}

	// 确保路径以 / 开头
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	targetURL := targetAddr + pathPrefix + path
	logger := logger.New("gateway", "info")
	logger.Info("proxying request", "target", targetURL, "method", c.Request.Method)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 读取请求体
	var body io.Reader
	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("failed to read request body", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
			return
		}
		body = bytes.NewReader(bodyBytes)
	}

	// 创建代理请求
	req, err := http.NewRequest(c.Request.Method, targetURL, body)
	if err != nil {
		logger.Error("failed to create proxy request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create proxy request"})
		return
	}

	// 复制请求头
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 设置目标主机
	req.Host = c.Request.Host

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed to proxy request", err)
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

	// 设置状态码
	c.Status(resp.StatusCode)

	// 复制响应体
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		logger.Error("failed to copy response body", err)
		return
	}

	logger.Info("proxy request completed", "target", targetURL, "status", resp.StatusCode)
}
