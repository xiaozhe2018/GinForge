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
	"goweb/services/file-api/internal/handler"
	"goweb/services/file-api/internal/model"
	"goweb/services/file-api/internal/repository"
	"goweb/services/file-api/internal/router"
	"goweb/services/file-api/internal/service"
)

// @title 文件上传服务 API
// @version 1.0
// @description GinForge 文件上传微服务 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8086
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 加载配置
	cfg := config.New()
	serviceName := "file-api"
	log := logger.New(serviceName, cfg.GetString("log.level"))

	log.Info("starting file-api service")

	// 初始化数据库
	database, err := db.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize database", err)
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(
		&model.FileRecord{},
		&model.FileUploadLog{},
		&model.FileDownloadLog{},
	); err != nil {
		log.Warn("failed to auto migrate database", "error", err)
	}

	// 获取存储配置，设置默认值
	storageType := cfg.GetString("storage.type")
	if storageType == "" {
		storageType = "local"
	}

	basePath := cfg.GetString("storage.local.base_path")
	if basePath == "" {
		basePath = "./uploads"
	}

	urlPrefix := cfg.GetString("storage.url_prefix")
	if urlPrefix == "" {
		urlPrefix = "http://localhost:8086/uploads"
	}

	maxFileSize := cfg.GetInt64("storage.max_file_size")
	if maxFileSize == 0 {
		maxFileSize = 104857600 // 100MB
	}

	// 初始化存储服务
	storageConfig := &service.StorageConfig{
		Type:          service.StorageType(storageType),
		LocalBasePath: basePath,
		URLPrefix:     urlPrefix,
		MaxFileSize:   maxFileSize,
	}

	storageService, err := service.NewStorageService(storageConfig, log)
	if err != nil {
		log.Fatal("failed to initialize storage service", err)
	}

	// 初始化仓储层
	fileRepo := repository.NewFileRepository(database)

	// 初始化服务层
	fileService := service.NewFileService(fileRepo, storageService, database, log)

	// 初始化处理器
	fileHandler := handler.NewFileHandler(fileService, cfg, log)

	// 初始化路由
	r := router.NewRouter(cfg, log, fileHandler)

	// 启动HTTP服务 - 文件服务固定使用8086端口
	port := 8086

	// 允许通过环境变量覆盖端口
	if envPort := os.Getenv("FILE_API_PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", &port)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  60 * time.Second, // 文件上传需要更长的超时时间
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info("file-api service starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("file-api service start error", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("file-api service shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("file-api service shutdown error", err)
	}

	log.Info("file-api service stopped")
}
