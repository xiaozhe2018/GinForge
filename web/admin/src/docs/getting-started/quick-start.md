# 快速开始

只需 **5 分钟**，你就能运行起 GinForge 框架！

## 📋 环境检查

首先确保你的开发环境满足以下要求：

```bash
# 检查 Go 版本（需要 1.21+）
go version
# 输出示例：go version go1.22.0 darwin/arm64

# 检查 Node.js 版本（需要 16+）
node --version
# 输出示例：v20.10.0

# 检查 Git
git --version
```

## ⚡ 三步启动

### 步骤 1：克隆项目

```bash
# 克隆仓库
git clone https://github.com/xiaozhe2018/GinForge.git
cd GinForge

# 安装 Go 依赖
go mod tidy
```

### 步骤 2：启动后端服务

```bash
# 使用默认配置启动 admin-api（使用 SQLite 数据库）
go run ./services/admin-api/cmd/server/main.go
```

**看到以下日志表示启动成功：**

```json
{"level":"info","ts":"2025-10-15T12:00:00+0700","msg":"admin-api service starting","port":8083}
```

服务启动后，你可以访问：
- 🏥 健康检查：http://localhost:8083/api/v1/admin/system/health
- 📚 API 文档：http://localhost:8083/swagger/index.html

### 步骤 3：启动前端服务

打开新的终端窗口：

```bash
# 进入前端目录
cd web/admin

# 安装依赖（首次需要）
npm install

# 启动开发服务器
npm run dev
```

**看到以下信息表示成功：**

```
VITE v5.x ready in xxx ms
➜  Local:   http://localhost:3000/
```

## 🎉 开始使用

### 1. 访问管理后台

打开浏览器访问：**http://localhost:3000**

### 2. 使用默认账号登录

```
用户名：admin
密码：admin123
```

### 3. 探索功能

登录后你可以看到：

- 📊 **仪表盘** - 系统概况和统计
- 👥 **后台用户管理** - 管理员账号管理
- 🎭 **角色管理** - 角色和权限分配
- 📑 **菜单管理** - 配置系统菜单
- 🔒 **权限管理** - 细粒度权限控制
- ⚙️ **系统管理** - 系统配置和监控
- 📖 **文档中心** - 在线文档（你正在看的就是）

## 🔧 配置数据库（可选）

### 选项 A：使用 SQLite（默认）

无需任何配置，开箱即用！数据存储在 `goweb.db` 文件中。

### 选项 B：使用 MySQL

如果你想使用 MySQL：

```bash
# 1. 启动 MySQL（使用 Docker）
docker run -d \
  --name mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=gin_forge \
  mysql:8.0

# 2. 导入初始化脚本
docker exec -i mysql mysql -uroot -p123456 gin_forge < database/migrations/001_create_admin_tables.sql
docker exec -i mysql mysql -uroot -p123456 gin_forge < database/migrations/002_create_system_tables.sql
```

修改配置文件 `configs/config.yaml`：

```yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  database: "gin_forge"
  username: "root"
  password: "123456"
```

## 🚀 启动 Redis（推荐）

Redis 用于缓存、消息队列、WebSocket 等功能：

```bash
# 使用 Docker 启动 Redis
docker run -d \
  --name redis \
  -p 6379:6379 \
  redis:7-alpine

# 验证 Redis 运行
docker exec redis redis-cli ping
# 输出：PONG
```

修改配置 `configs/config.yaml`：

```yaml
redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

## 📝 测试 API

### 使用 curl 测试

```bash
# 测试登录 API
curl -X POST http://localhost:8083/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

**成功响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGc...",
    "user": {
      "id": 1,
      "username": "admin",
      "name": "管理员"
    },
    "menus": [...],
    "permissions": [...]
  }
}
```

### 使用 Swagger 测试

访问 **http://localhost:8083/swagger/index.html** 查看完整的 API 文档并在线测试。

## 🎯 下一步

恭喜！你已经成功运行了 GinForge 框架。接下来：

1. [了解项目结构](./project-structure) - 熟悉框架的目录组织
2. [学习配置系统](../core-concepts/configuration) - 掌握配置管理
3. [创建第一个 API](../core-concepts/routing) - 开始开发

## ⚠️ 常见问题

### 端口被占用

```bash
# 查看占用端口的进程
lsof -ti :8083

# 杀死进程
kill -9 $(lsof -ti :8083)
```

### 数据库连接失败

```bash
# 检查 MySQL 是否运行
docker ps | grep mysql

# 查看 MySQL 日志
docker logs mysql
```

### 前端无法访问后端

检查 CORS 配置和前端代理配置：

```typescript
// web/admin/vite.config.ts
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8083',
      changeOrigin: true
    }
  }
}
```

---

**遇到问题？** 查看 [完整的故障排查指南](../../GETTING_STARTED.md#常见问题与解决方案)

