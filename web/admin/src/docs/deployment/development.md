# 开发环境

搭建 GinForge 本地开发环境的完整指南。

## 🛠️ 环境要求

### 必需工具

| 工具 | 版本要求 | 用途 |
|------|---------|------|
| Go | 1.21+ | 后端开发 |
| Node.js | 16+ | 前端开发 |
| Git | 2.0+ | 版本控制 |

### 可选工具

| 工具 | 版本 | 用途 |
|------|------|------|
| MySQL | 8.0+ | 数据库（推荐） |
| Redis | 6.0+ | 缓存和队列（推荐） |
| Docker | 20+ | 容器化开发 |

## 🚀 快速搭建

### 方式一：本地开发（推荐新手）

```bash
# 1. 克隆项目
git clone https://github.com/xiaozhe2018/GinForge.git
cd GinForge

# 2. 安装 Go 依赖
go mod tidy

# 3. 启动后端（使用 SQLite）
go run ./services/admin-api/cmd/server/main.go

# 4. 新开终端，启动前端
cd web/admin
npm install
npm run dev
```

### 方式二：使用 Docker Compose（推荐）

```bash
# 1. 启动所有服务
docker-compose up -d

# 2. 查看服务状态
docker-compose ps

# 3. 查看日志
docker-compose logs -f

# 4. 停止服务
docker-compose down
```

## 🔧 详细配置步骤

### 1. 配置 MySQL（可选）

```bash
# 使用 Docker 启动 MySQL
docker run -d \
  --name mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=gin_forge \
  -v $(pwd)/data/mysql:/var/lib/mysql \
  mysql:8.0

# 导入初始化脚本
docker exec -i mysql mysql -uroot -p123456 gin_forge < database/migrations/001_create_admin_tables.sql
docker exec -i mysql mysql -uroot -p123456 gin_forge < database/migrations/002_create_system_tables.sql
```

修改配置：

```yaml
# configs/config.yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  database: "gin_forge"
  username: "root"
  password: "123456"
```

### 2. 配置 Redis（推荐）

```bash
# 使用 Docker 启动 Redis
docker run -d \
  --name redis \
  -p 6379:6379 \
  -v $(pwd)/data/redis:/data \
  redis:7-alpine redis-server --appendonly yes

# 验证 Redis
docker exec redis redis-cli ping
# 输出：PONG
```

修改配置：

```yaml
# configs/config.yaml
redis:
  enabled: true
  host: "localhost"
  port: 6379
```

### 3. 配置环境变量

```bash
# 创建 .env 文件
cat > .env << EOF
APP_ENV=development
APP_PORT=8083
DB_TYPE=mysql
DB_PASSWORD=123456
REDIS_ENABLED=true
JWT_SECRET=dev-secret-key-$(openssl rand -hex 16)
LOG_LEVEL=debug
EOF
```

## 🔍 开发工具配置

### 1. IDE 配置

#### VS Code

安装扩展：
- Go (Official)
- Vue Language Features (Volar)
- ESLint
- Prettier

配置文件 (`.vscode/settings.json`)：

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  },
  "[vue]": {
    "editor.defaultFormatter": "Vue.volar"
  }
}
```

#### GoLand / WebStorm

1. 打开项目目录
2. 自动识别 Go 和 Vue 项目
3. 配置 GOPATH 和 Go SDK

### 2. Git Hooks

配置提交前检查：

```bash
# .git/hooks/pre-commit
#!/bin/sh

# 运行 Go lint
golangci-lint run

# 运行 Go tests
go test ./...

# 运行前端 lint
cd web/admin && npm run lint
```

### 3. 调试配置

#### Delve 调试器

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试服务
dlv debug ./services/admin-api/cmd/server/main.go

# 设置断点
(dlv) break main.main
(dlv) continue
```

#### VS Code 调试配置 (`.vscode/launch.json`)

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch admin-api",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/services/admin-api/cmd/server",
      "env": {
        "APP_ENV": "development"
      },
      "args": []
    }
  ]
}
```

## 🧪 测试环境

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/utils

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔧 开发脚本

### 创建启动脚本

```bash
#!/bin/bash
# scripts/dev.sh

# 启动 MySQL
docker start mysql || docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_DATABASE=gin_forge mysql:8.0

# 启动 Redis
docker start redis || docker run -d --name redis -p 6379:6379 redis:7-alpine

# 等待服务就绪
sleep 5

# 启动后端
go run ./services/admin-api/cmd/server/main.go &

# 启动前端
cd web/admin && npm run dev
```

使用脚本：

```bash
chmod +x scripts/dev.sh
./scripts/dev.sh
```

## 📊 日志查看

### 实时查看日志

```bash
# 查看后端日志
tail -f logs/admin-api.log

# 查看所有日志
tail -f logs/*.log

# 使用 grep 过滤
tail -f logs/admin-api.log | grep ERROR
```

### 日志级别

- `DEBUG`: 详细的调试信息
- `INFO`: 一般信息（默认）
- `WARN`: 警告信息
- `ERROR`: 错误信息
- `FATAL`: 致命错误

## 🛡️ 开发最佳实践

### 1. 使用 Air 热重载

安装 Air：

```bash
go install github.com/cosmtrek/air@latest
```

配置文件 (`.air.toml`)：

```toml
[build]
  cmd = "go build -o ./bin/admin-api ./services/admin-api/cmd/server/main.go"
  bin = "bin/admin-api"
  include_ext = ["go"]
  exclude_dir = ["web", "bin", "logs"]
```

使用：

```bash
# 在项目根目录运行
air
```

### 2. 使用 Make 命令

```makefile
# Makefile
.PHONY: dev test build clean

dev:
	go run ./services/admin-api/cmd/server/main.go

test:
	go test -v ./...

build:
	go build -o bin/admin-api ./services/admin-api/cmd/server/main.go

clean:
	rm -rf bin/ logs/*.log
```

使用：

```bash
make dev    # 启动开发服务
make test   # 运行测试
make build  # 构建
make clean  # 清理
```

## 🎯 下一步

- [生产部署](./production) - 生产环境配置
- [Docker 部署](./docker) - 容器化部署

---

**提示**: 开发环境注重便利性和调试体验，生产环境注重性能和安全。

