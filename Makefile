# GinForge 微服务框架 Makefile

.PHONY: help build run stop restart status test clean swagger docker compose init generate-config db-init db-reset db-status

# 默认目标
help:
	@echo "GinForge 微服务框架 - 可用命令:"
	@echo "  make build     - 构建所有服务"
	@echo "  make run       - 启动所有服务"
	@echo "  make stop      - 停止所有以 go run 启动的服务"
	@echo "  make restart   - 停止并重新启动所有服务"
	@echo "  make status    - 查看端口占用与服务状态"
	@echo "  make test      - 运行测试"
	@echo "  make clean     - 清理构建文件"
	@echo "  make swagger   - 生成 Swagger 文档"
	@echo "  make docker    - 构建 Docker 镜像"
	@echo "  make compose   - 启动 Docker Compose"
	@echo ""
	@echo "数据库命令:"
	@echo "  make init           - 初始化数据库（推荐：生成配置文件 + 初始化数据库）"
	@echo "  make generate-config - 从 .env 生成 configs/config.yaml"
	@echo "  make db-init        - 初始化数据库（执行所有迁移文件）"
	@echo "  make db-reset       - 重置数据库（删除并重新创建）"
	@echo "  make db-status      - 查看数据库状态"
	@echo ""
	@echo "配置说明:"
	@echo "  1. 复制 env.example 为 .env: cp env.example .env"
	@echo "  2. 修改 .env 文件中的配置"
	@echo "  3. 运行 make init 生成配置文件并初始化数据库"

# 构建所有服务
build:
	@echo "构建所有微服务..."
	@go build -o bin/user-api ./services/user-api/cmd/server
	@go build -o bin/merchant-api ./services/merchant-api/cmd/server
	@go build -o bin/admin-api ./services/admin-api/cmd/server
	@go build -o bin/gateway ./services/gateway/cmd/server
	@go build -o bin/websocket-gateway ./services/websocket-gateway/cmd/server
	@go build -o bin/gateway-worker ./services/gateway-worker/cmd/server
	@go build -o bin/demo ./services/demo/cmd/server
	@go build -o bin/file-api ./services/file-api/cmd/server
	@go build -o bin/ginforge ./cmd/cli
	@echo "构建完成！"

# 安装CLI工具
install-cli:
	@echo "安装CLI工具..."
	@go build -o bin/ginforge ./cmd/cli
	@sudo cp bin/ginforge /usr/local/bin/
	@echo "CLI工具已安装到 /usr/local/bin/ginforge"

# 启动所有服务
run:
	@echo "启动所有微服务..."
	@go run ./services/user-api/cmd/server &
	@go run ./services/merchant-api/cmd/server &
	@go run ./services/admin-api/cmd/server &
	@go run ./services/gateway/cmd/server &
	@go run ./services/websocket-gateway/cmd/server &
	@go run ./services/gateway-worker/cmd/server &
	@go run ./services/demo/cmd/server &
	@go run ./services/file-api/cmd/server &
	@echo "所有服务已启动！"
	@echo "API网关: http://localhost:8080"
	@echo "用户端API: http://localhost:8081"
	@echo "商户端API: http://localhost:8082"
	@echo "管理后台API: http://localhost:8083"
	@echo "网关工作器: http://localhost:8084"
	@echo "文件服务API: http://localhost:8086"
	@echo "WebSocket网关: ws://localhost:8087"

# 需要清理/检查的端口（开发环境）
PORTS=8080 8081 8082 8083 8084 8085 8086 8087

# 停止所有以 go run 启动的微服务（开发环境）
stop:
	@echo "停止所有微服务..."
	@-pkill -f "go run ./services/.*/cmd/server" 2>/dev/null || true
	@-pkill -f "services/.*/cmd/server" 2>/dev/null || true
	@for p in $(PORTS); do \
		pid=$$(lsof -ti :$$p 2>/dev/null || true); \
		if [ -n "$$pid" ]; then \
			echo "杀死端口 $$p 的进程: $$pid"; \
			kill -9 $$pid 2>/dev/null || true; \
		fi; \
	done
	@echo "所有微服务已停止（若仍有端口占用请手动 lsof -i :<port> 检查）。"

# 重启所有微服务
restart: stop run
	@echo "已完成重启。"

