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
	"goweb/pkg/redis"
	"goweb/services/gateway-worker/internal/handler"
	"goweb/services/gateway-worker/internal/service"
)

func main() {
	// 加载配置
	cfg := config.New()
	serviceName := "gateway-worker"
	log := logger.New(
		serviceName,
		cfg.GetString("log.level"),
		cfg.GetString("log.output"),
		cfg.GetString("log.dir"),
	)

	// 创建 Redis 管理器
	redisConfig := cfg.GetRedisConfig()
	redisManager := redis.NewManager(&redisConfig, log)

	// 检查 Redis 连接
	ctx := context.Background()
	if err := redisManager.Ping(ctx); err != nil {
		log.Fatal("Redis connection failed", err)
	}

	// 初始化服务
	workerService := service.NewWorkerService(redisManager, log)
	workerHandler := handler.NewWorkerHandler(workerService, log)

	// 启动消息消费者
	if err := workerService.StartConsumers(ctx); err != nil {
		log.Fatal("Failed to start message consumers", err)
	}

	// 启动健康检查服务
	healthServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.GetInt("services.gateway_worker.port")),
		Handler: workerHandler.GetHealthHandler(),
	}

	go func() {
		log.Info("gateway-worker health server starting", "port", cfg.GetInt("services.gateway_worker.port"))
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("health server start error", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("gateway-worker shutting down...")

	// 停止消息消费者
	workerService.StopConsumers()

	// 关闭健康检查服务
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := healthServer.Shutdown(shutdownCtx); err != nil {
		log.Error("health server shutdown error", err)
	}

	log.Info("gateway-worker stopped")
}
