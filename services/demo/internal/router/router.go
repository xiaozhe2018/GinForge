package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/middleware"
	"goweb/services/demo/internal/handler"
)

func NewRouter(cfg *config.Config, log logger.Logger, h *handler.DemoHandler) *gin.Engine {
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

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "service": "demo"}) })

	// 设置处理器日志
	h.SetLogger(log)

	// API 路由
	api := r.Group("/api/v1")
	{
		api.GET("/data", h.GetData)
		api.GET("/user/:user_id", h.GetUserInfo)
	}

	return r
}
