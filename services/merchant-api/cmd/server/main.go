package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/services/merchant-api/internal/handler"
	"goweb/services/merchant-api/internal/router"
	"goweb/services/merchant-api/internal/service"

	_ "goweb/services/merchant-api/docs" // 导入生成的 docs 包
)

// @title          商户端 API
// @version         1.0
// @description     商户端相关接口文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8082
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 加载配置（新版）
	cfg := config.New()
	serviceName := "merchant-api"
	log := logger.New(serviceName, cfg.GetString("log.level"))

	// 初始化服务
	merchantService := service.NewMerchantService()
	productService := service.NewProductService()
	orderService := service.NewOrderService()

	// 初始化处理器
	merchantHandler := handler.NewMerchantHandler(merchantService, productService, orderService)

	// 初始化路由
	r := router.NewRouter(cfg, log, merchantHandler)

	// 启动HTTP服务
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.GetInt("services.merchant_api.port")),
		Handler:      r,
		ReadTimeout:  cfg.GetDuration("app.read_timeout"),
		WriteTimeout: cfg.GetDuration("app.write_timeout"),
		IdleTimeout:  cfg.GetDuration("app.idle_timeout"),
	}

	go func() {
		log.Info("merchant-api service starting", "port", cfg.GetInt("services.merchant_api.port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("merchant-api service start error", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("merchant-api service shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("merchant-api service shutdown error", err)
	}
}