# 查看端口占用与服务状态
status:
	@for p in $(PORTS); do \
		echo "端口: $$p"; \
		lsof -i :$$p | cat || true; \
		echo; \
	done

# 运行测试
test:
	@echo "运行测试..."
	@go test ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 运行基准测试
benchmark:
	@echo "运行基准测试..."
	@go test -bench=. ./...

# 运行集成测试
test-integration:
	@echo "运行集成测试..."
	@go test -tags=integration ./...

# 运行CLI测试
test-cli:
	@echo "运行CLI测试..."
	@go test ./cmd/cli/...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf bin/
	@rm -rf services/*/docs/
	@go clean

# 生成 Swagger 文档
swagger:
	@echo "生成 Swagger 文档..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@swag init -g services/user-api/cmd/server/main.go -o services/user-api/docs
	@swag init -g services/merchant-api/cmd/server/main.go -o services/merchant-api/docs
	@swag init -g services/admin-api/cmd/server/main.go -o services/admin-api/docs
	@swag init -g services/gateway/cmd/server/main.go -o services/gateway/docs
	@swag init -g services/file-api/cmd/server/main.go -o services/file-api/docs
	@echo "Swagger 文档生成完成！"
	@echo "访问文档:"
	@echo "  用户端: http://localhost:8081/swagger/index.html"
	@echo "  商户端: http://localhost:8082/swagger/index.html"
	@echo "  管理后台: http://localhost:8083/swagger/index.html"
	@echo "  文件服务: http://localhost:8086/swagger/index.html"

# 构建 Docker 镜像
docker:
	@echo "构建 Docker 镜像..."
	@docker build -f deployments/docker/Dockerfile -t goease:latest .

# 启动 Docker Compose
compose:
	@echo "启动 Docker Compose..."
	@docker-compose -f deployments/docker-compose.yml up -d

# 停止 Docker Compose
compose-down:
	@echo "停止 Docker Compose..."
	@docker-compose -f deployments/docker-compose.yml down

# 前端开发命令
web-dev:
	@echo "启动前端开发服务器..."
	@cd web/admin && npm run dev

web-build:
	@echo "构建前端项目..."
	@cd web/admin && npm run build

web-install:
	@echo "安装前端依赖..."
	@cd web/admin && npm install

# 开发环境快速启动
dev: swagger run
	@echo "开发环境已启动！"
	@echo "Swagger 文档已生成，服务已启动"
# 完整开发环境（后端+前端）
dev-full: swagger run web-dev
	@echo "完整开发环境已启动！"
	@echo "后端服务: http://localhost:8080"
	@echo "前端管理后台: http://localhost:3000"

# ============================================
# 数据库管理命令
# ============================================

# 生成配置文件（从 .env 生成 configs/config.yaml）
generate-config:
	@echo "生成配置文件..."
	@if [ ! -f .env ]; then \
		echo "错误: .env 文件不存在"; \
		echo ""; \
		echo "请先复制 env.example 为 .env 并修改配置:"; \
		echo "  cp env.example .env"; \
		echo "  vim .env"; \
		exit 1; \
	fi
	@go run scripts/generate-config.go

# 初始化数据库（别名，推荐使用）
# 步骤：1. 生成配置文件  2. 初始化数据库
init: generate-config db-init
	@echo ""
	@echo "✓ 数据库初始化完成！"

# 从配置文件读取数据库信息
DB_HOST := $(shell grep -A 10 "^database:" configs/config.yaml | grep "^  host:" | awk '{print $$2}' | tr -d '"')
DB_PORT := $(shell grep -A 10 "^database:" configs/config.yaml | grep "^  port:" | awk '{print $$2}')
DB_NAME := $(shell grep -A 10 "^database:" configs/config.yaml | grep "^  database:" | awk '{print $$2}' | tr -d '"')
DB_USER := $(shell grep -A 10 "^database:" configs/config.yaml | grep "^  username:" | awk '{print $$2}' | tr -d '"')
DB_PASS := $(shell grep -A 10 "^database:" configs/config.yaml | grep "^  password:" | awk '{print $$2}' | tr -d '"')

# 初始化数据库（执行所有迁移文件）
db-init:
	@echo "开始初始化数据库..."
	@if [ -z "$(DB_HOST)" ] || [ -z "$(DB_NAME)" ] || [ -z "$(DB_USER)" ]; then \
		echo "错误: 无法读取数据库配置"; \
		echo ""; \
		echo "请确保已运行 'make generate-config' 或 'make init' 生成配置文件"; \
		echo "配置文件路径: configs/config.yaml"; \
		echo ""; \
		echo "如果配置文件不存在，请:"; \
		echo "  1. 复制 env.example 为 .env: cp env.example .env"; \
		echo "  2. 修改 .env 文件中的数据库配置"; \
		echo "  3. 运行 make init 生成配置文件并初始化数据库"; \
		exit 1; \
	fi
	@echo "数据库: $(DB_NAME)"
	@echo "主机: $(DB_HOST):$(DB_PORT)"
	@echo "用户: $(DB_USER)"
	@echo ""
	@echo "检测 MySQL 连接方式..."
	@MYSQL_CONTAINER=$$(docker ps --format "{{.Names}}" | grep -i mysql | head -n1); \
	if [ -n "$$MYSQL_CONTAINER" ]; then \
		echo "检测到 Docker MySQL 容器: $$MYSQL_CONTAINER"; \
		echo "使用 Docker 方式执行..."; \
		echo "创建数据库（如果不存在）..."; \
		docker exec $$MYSQL_CONTAINER mysql -u$(DB_USER) -p$(DB_PASS) -e "CREATE DATABASE IF NOT EXISTS $(DB_NAME) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null || \
		(echo "无法创建数据库，请检查 MySQL 连接"; exit 1); \
		echo "执行数据库初始化脚本..."; \
		if [ -f "database/migrations/init.sql" ]; then \
			echo "  → init.sql"; \
			docker exec -i $$MYSQL_CONTAINER mysql -u$(DB_USER) -p$(DB_PASS) $(DB_NAME) < database/migrations/init.sql 2>/dev/null || \
			(echo "     执行失败，请检查 SQL 文件"; exit 1); \
		else \
			echo "     初始化脚本不存在: database/migrations/init.sql"; \
			exit 1; \
		fi; \
	else \
		echo "使用本地 MySQL 客户端..."; \
		echo "创建数据库（如果不存在）..."; \
		MYSQL_PWD=$(DB_PASS) mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) -e "CREATE DATABASE IF NOT EXISTS $(DB_NAME) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null || \
		(echo "无法创建数据库，请检查 MySQL 连接"; exit 1); \
		echo "执行数据库初始化脚本..."; \
		if [ -f "database/migrations/init.sql" ]; then \
			echo "  → init.sql"; \
			MYSQL_PWD=$(DB_PASS) mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) $(DB_NAME) < database/migrations/init.sql 2>/dev/null || \
			(echo "     执行失败，请检查 SQL 文件"; exit 1); \
		else \
			echo "     初始化脚本不存在: database/migrations/init.sql"; \
			exit 1; \
		fi; \
	fi
	@echo ""
	@echo "数据库初始化完成！"
	@echo "默认管理员账号: admin / admin123"

# 重置数据库（删除并重新创建）
db-reset:
	@if [ -z "$(DB_NAME)" ]; then \
		echo " 无法从配置文件读取数据库名称"; \
		exit 1; \
	fi
	@echo "警告: 这将删除数据库 $(DB_NAME) 并重新创建！"
	@echo -n "确认继续? (y/N): " && read confirm && [ "$$confirm" = "y" ] || exit 1
	@echo "删除数据库..."
	@MYSQL_PWD=$(DB_PASS) mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) -e "DROP DATABASE IF EXISTS $(DB_NAME);" 2>/dev/null || true
	@echo "创建数据库..."
	@MYSQL_PWD=$(DB_PASS) mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) -e "CREATE DATABASE $(DB_NAME) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null || true
	@echo "执行迁移文件..."
	@$(MAKE) db-init

# 查看数据库状态
db-status:
	@echo "数据库状态:"
	@echo "  主机: $(DB_HOST):$(DB_PORT)"
	@echo "  数据库: $(DB_NAME)"
	@echo "  用户: $(DB_USER)"
	@echo ""
	@echo "数据库表列表:"
	@MYSQL_PWD=$(DB_PASS) mysql -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) $(DB_NAME) -e "SHOW TABLES;" 2>/dev/null || echo "无法连接到数据库，请检查配置"

