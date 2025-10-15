# 安装指南

本指南将帮助你快速安装和配置 GinForge 框架。

## 环境要求

在开始之前，请确保你的系统满足以下要求：

### 必需组件

| 组件 | 最低版本 | 推荐版本 | 说明 |
|------|---------|---------|------|
| Go | 1.21+ | 1.22+ | 编程语言 |
| MySQL | 5.7+ | 8.0+ | 数据库 |
| Redis | 6.0+ | 7.0+ | 缓存服务 |
| Node.js | 16+ | 20+ | 前端构建 |

### 可选组件

- Docker 20+ (用于容器化部署)
- Git 2.0+ (版本控制)
- Make (构建工具)

## 安装步骤

### 1. 克隆项目

```bash
# 使用 Git 克隆
git clone https://github.com/ginforge/ginforge.git

# 进入项目目录
cd ginforge
```

### 2. 安装 Go 依赖

```bash
# 下载依赖
go mod download

# 验证依赖
go mod verify
```

### 3. 配置环境

复制配置文件模板：

```bash
# 复制环境配置
cp env.example .env

# 复制应用配置
cp configs/config.yaml.example configs/config.yaml
```

编辑 `.env` 文件，设置基本环境变量：

```bash
# 应用环境
APP_ENV=development
APP_PORT=8080

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=ginforge
DB_USERNAME=root
DB_PASSWORD=your_password

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT 密钥
JWT_SECRET=your-secret-key-change-in-production
```

### 4. 初始化数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS ginforge CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行迁移脚本
make migrate

# 或者手动执行
go run cmd/cli/main.go migrate
```

### 5. 安装前端依赖

```bash
# 进入前端目录
cd web/admin

# 安装依赖
npm install

# 返回根目录
cd ../..
```

## 验证安装

### 启动后端服务

```bash
# 启动 admin-api 服务
go run services/admin-api/cmd/server/main.go
```

你应该看到类似输出：

```
{"level":"info","ts":"2025-10-15T12:00:00.000+0700","msg":"admin-api service starting","port":8083}
```

### 启动前端服务

在新的终端窗口：

```bash
cd web/admin
npm run dev
```

访问 http://localhost:3000，你应该看到登录页面。

### 测试 API

```bash
# 测试健康检查
curl http://localhost:8083/api/v1/admin/health

# 预期输出
{
  "status": "ok",
  "service": "admin-api"
}
```

## 快速启动脚本

为了方便开发，我们提供了启动脚本：

```bash
# 启动所有服务
make start

# 或使用脚本
./scripts/start-services.sh
```

这将同时启动：
- ✅ Admin API (端口: 8083)
- ✅ User API (端口: 8081)
- ✅ WebSocket Gateway (端口: 8087)
- ✅ Gateway (端口: 8080)

## 使用 Docker

如果你更喜欢使用 Docker：

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 常见问题

### Go 依赖下载失败

如果遇到网络问题，可以设置代理：

```bash
# 使用 GOPROXY
export GOPROXY=https://goproxy.cn,direct

# 或者
export GOPROXY=https://goproxy.io,direct
```

### MySQL 连接失败

检查配置：

```bash
# 测试 MySQL 连接
mysql -h localhost -u root -p

# 检查用户权限
GRANT ALL PRIVILEGES ON ginforge.* TO 'root'@'localhost';
FLUSH PRIVILEGES;
```

### Redis 连接失败

```bash
# 检查 Redis 是否运行
redis-cli ping

# 预期输出: PONG
```

### 端口被占用

修改配置文件中的端口号：

```yaml
# configs/config.yaml
services:
  admin_api:
    port: 8083  # 改为其他端口
```

## 开发工具推荐

### IDE

- **GoLand** - JetBrains 出品的 Go IDE (推荐)
- **VS Code** - 安装 Go 扩展
- **Vim/Neovim** - 使用 vim-go 插件

### 调试工具

```bash
# 安装 Delve 调试器
go install github.com/go-delve/delve/cmd/dlv@latest

# 使用 Delve 调试
dlv debug services/admin-api/cmd/server/main.go
```

### 数据库工具

- **TablePlus** - 现代化的数据库客户端
- **DBeaver** - 开源数据库工具
- **MySQL Workbench** - 官方工具

## 下一步

安装完成后，你可以：

- [快速开始](./quick-start) - 创建你的第一个应用
- [项目结构](./project-structure) - 了解目录结构
- [配置系统](../core-concepts/configuration) - 学习如何配置

---

**遇到问题？**
- 查看 [常见问题](../faq)
- 提交 [GitHub Issue](https://github.com/ginforge/ginforge/issues)

