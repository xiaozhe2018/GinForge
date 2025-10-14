package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type InitCommand struct {
	name        string
	description string
	author      string
	version     string
	output      string
}

func NewInitCommand() *InitCommand {
	return &InitCommand{}
}

func (c *InitCommand) Run(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	fs.StringVar(&c.name, "name", "", "项目名称 (必需)")
	fs.StringVar(&c.description, "description", "", "项目描述")
	fs.StringVar(&c.author, "author", "", "作者")
	fs.StringVar(&c.version, "version", "1.0.0", "版本")
	fs.StringVar(&c.output, "output", ".", "输出目录")

	fs.Parse(args)

	if c.name == "" {
		fmt.Println("错误: 项目名称不能为空")
		fmt.Println("用法: ginforge init --name=<project-name> [flags]")
		os.Exit(1)
	}

	// 生成项目
	if err := c.generateProject(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 项目 '%s' 初始化成功！\n", c.name)
	fmt.Printf("📁 目录: %s\n", c.output)
	fmt.Printf("🚀 启动: make run\n")
}

func (c *InitCommand) generateProject() error {
	projectDir := filepath.Join(c.output, c.name)

	// 创建目录结构
	dirs := []string{
		filepath.Join(projectDir, "cmd", "cli"),
		filepath.Join(projectDir, "pkg"),
		filepath.Join(projectDir, "services"),
		filepath.Join(projectDir, "configs"),
		filepath.Join(projectDir, "deployments", "docker"),
		filepath.Join(projectDir, "deployments", "k8s"),
		filepath.Join(projectDir, "docs"),
		filepath.Join(projectDir, "scripts"),
		filepath.Join(projectDir, "templates"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// 生成文件
	files := map[string]string{
		"go.mod":                         c.generateGoMod(),
		"Makefile":                       c.generateMakefile(),
		"README.md":                      c.generateReadme(),
		"configs/config.yaml":            c.generateConfigYaml(),
		"configs/config.prod.yaml":       c.generateConfigProdYaml(),
		"configs/config.test.yaml":       c.generateConfigTestYaml(),
		"env.example":                    c.generateEnvExample(),
		"cmd/cli/main.go":                c.generateCliMain(),
		"deployments/docker-compose.yml": c.generateDockerCompose(),
		"deployments/docker/Dockerfile":  c.generateDockerfile(),
		"scripts/start-redis.sh":         c.generateStartRedisScript(),
		"docs/README.md":                 c.generateDocsReadme(),
	}

	for filePath, content := range files {
		fullPath := filepath.Join(projectDir, filePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (c *InitCommand) generateGoMod() string {
	return fmt.Sprintf(`module %s

go 1.20

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-contrib/requestid v0.0.6
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/prometheus/client_golang v1.17.0
	github.com/spf13/viper v1.16.0
	github.com/stretchr/testify v1.8.4
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.16.2
	go.uber.org/zap v1.25.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/driver/postgres v1.5.2
	gorm.io/driver/sqlite v1.5.3
	gorm.io/gorm v1.25.4
)
`, c.name)
}

func (c *InitCommand) generateMakefile() string {
	return `# GinForge 微服务框架 Makefile

.PHONY: help build run stop restart status test clean swagger docker compose

# 默认目标
help:
	@echo "GinForge 微服务框架 - 可用命令:"
	@echo "  make build     - 构建所有服务"
	@echo "  make run       - 启动所有服务"
	@echo "  make stop      - 停止所有服务"
	@echo "  make restart   - 重启所有服务"
	@echo "  make status    - 查看服务状态"
	@echo "  make test      - 运行测试"
	@echo "  make clean     - 清理构建文件"
	@echo "  make swagger   - 生成 Swagger 文档"
	@echo "  make docker    - 构建 Docker 镜像"
	@echo "  make compose   - 启动 Docker Compose"

# 构建所有服务
build:
	@echo "构建所有微服务..."
	@go build -o bin/user-api ./services/user-api/cmd/server
	@go build -o bin/merchant-api ./services/merchant-api/cmd/server
	@go build -o bin/admin-api ./services/admin-api/cmd/server
	@go build -o bin/gateway ./services/gateway/cmd/server
	@echo "构建完成！"

# 启动所有服务
run:
	@echo "启动所有微服务..."
	@go run ./services/user-api/cmd/server &
	@go run ./services/merchant-api/cmd/server &
	@go run ./services/admin-api/cmd/server &
	@go run ./services/gateway/cmd/server &
	@echo "所有服务已启动！"

# 停止所有服务
stop:
	@echo "停止所有微服务..."
	@pkill -f "go run ./services/.*/cmd/server" 2>/dev/null || true
	@echo "所有微服务已停止"

# 重启所有服务
restart: stop run
	@echo "已完成重启"

# 查看服务状态
status:
	@echo "服务状态:"
	@ps aux | grep "go run ./services" | grep -v grep || echo "没有运行的服务"

# 运行测试
test:
	@echo "运行测试..."
	@go test ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf bin/
	@go clean

# 生成 Swagger 文档
swagger:
	@echo "生成 Swagger 文档..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@swag init -g services/user-api/cmd/server/main.go -o services/user-api/docs
	@swag init -g services/merchant-api/cmd/server/main.go -o services/merchant-api/docs
	@swag init -g services/admin-api/cmd/server/main.go -o services/admin-api/docs
	@swag init -g services/gateway/cmd/server/main.go -o services/gateway/docs
	@echo "Swagger 文档生成完成！"

# 构建 Docker 镜像
docker:
	@echo "构建 Docker 镜像..."
	@docker build -f deployments/docker/Dockerfile -t ginforge:latest .

# 启动 Docker Compose
compose:
	@echo "启动 Docker Compose..."
	@docker-compose -f deployments/docker-compose.yml up -d
`
}

func (c *InitCommand) generateReadme() string {
	return fmt.Sprintf(`# %s

%s

## 功能特性

- 微服务架构
- RESTful API
- 统一配置管理
- 统一日志记录
- 统一响应格式
- Swagger文档
- Docker支持
- Kubernetes支持

## 快速开始

### 环境要求

- Go 1.20+
- Docker (可选)
- Redis (可选)

### 安装依赖

`+"`"+`bash
go mod tidy
`+"`"+`

### 启动服务

`+"`"+`bash
make run
`+"`"+`

### 访问API

- 用户端API: http://localhost:8081
- 商户端API: http://localhost:8082
- 管理后台API: http://localhost:8083
- API网关: http://localhost:8080

### 生成Swagger文档

`+"`"+`bash
make swagger
`+"`"+`

## 开发指南

### 创建新服务

`+"`"+`bash
go run ./cmd/cli service --name=payment --port=8086
`+"`"+`

### 运行测试

`+"`"+`bash
make test
`+"`"+`

## 部署

### Docker

`+"`"+`bash
make docker
docker run -p 8080:8080 ginforge:latest
`+"`"+`

### Docker Compose

`+"`"+`bash
make compose
`+"`"+`

## 许可证

MIT
`, c.name, c.description)
}

func (c *InitCommand) generateConfigYaml() string {
	return `app:
  name: "GinForge"
  version: "1.0.0"
  env: "development"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

log:
  level: "info"
  format: "json"
  output: "stdout"

database:
  driver: "sqlite"
  host: "localhost"
  port: 3306
  database: "ginforge.db"
  username: ""
  password: ""
  charset: "utf8mb4"
  parse_time: true
  loc: "Local"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: "1h"
  log_level: "silent"

redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  database: 0
  pool_size: 10
  min_idle_conns: 5
  max_retries: 3
  dial_timeout: "5s"
  read_timeout: "3s"
  write_timeout: "3s"
  pool_timeout: "4s"
  idle_timeout: "5m"
  idle_check_freq: "1m"

jwt:
  secret: "dev-secret-key"
  expire: "24h"
  issuer: "GinForge"
  audience: "GinForge-Users"

services:
  user_api:
    port: 8081
  merchant_api:
    port: 8082
  admin_api:
    port: 8083
  gateway:
    port: 8080
  gateway_worker:
    port: 8084
  demo:
    port: 8085

external_services:
  user_api_url: "http://localhost:8081"
  merchant_api_url: "http://localhost:8082"
  admin_api_url: "http://localhost:8083"
  gateway_worker_url: "http://localhost:8084"
  demo_url: "http://localhost:8085"

gateway:
  base_url: "http://localhost:8080"
  timeout: "30s"
  retry_count: 3
  retry_delay: "1s"

istio:
  enabled: false
  namespace: "default"
  service_account: "default"
  sidecar_image: "docker.io/istio/proxyv2"
  sidecar_version: "1.19.0"
  proxy_cpu: "100m"
  proxy_memory: "128Mi"
  log_level: "info"
  trace_sampling: "1.0"
  access_log_format: ""
  enable_tracing: false
  enable_metrics: true
  enable_access_log: true
  enable_prometheus: true
  enable_jaeger: false
  jaeger_endpoint: ""
  prometheus_port: 15020
  zipkin_endpoint: ""
  enable_zipkin: false
`
}

func (c *InitCommand) generateConfigProdYaml() string {
	return `app:
  env: "production"
  log_level: "warn"

database:
  driver: "mysql"
  host: "mysql"
  port: 3306
  database: "ginforge_prod"
  username: "ginforge"
  password: "ginforge_password"

redis:
  host: "redis"
  port: 6379
  password: "redis_password"

jwt:
  secret: "prod-secret-key"

external_services:
  user_api_url: "http://user-api:8081"
  merchant_api_url: "http://merchant-api:8082"
  admin_api_url: "http://admin-api:8083"
  gateway_worker_url: "http://gateway-worker:8084"
  demo_url: "http://demo:8085"

istio:
  enabled: true
  namespace: "ginforge"
  service_account: "ginforge"
`
}

func (c *InitCommand) generateConfigTestYaml() string {
	return `app:
  env: "test"
  log_level: "debug"

database:
  driver: "sqlite"
  database: ":memory:"
  log_level: "silent"

redis:
  database: 15

jwt:
  secret: "test-secret-key"
  issuer: "GinForge-Test"
  audience: "GinForge-Test-Users"
`
}

func (c *InitCommand) generateEnvExample() string {
	return `# GinForge 环境变量配置示例

# 应用配置
APP_NAME=GinForge
APP_VERSION=1.0.0
APP_ENV=development
APP_PORT=8080
APP_READ_TIMEOUT=30s
APP_WRITE_TIMEOUT=30s
APP_IDLE_TIMEOUT=120s

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout

# 数据库配置
DATABASE_DRIVER=sqlite
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_DATABASE=ginforge.db
DATABASE_USERNAME=
DATABASE_PASSWORD=
DATABASE_CHARSET=utf8mb4
DATABASE_PARSE_TIME=true
DATABASE_LOC=Local
DATABASE_MAX_IDLE_CONNS=10
DATABASE_MAX_OPEN_CONNS=100
DATABASE_CONN_MAX_LIFETIME=1h
DATABASE_LOG_LEVEL=silent

# Redis配置
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=5
REDIS_MAX_RETRIES=3
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s
REDIS_POOL_TIMEOUT=4s
REDIS_IDLE_TIMEOUT=5m
REDIS_IDLE_CHECK_FREQ=1m

# JWT配置
JWT_SECRET=dev-secret-key
JWT_EXPIRE=24h
JWT_ISSUER=GinForge
JWT_AUDIENCE=GinForge-Users

# 服务端口配置
SERVICES_USER_API_PORT=8081
SERVICES_MERCHANT_API_PORT=8082
SERVICES_ADMIN_API_PORT=8083
SERVICES_GATEWAY_PORT=8080
SERVICES_GATEWAY_WORKER_PORT=8084
SERVICES_DEMO_PORT=8085

# 外部服务URL
EXTERNAL_SERVICES_USER_API_URL=http://localhost:8081
EXTERNAL_SERVICES_MERCHANT_API_URL=http://localhost:8082
EXTERNAL_SERVICES_ADMIN_API_URL=http://localhost:8083
EXTERNAL_SERVICES_GATEWAY_WORKER_URL=http://localhost:8084
EXTERNAL_SERVICES_DEMO_URL=http://localhost:8085

# 网关配置
GATEWAY_BASE_URL=http://localhost:8080
GATEWAY_TIMEOUT=30s
GATEWAY_RETRY_COUNT=3
GATEWAY_RETRY_DELAY=1s

# Istio配置
ISTIO_ENABLED=false
ISTIO_NAMESPACE=default
ISTIO_SERVICE_ACCOUNT=default
ISTIO_SIDECAR_IMAGE=docker.io/istio/proxyv2
ISTIO_SIDECAR_VERSION=1.19.0
ISTIO_PROXY_CPU=100m
ISTIO_PROXY_MEMORY=128Mi
ISTIO_LOG_LEVEL=info
ISTIO_TRACE_SAMPLING=1.0
ISTIO_ACCESS_LOG_FORMAT=
ISTIO_ENABLE_TRACING=false
ISTIO_ENABLE_METRICS=true
ISTIO_ENABLE_ACCESS_LOG=true
ISTIO_ENABLE_PROMETHEUS=true
ISTIO_ENABLE_JAEGER=false
ISTIO_JAEGER_ENDPOINT=
ISTIO_PROMETHEUS_PORT=15020
ISTIO_ZIPKIN_ENDPOINT=
ISTIO_ENABLE_ZIPKIN=false
`
}

func (c *InitCommand) generateCliMain() string {
	return `package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"goweb/cmd/cli/commands"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "service":
		commands.NewServiceCommand().Run(args)
	case "handler":
		commands.NewHandlerCommand().Run(args)
	case "model":
		commands.NewModelCommand().Run(args)
	case "middleware":
		commands.NewMiddlewareCommand().Run(args)
	case "config":
		commands.NewConfigCommand().Run(args)
	case "test":
		commands.NewTestCommand().Run(args)
	case "deploy":
		commands.NewDeployCommand().Run(args)
	case "init":
		commands.NewInitCommand().Run(args)
	case "version":
		commands.NewVersionCommand().Run(args)
	case "help", "-h", "--help":
		showHelp()
	default:
		fmt.Printf("未知命令: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("GinForge CLI - 微服务开发脚手架")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  ginforge <command> [flags]")
	fmt.Println()
	fmt.Println("可用命令:")
	fmt.Println("  service    创建新的微服务")
	fmt.Println("  handler    创建新的处理器")
	fmt.Println("  model      创建新的数据模型")
	fmt.Println("  middleware 创建新的中间件")
	fmt.Println("  config     管理配置文件")
	fmt.Println("  test       运行测试")
	fmt.Println("  deploy     部署服务")
	fmt.Println("  init       初始化项目")
	fmt.Println("  version    显示版本信息")
	fmt.Println("  help       显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  ginforge service --name=payment --port=8086")
	fmt.Println("  ginforge handler --service=user --name=profile")
	fmt.Println("  ginforge model --name=user --fields=name,email,age")
	fmt.Println("  ginforge test --service=user --coverage")
	fmt.Println("  ginforge deploy --env=production")
}
`
}

func (c *InitCommand) generateDockerCompose() string {
	return `version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: ginforge-mysql
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: ginforge
      MYSQL_USER: ginforge
      MYSQL_PASSWORD: ginforge_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - ginforge-network

  redis:
    image: redis:7-alpine
    container_name: ginforge-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - ginforge-network

  user-api:
    build: .
    container_name: ginforge-user-api
    ports:
      - "8081:8081"
    environment:
      - APP_ENV=production
      - DATABASE_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    networks:
      - ginforge-network

  merchant-api:
    build: .
    container_name: ginforge-merchant-api
    ports:
      - "8082:8082"
    environment:
      - APP_ENV=production
      - DATABASE_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    networks:
      - ginforge-network

  admin-api:
    build: .
    container_name: ginforge-admin-api
    ports:
      - "8083:8083"
    environment:
      - APP_ENV=production
      - DATABASE_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    networks:
      - ginforge-network

  gateway:
    build: .
    container_name: ginforge-gateway
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - REDIS_HOST=redis
    depends_on:
      - redis
    networks:
      - ginforge-network

volumes:
  mysql_data:
  redis_data:

networks:
  ginforge-network:
    driver: bridge
`
}

func (c *InitCommand) generateDockerfile() string {
	return `FROM golang:1.20-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 复制构建的二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/configs ./configs

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
`
}

func (c *InitCommand) generateStartRedisScript() string {
	return `#!/bin/bash

# 启动Redis容器
echo "启动Redis容器..."

docker run -d \
  --name ginforge-redis \
  -p 6379:6379 \
  -v ginforge-redis-data:/data \
  redis:7-alpine

echo "Redis容器已启动！"
echo "连接信息:"
echo "  主机: localhost"
echo "  端口: 6379"
echo "  数据库: 0"
`
}

func (c *InitCommand) generateDocsReadme() string {
	return `# GinForge 文档

欢迎使用 GinForge 微服务框架！

## 文档目录

- [快速开始](./QUICK_START.md) - 5分钟快速上手
- [框架概述](./FRAMEWORK.md) - 详细框架说明
- [API文档](./API.md) - API接口文档
- [部署指南](./DEPLOYMENT.md) - 部署相关文档
- [开发指南](./DEVELOPMENT.md) - 开发相关文档
- [测试指南](./TESTING.md) - 测试相关文档
- [故障排除](./TROUBLESHOOTING.md) - 常见问题解决

## 示例代码

- [配置使用](./demo/config.md)
- [中间件使用](./demo/middleware.md)
- [数据库使用](./demo/db.md)
- [缓存使用](./demo/cache.md)
- [路由响应](./demo/router_response.md)
- [Swagger使用](./demo/swagger.md)
- [参数校验](./demo/validation.md)
- [Redis使用](./demo/redis_usage.md)
- [延时队列](./demo/delayed_queue_usage.md)
- [网关工作器](./demo/gateway_worker_usage.md)

## 贡献指南

欢迎贡献代码和文档！

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License
`
}
