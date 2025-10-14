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
	"goweb/pkg/gateway"
	"goweb/pkg/logger"
	"goweb/services/demo/internal/handler"
	"goweb/services/demo/internal/router"
	"goweb/services/demo/internal/service"
)

func main() {
	// 加载配置（新版）
	cfg := config.New()
	serviceName := "demo"
	log := logger.New(serviceName, cfg.GetString("log.level"))

	// 初始化 Gateway 客户端
	gatewayClient := gateway.NewClient(cfg, log)

	// 初始化服务
	demoService := service.NewDemoService(gatewayClient, log)
	demoHandler := handler.NewDemoHandler(demoService, gatewayClient, log)

	// 初始化路由
	r := router.NewRouter(cfg, log, demoHandler)

	// 启动HTTP服务
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.GetInt("services.demo.port")),
		Handler:      r,
		ReadTimeout:  cfg.GetDuration("app.read_timeout"),
		WriteTimeout: cfg.GetDuration("app.write_timeout"),
		IdleTimeout:  cfg.GetDuration("app.idle_timeout"),
	}

	go func() {
		log.Info("demo service starting", "port", cfg.GetInt("services.demo.port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("demo service start error", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("demo service shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("demo service shutdown error", err)
	}
}
