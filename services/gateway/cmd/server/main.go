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
	"goweb/services/gateway/internal/router"
)

func main() {
	// 加载配置（新版）
	cfg := config.New()
	serviceName := "gateway"
	log := logger.New(serviceName, cfg.GetString("log.level"))

	// 初始化路由
	r := router.NewRouter(cfg, log)

	// 启动HTTP服务
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.GetInt("services.gateway.port")),
		Handler:      r,
		ReadTimeout:  cfg.GetDuration("app.read_timeout"),
		WriteTimeout: cfg.GetDuration("app.write_timeout"),
		IdleTimeout:  cfg.GetDuration("app.idle_timeout"),
	}

	go func() {
		log.Info("gateway service starting", "port", cfg.GetInt("services.gateway.port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("gateway service start error", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("gateway service shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("gateway service shutdown error", err)
	}
}
