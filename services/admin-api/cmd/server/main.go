// @title GinForge Admin API
// @version 1.0
// @description 管理后台API接口文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8083
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT认证，格式: Bearer {token}

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
	"goweb/pkg/db"
	"goweb/pkg/logger"
	"goweb/pkg/notification"
	"goweb/pkg/redis"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/router"
)

func main() {
	// 加载配置（新版）
	cfg := config.New()
	serviceName := "admin-api"
	log := logger.New(serviceName, cfg.GetString("log.level"))

	log.Info("starting admin-api service")

	// 初始化数据库
	database, err := db.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize database", err)
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(
		&model.AdminUser{},
		&model.AdminRole{},
		&model.AdminPermission{},
		&model.AdminMenu{},
		&model.AdminUserRole{},
		&model.AdminRolePermission{},
		&model.AdminRoleMenu{},
		&model.AdminOperationLog{},
		&model.AdminSystemConfig{},
	); err != nil {
		log.Warn("failed to auto migrate database", "error", err)
	}

	// 初始化 Redis 客户端
	var redisClient *redis.Client
	redisConfig := cfg.GetRedisConfig()
	if redisConfig.Enabled {
		redisClient = redis.NewClient(&redisConfig, log)
		log.Info("redis client initialized successfully")
	}

	// 初始化通知服务
	var notifyService *notification.Service
	// 创建通知配置
	notificationConfig := &notification.Config{
		Email: nil, // 可以从配置文件加载
		SMS:   nil, // 可以从配置文件加载
	}
	
	notifyService, err = notification.NewServiceFromConfig(notificationConfig, log)
	if err != nil {
		log.Warn("failed to initialize notification service", "error", err)
	} else {
		log.Info("notification service initialized successfully")
	}

	// 初始化路由
	r := router.NewRouter(database, redisClient, notifyService, log, cfg)

	// 启动HTTP服务
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.GetInt("services.admin_api.port")),
		Handler:      r,
		ReadTimeout:  cfg.GetDuration("app.read_timeout"),
		WriteTimeout: cfg.GetDuration("app.write_timeout"),
		IdleTimeout:  cfg.GetDuration("app.idle_timeout"),
	}

	go func() {
		log.Info("admin-api service starting", "port", cfg.GetInt("services.admin_api.port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("admin-api service start error", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("admin-api service shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("admin-api service shutdown error", err)
	}
}
