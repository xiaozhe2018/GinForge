package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/middleware"
	"goweb/services/user-api/internal/handler"
)

func NewRouter(cfg *config.Config, log logger.Logger, userHandler *handler.UserHandler) *gin.Engine {
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
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "user-api"})
	})

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 设置处理器日志
	userHandler.SetLogger(log)

	// API路由
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		user := api.Group("/user")
		user.Use(middleware.JWTAuth(cfg.GetString("jwt.secret"))) // JWT认证中间件
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.GET("/orders", userHandler.GetOrders)
		}
	}

	return r
}
