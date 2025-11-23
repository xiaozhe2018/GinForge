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
	"goweb/pkg/websocket"
	"goweb/pkg/websocket/group"
	"goweb/pkg/websocket/session"
	"goweb/services/websocket-gateway/internal/router"
)

func main() {
	// 加载配置
	cfg := config.New()
	serviceName := "websocket-gateway"
	log := logger.New(
		serviceName,
		cfg.GetString("log.level"),
		cfg.GetString("log.output"),
		cfg.GetString("log.dir"),
	)

	// 初始化 Redis 客户端（如果启用）
	var redisClient *redis.Client
	redisConfig := cfg.GetRedisConfig()
	if redisConfig.Enabled {
		redisClient = redis.NewClient(&redisConfig, log)
		log.Info("Redis 客户端初始化成功")
	} else {
		log.Warn("Redis 未启用，将使用内存存储")
	}

	// 创建 WebSocket 管理器
	wsManager := websocket.NewManager(log)

	// 启动 WebSocket 管理器
	go wsManager.Run()

	// 创建会话管理器
	var sessionMgr *session.SessionManager
	if redisClient != nil {
		// 使用 Redis 会话管理
		sessionMgr = session.NewRedisSessionManager(redisClient.GetClient(), "ws", 24*time.Hour)
		log.Info("使用 Redis 会话管理")
	} else {
		// 使用内存会话管理
		sessionMgr = session.NewMemorySessionManager()
		log.Info("使用内存会话管理")
	}

	// 创建分组管理器
	var groupMgr *group.GroupManager
	if redisClient != nil {
		// 使用 Redis 分组管理
		groupMgr = group.NewRedisGroupManager(redisClient.GetClient(), "ws", 24*time.Hour)
		log.Info("使用 Redis 分组管理")
	} else {
		// 使用内存分组管理
		groupMgr = group.NewMemoryGroupManager()
		log.Info("使用内存分组管理")
	}

	// 初始化路由
	r := router.NewRouter(cfg, log, wsManager, sessionMgr, groupMgr, redisClient)

	// 启动 HTTP 服务
	port := cfg.GetInt("services.websocket_gateway.port")
	if port == 0 {
		port = 8087 // 默认端口
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  cfg.GetDuration("app.read_timeout"),
		WriteTimeout: cfg.GetDuration("app.write_timeout"),
		IdleTimeout:  120 * time.Second, // WebSocket 需要更长的空闲超时
	}

	go func() {
		log.Info("websocket-gateway service starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("websocket-gateway service start error", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("websocket-gateway service shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("websocket-gateway service shutdown error", err)
	}

	log.Info("websocket-gateway service stopped")
}
