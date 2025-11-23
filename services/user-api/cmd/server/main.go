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
	"goweb/services/user-api/internal/handler"
	"goweb/services/user-api/internal/router"
	"goweb/services/user-api/internal/service"

	_ "goweb/services/user-api/docs" // 导入生成的 docs 包
)

// @title          用户端 API
// @version         1.0
// @description     用户端相关接口文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 加载配置（新版）
	cfg := config.New()
	serviceName := "user-api"
	log := logger.New(
		serviceName,
		cfg.GetString("log.level"),
		cfg.GetString("log.output"),
		cfg.GetString("log.dir"),
	)

	// 初始化服务
	userService := service.NewUserService()
	userHandler := handler.NewUserHandler(userService)

	// 初始化路由
	r := router.NewRouter(cfg, log, userHandler)

	// 启动HTTP服务
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.GetInt("services.user_api.port")),
		Handler:      r,
		ReadTimeout:  cfg.GetDuration("app.read_timeout"),
		WriteTimeout: cfg.GetDuration("app.write_timeout"),
		IdleTimeout:  cfg.GetDuration("app.idle_timeout"),
	}

	go func() {
		log.Info("user-api service starting", "port", cfg.GetInt("services.user_api.port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("user-api service start error", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("user-api service shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("user-api service shutdown error", err)
	}
}
