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
	"goweb/services/merchant-api/internal/handler"
)

func NewRouter(cfg *config.Config, log logger.Logger, merchantHandler *handler.MerchantHandler) *gin.Engine {
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
		c.JSON(200, gin.H{"status": "ok", "service": "merchant-api"})
	})

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 设置处理器日志
	merchantHandler.SetLogger(log)

	// API路由
	api := r.Group("/api/v1")
	{
		// 商户相关路由
		merchant := api.Group("/merchant")
		merchant.Use(middleware.JWTAuth(cfg.GetString("jwt.secret"))) // JWT认证中间件
		{
			merchant.GET("/info", merchantHandler.GetMerchantInfo)
			merchant.PUT("/info", merchantHandler.UpdateMerchantInfo)
		}

		// 商品相关路由
		product := api.Group("/product")
		product.Use(middleware.JWTAuth(cfg.GetString("jwt.secret")))
		{
			product.GET("/list", merchantHandler.GetProducts)
			product.POST("/create", merchantHandler.CreateProduct)
		}

		// 订单相关路由
		order := api.Group("/order")
		order.Use(middleware.JWTAuth(cfg.GetString("jwt.secret")))
		{
			order.GET("/list", merchantHandler.GetOrders)
		}
	}

	return r
}
