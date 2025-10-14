package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ServiceCommand struct {
	name        string
	port        int
	description string
	author      string
	version     string
	output      string
}

func NewServiceCommand() *ServiceCommand {
	return &ServiceCommand{}
}

func (c *ServiceCommand) Run(args []string) {
	fs := flag.NewFlagSet("service", flag.ExitOnError)
	fs.StringVar(&c.name, "name", "", "服务名称 (必需)")
	fs.IntVar(&c.port, "port", 0, "服务端口")
	fs.StringVar(&c.description, "description", "", "服务描述")
	fs.StringVar(&c.author, "author", "", "作者")
	fs.StringVar(&c.version, "version", "1.0.0", "版本")
	fs.StringVar(&c.output, "output", "services", "输出目录")

	fs.Parse(args)

	if c.name == "" {
		fmt.Println("错误: 服务名称不能为空")
		fmt.Println("用法: ginforge service --name=<service-name> [flags]")
		os.Exit(1)
	}

	// 生成端口号
	if c.port == 0 {
		c.port = c.generatePort()
	}

	// 生成服务
	if err := c.generateService(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 服务 '%s' 创建成功！\n", c.name)
	fmt.Printf("📁 目录: %s/%s\n", c.output, c.name)
	fmt.Printf("🌐 端口: %d\n", c.port)
	fmt.Printf("🚀 启动: go run ./%s/%s/cmd/server\n", c.output, c.name)
}

func (c *ServiceCommand) generatePort() int {
	// 简单的端口分配策略
	basePort := 8080
	serviceCount := c.countExistingServices()
	return basePort + serviceCount + 1
}

func (c *ServiceCommand) countExistingServices() int {
	servicesDir := "services"
	if _, err := os.Stat(servicesDir); os.IsNotExist(err) {
		return 0
	}

	files, err := os.ReadDir(servicesDir)
	if err != nil {
		return 0
	}

	count := 0
	for _, file := range files {
		if file.IsDir() && strings.HasSuffix(file.Name(), "-api") {
			count++
		}
	}
	return count
}

func (c *ServiceCommand) generateService() error {
	serviceDir := filepath.Join(c.output, c.name)

	// 创建目录结构
	dirs := []string{
		filepath.Join(serviceDir, "cmd", "server"),
		filepath.Join(serviceDir, "internal", "handler"),
		filepath.Join(serviceDir, "internal", "service"),
		filepath.Join(serviceDir, "internal", "router"),
		filepath.Join(serviceDir, "internal", "model"),
		filepath.Join(serviceDir, "docs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// 生成文件
	files := map[string]string{
		"cmd/server/main.go":          c.generateMain(),
		"internal/handler/handler.go": c.generateHandler(),
		"internal/service/service.go": c.generateServiceTemplate(),
		"internal/router/router.go":   c.generateRouter(),
		"internal/model/model.go":     c.generateModel(),
		"README.md":                   c.generateReadme(),
		"go.mod":                      c.generateGoMod(),
	}

	for filePath, content := range files {
		fullPath := filepath.Join(serviceDir, filePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (c *ServiceCommand) generateMain() string {
	return fmt.Sprintf(`package main

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
	"goweb/%s/internal/handler"
	"goweb/%s/internal/router"
	"goweb/%s/internal/service"
)

func main() {
	// 加载配置
	cfg := config.New()
	serviceName := "%s"
	log := logger.New(serviceName, cfg.GetString("log.level"))

	// 初始化服务
	%sService := service.New%sService()

	// 初始化处理器
	%sHandler := handler.New%sHandler(%sService)

	// 初始化路由
	r := router.NewRouter(cfg, log, %sHandler)

	// 启动HTTP服务
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", %d),
		Handler:      r,
		ReadTimeout:  cfg.GetDuration("app.read_timeout"),
		WriteTimeout: cfg.GetDuration("app.write_timeout"),
		IdleTimeout:  cfg.GetDuration("app.idle_timeout"),
	}

	go func() {
		log.Info("%s service starting", "port", %d)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("%s service start error", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("%s service shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("%s service shutdown error", err)
	}
}
`, c.name, c.name, c.name, c.name, c.name, c.name, c.name, c.name, c.name, c.name, c.port, c.name, c.port, c.name, c.name, c.name)
}

func (c *ServiceCommand) generateHandler() string {
	return fmt.Sprintf(`package handler

import (
	"github.com/gin-gonic/gin"
	"goweb/pkg/response"
	"goweb/%s/internal/service"
)

type %sHandler struct {
	svc *service.%sService
}

func New%sHandler(svc *service.%sService) *%sHandler {
	return &%sHandler{svc: svc}
}

// GetData 获取数据
// @Summary 获取数据
// @Description 获取示例数据
// @Tags %s
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=object}
// @Router /data [get]
func (h *%sHandler) GetData(c *gin.Context) {
	data, err := h.svc.GetData()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// CreateData 创建数据
// @Summary 创建数据
// @Description 创建示例数据
// @Tags %s
// @Accept json
// @Produce json
// @Param data body object true "数据"
// @Success 200 {object} response.Response{data=object}
// @Router /data [post]
func (h *%sHandler) CreateData(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := h.svc.CreateData(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}
`, c.name, c.name, c.name, c.name, c.name, c.name, c.name, c.name, c.name, c.name, c.name, c.name)
}

func (c *ServiceCommand) generateServiceTemplate() string {
	return fmt.Sprintf(`package service

import (
	"goweb/pkg/base"
)

type %sService struct {
	*base.BaseService
}

func New%sService() *%sService {
	return &%sService{
		BaseService: base.NewBaseService("%s"),
	}
}

// GetData 获取数据
func (s *%sService) GetData() (interface{}, error) {
	s.LogInfo("获取数据")
	return map[string]interface{}{
		"message": "Hello from %s service",
		"version": "1.0.0",
	}, nil
}

// CreateData 创建数据
func (s *%sService) CreateData(data map[string]interface{}) (interface{}, error) {
	s.LogInfo("创建数据", "data", data)
	return map[string]interface{}{
		"id":      1,
		"message": "数据创建成功",
		"data":    data,
	}, nil
}
`, c.name, c.name, c.name, c.name, c.name, c.name, c.name, c.name)
}

func (c *ServiceCommand) generateRouter() string {
	return fmt.Sprintf(`package router

import (
	"github.com/gin-gonic/gin"
	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/middleware"
	"goweb/pkg/swagger"
	"goweb/%s/internal/handler"
)

func NewRouter(cfg *config.Config, log *logger.Logger, %sHandler *handler.%sHandler) *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestID())
	r.Use(middleware.AccessLogger(log))
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger文档
	if !cfg.IsProduction() {
		swagger.SetupSwagger(r, "/swagger")
	}

	// API路由
	api := r.Group("/api/v1")
	{
		api.GET("/data", %sHandler.GetData)
		api.POST("/data", %sHandler.CreateData)
	}

	return r
}
`, c.name, c.name, c.name, c.name, c.name)
}

func (c *ServiceCommand) generateModel() string {
	return fmt.Sprintf(`package model

import (
	"time"
	"goweb/pkg/model"
)

// %s %s数据模型
type %s struct {
	model.BaseModel
	Name        string    `+"`json:\"name\" gorm:\"size:100;not null\"`"+`
	Description string    `+"`json:\"description\" gorm:\"size:500\"`"+`
	Status      int       `+"`json:\"status\" gorm:\"default:1\"`"+`
	CreatedAt   time.Time `+"`json:\"created_at\"`"+`
	UpdatedAt   time.Time `+"`json:\"updated_at\"`"+`
}

// TableName 表名
func (%s) TableName() string {
	return "%s"
}
`, c.name, c.name, c.name, c.name, c.name, c.name)
}

func (c *ServiceCommand) generateReadme() string {
	return fmt.Sprintf(`# %s Service

%s

## 功能特性

- RESTful API
- 健康检查
- Swagger文档
- 统一日志
- 统一响应

## 快速开始

### 启动服务

`+"`"+`bash
go run ./cmd/server
`+"`"+`

### 访问API

- 健康检查: http://localhost:%d/healthz
- API文档: http://localhost:%d/swagger/index.html
- 示例API: http://localhost:%d/api/v1/data

## API接口

### GET /api/v1/data
获取数据

### POST /api/v1/data
创建数据

## 开发指南

### 添加新的API

1. 在 internal/service 中添加业务逻辑
2. 在 internal/handler 中添加HTTP处理器
3. 在 internal/router 中注册路由
4. 添加Swagger注解

### 测试

`+"`"+`bash
go test ./...
`+"`"+`

## 部署

### Docker

`+"`"+`bash
docker build -t %s:latest .
docker run -p %d:%d %s:latest
`+"`"+`
`, c.name, c.description, c.port, c.port, c.port, c.name, c.port, c.port, c.name)
}

func (c *ServiceCommand) generateGoMod() string {
	return fmt.Sprintf(`module goweb/%s

go 1.20

require (
	github.com/gin-gonic/gin v1.9.1
	goweb v0.0.0
)

replace goweb => ../../
`)
}
