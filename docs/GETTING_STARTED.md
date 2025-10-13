# 🚀 GinForge 快速上手指南

> 从零到运行，只需 5 分钟！本指南将带你快速启动 GinForge 微服务框架。

## 📋 前置准备

### 检查开发环境

```bash
# 检查 Go 版本（需要 1.20+）
go version

# 检查 Node.js 版本（需要 16+）
node --version
npm --version

# 检查 Git
git --version
```

如果缺少任何工具，请先安装：
- **Go**: https://golang.org/dl/
- **Node.js**: https://nodejs.org/
- **Git**: https://git-scm.com/

### 可选组件

以下组件是可选的，框架会根据配置自动适配：

```bash
# MySQL（推荐生产环境使用）
# 开发环境可以使用 SQLite（无需安装）
mysql --version

# Redis（推荐，用于缓存和队列）
redis-cli --version

# Docker（用于容器化部署）
docker --version
docker-compose --version
```

## ⚡ 快速开始（三步启动）

### 第一步：克隆项目

```bash
# 克隆仓库
git clone https://github.com/xiaozhe2018/GinForge.git
cd GinForge

# 安装 Go 依赖
go mod tidy
```

### 第二步：启动后端服务

```bash
# 方式1：直接运行（使用默认配置和 SQLite）
go run ./services/admin-api/cmd/server

# 方式2：使用 Make 命令
make run

# 方式3：使用环境变量自定义配置
export APP_PORT=8083
export DB_TYPE=sqlite
go run ./services/admin-api/cmd/server
```

**看到以下日志表示启动成功：**
```
✅ Server started at http://localhost:8083
✅ Swagger docs at http://localhost:8083/swagger/index.html
```

### 第三步：启动前端服务

```bash
# 打开新终端，进入前端目录
cd web/admin

# 安装依赖（首次运行需要）
npm install

# 启动开发服务器
npm run dev
```

**看到以下信息表示成功：**
```
  VITE v5.0.8  ready in 500 ms

  ➜  Local:   http://localhost:3000/
  ➜  Network: use --host to expose
  ➜  press h to show help
```

### 🎉 开始使用

1. **打开浏览器**，访问：http://localhost:3000
2. **使用默认账号登录**：
   - 用户名：`admin`
   - 密码：`admin123`
3. **探索功能**：
   - ✅ 用户管理
   - ✅ 角色管理
   - ✅ 菜单管理
   - ✅ 权限管理
   - ✅ 个人设置

## 🔧 环境配置详解

### 1️⃣ 数据库配置

#### 选项 A：使用 SQLite（推荐新手）

默认配置，无需额外安装，开箱即用。数据存储在 `goweb.db` 文件中。

```yaml
# configs/config.yaml
database:
  type: "sqlite"
  database: "goweb.db"
```

#### 选项 B：使用 MySQL（推荐生产）

```bash
# 1. 启动 MySQL（使用 Docker）
docker run -d \
  --name mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=gin_forge \
  mysql:8.0

# 2. 等待 MySQL 启动（约 30 秒）
docker logs -f mysql

# 3. 导入初始化 SQL
docker exec -i mysql mysql -uroot -p123456 gin_forge < database/migrations/001_create_admin_tables.sql
```

**修改配置文件**：
```yaml
# configs/config.yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  database: "gin_forge"
  username: "root"
  password: "123456"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
```

### 2️⃣ Redis 配置（可选）

Redis 用于缓存、Token 黑名单、消息队列等功能。

```bash
# 启动 Redis（使用 Docker）
docker run -d \
  --name redis \
  -p 6379:6379 \
  redis:7-alpine

# 验证 Redis 运行
docker exec redis redis-cli ping
# 输出：PONG
```

**修改配置文件**：
```yaml
# configs/config.yaml
redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
```

### 3️⃣ 环境变量配置

创建 `.env` 文件（基于 `env.example`）：

```bash
# 复制示例文件
cp env.example .env

# 编辑配置
vim .env
```

`.env` 文件示例：
```bash
# 应用配置
APP_ENV=development
APP_PORT=8083
APP_DEBUG=true

# 数据库配置
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=gin_forge
DB_USERNAME=root
DB_PASSWORD=123456

# Redis 配置
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT 配置
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRE_HOURS=24
```

**配置优先级**：环境变量 > .env 文件 > config.yaml > 默认值

## 🎯 验证安装

### 检查后端服务

```bash
# 1. 测试健康检查
curl http://localhost:8083/health

# 预期输出：
# {"status":"ok","database":"connected","redis":"connected"}

# 2. 查看 Swagger 文档
curl http://localhost:8083/swagger/doc.json

# 3. 测试登录 API
curl -X POST http://localhost:8083/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 预期输出：
# {
#   "code": 0,
#   "message": "success",
#   "data": {
#     "token": "eyJhbGc...",
#     "user": {...},
#     "menus": [...],
#     "permissions": [...]
#   }
# }
```

### 检查前端服务

1. 访问 http://localhost:3000
2. 应该看到登录页面
3. 输入账号密码后能成功登录
4. 浏览器控制台无错误信息

## 🎯 核心功能

登录后你可以：

- ✅ **仪表盘** - 查看系统概况和统计信息
- ✅ **用户管理** - 管理系统用户，分配角色
- ✅ **角色管理** - 配置角色和权限
- ✅ **菜单管理** - 配置系统菜单结构
- ✅ **权限管理** - 细粒度权限控制
- ✅ **系统管理** - 系统配置和监控
- ✅ **个人设置** - 修改个人信息和密码

## 📚 下一步

### 查看API文档
访问Swagger文档: **http://localhost:8083/swagger/index.html**

### 测试API
```bash
# 测试登录API
curl -X POST http://localhost:8083/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 阅读完整文档
- 项目总结: `PROJECT_STATUS.md`
- API对接: `web/admin/API_INTEGRATION.md`
- 框架文档: `docs/FRAMEWORK.md`
- 文档索引: `docs/INDEX.md`

## 🛠️ 开发新功能

### 创建新服务
```bash
go run ./cmd/generator -command=service -name=payment
```

### 生成API文档
```bash
make swagger
```

### 运行测试
```bash
make test
```

## ⚠️ 常见问题与解决方案

### 问题 1：后端启动失败

#### 症状：`port already in use`
```
Error: listen tcp :8083: bind: address already in use
```

**解决方案**：
```bash
# 查看占用端口的进程
lsof -ti :8083

# 杀死进程
kill -9 $(lsof -ti :8083)

# 或者修改端口
export APP_PORT=8084
go run ./services/admin-api/cmd/server
```

#### 症状：`database connection failed`
```
Error: failed to connect to database
```

**解决方案**：
```bash
# 1. 检查 MySQL 是否运行
docker ps | grep mysql

# 2. 如果使用 SQLite，确保有写权限
ls -la goweb.db
chmod 666 goweb.db

# 3. 测试数据库连接
mysql -h localhost -u root -p123456 gin_forge
```

#### 症状：`go: module not found`
```
Error: package goweb/pkg/config not found
```

**解决方案**：
```bash
# 重新下载依赖
go mod tidy
go mod download

# 清理缓存
go clean -modcache
go mod tidy
```

### 问题 2：前端启动失败

#### 症状：`npm install` 失败
```
npm ERR! code ERESOLVE
```

**解决方案**：
```bash
# 方案 1：清理缓存重新安装
cd web/admin
rm -rf node_modules package-lock.json
npm cache clean --force
npm install

# 方案 2：使用淘宝镜像
npm config set registry https://registry.npmmirror.com
npm install

# 方案 3：使用 cnpm
npm install -g cnpm --registry=https://registry.npmmirror.com
cnpm install
```

#### 症状：TypeScript 编译错误
```
Type error: Cannot find module 'element-plus'
```

**解决方案**：
```bash
# 重新安装类型定义
cd web/admin
npm install --save-dev @types/node
npm install

# 或者删除 node_modules 重新安装
rm -rf node_modules package-lock.json
npm install
```

#### 症状：`Vite` 启动失败
```
Error: Cannot find module 'vite'
```

**解决方案**：
```bash
# 确保在正确的目录
cd web/admin
pwd  # 应该显示 .../GinForge/web/admin

# 重新安装
npm install
```

### 问题 3：登录相关问题

#### 症状：登录提示"用户名或密码错误"
**原因**：可能是数据库未初始化或密码错误

**解决方案**：
```bash
# 1. 确认使用正确的账号密码
用户名：admin
密码：admin123

# 2. 检查数据库中的用户
# SQLite
sqlite3 goweb.db "SELECT username, status FROM admin_users;"

# MySQL
docker exec mysql mysql -uroot -p123456 gin_forge \
  -e "SELECT username, status FROM admin_users;"

# 3. 重置管理员密码（如果需要）
# 在后端代码中临时添加重置逻辑或重新导入初始化 SQL
```

#### 症状：登录后立即退出
**原因**：Token 存储或验证问题

**解决方案**：
```bash
# 1. 打开浏览器控制台（F12）查看错误
# 2. 检查 localStorage
localStorage.getItem('token')
localStorage.getItem('user')

# 3. 清除浏览器缓存
# Chrome: Ctrl+Shift+Delete
# 或者在控制台执行
localStorage.clear()

# 4. 检查后端日志
# 看是否有 JWT 验证错误
```

#### 症状：Token 过期太快
**原因**：JWT 过期时间配置问题

**解决方案**：
```yaml
# 修改 configs/config.yaml
jwt:
  secret: "your-secret-key"
  expire_hours: 24  # 修改为 24 小时或更长
```

### 问题 4：API 请求失败

#### 症状：`404 Not Found`
```
GET http://localhost:8083/api/v1/admin/users 404 (Not Found)
```

**解决方案**：
```bash
# 1. 确认后端服务运行正常
curl http://localhost:8083/health

# 2. 检查路由配置
# 查看 services/admin-api/internal/router/router.go

# 3. 查看 Swagger 文档确认正确的 API 路径
# http://localhost:8083/swagger/index.html
```

#### 症状：`CORS` 跨域错误
```
Access to XMLHttpRequest has been blocked by CORS policy
```

**解决方案**：
```go
// 检查 CORS 中间件配置
// pkg/middleware/cors.go 或路由配置中
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

#### 症状：`401 Unauthorized`
```
{"code":401,"message":"未授权访问"}
```

**解决方案**：
```bash
# 1. 检查 Token 是否正确
# 在浏览器控制台查看请求头
# Network -> 选择请求 -> Headers -> Authorization

# 2. 重新登录获取新 Token

# 3. 检查 Redis 黑名单（如果使用 Redis）
docker exec redis redis-cli KEYS "token:blacklist:*"
```

### 问题 5：Docker 部署问题

#### 症状：Docker 容器启动失败
```
Error response from daemon: driver failed programming external connectivity
```

**解决方案**：
```bash
# 1. 检查端口占用
lsof -i :8083
lsof -i :3306
lsof -i :6379

# 2. 重启 Docker
# Mac
restart Docker Desktop

# Linux
sudo systemctl restart docker

# 3. 清理未使用的容器和网络
docker system prune
```

### 问题 6：性能问题

#### 症状：API 响应慢
**解决方案**：
```bash
# 1. 启用 Redis 缓存
# 修改 configs/config.yaml
redis:
  enabled: true

# 2. 优化数据库查询
# 检查慢查询日志
# 添加适当的索引

# 3. 增加数据库连接池
database:
  max_open_conns: 100
  max_idle_conns: 10

# 4. 检查系统资源
top
htop
```

## 📚 下一步学习

### 1. 阅读核心文档
- [📖 框架使用指南](./docs/FRAMEWORK.md) - 详细的框架使用说明
- [⚡ 快速开始](./docs/QUICK_START.md) - 5分钟快速入门
- [🚀 高级功能](./docs/ADVANCED_FEATURES.md) - 高级功能详解

### 2. 查看示例代码
- [基类使用](./docs/demo/base_classes_usage.md)
- [缓存使用](./docs/demo/cache.md)
- [队列使用](./docs/demo/queue_usage.md)
- [Redis 使用](./docs/demo/redis_usage.md)

### 3. 开发新功能
```bash
# 创建新服务
go run ./cmd/generator -command=service -name=payment

# 生成 API 文档
make swagger

# 运行测试
make test
```

### 4. 生产部署
- [Docker 部署](./deployments/docker/)
- [Kubernetes 部署](./deployments/k8s/)
- [Istio 服务网格](./deployments/k8s/istio/)

## 📞 获取帮助

### 文档资源
- **完整文档索引**: [docs/INDEX.md](./docs/INDEX.md)
- **API 文档**: http://localhost:8083/swagger/index.html
- **前端对接文档**: [web/admin/API_INTEGRATION.md](./web/admin/API_INTEGRATION.md)
- **故障排查**: [web/admin/TROUBLESHOOTING.md](./web/admin/TROUBLESHOOTING.md)

### 示例代码
- **基础示例**: [docs/demo/](./docs/demo/)
- **测试用例**: `pkg/*/test` 和 `services/*/test`

### 社区支持
- **GitHub Issues**: 提交 Bug 或功能请求
- **GitHub Discussions**: 技术讨论和问答
- **项目 Wiki**: 更多教程和最佳实践

## 🎓 学习路径建议

### 初学者（第 1-3 天）
1. ✅ 完成快速开始，运行起来
2. ✅ 浏览管理后台，熟悉功能
3. ✅ 查看 Swagger 文档，了解 API
4. ✅ 阅读框架核心文档

### 进阶开发（第 4-7 天）
1. 📝 创建第一个自定义服务
2. 📝 编写 API 接口和业务逻辑
3. 📝 集成 Redis 缓存
4. 📝 添加单元测试

### 高级应用（第 2-4 周）
1. 🚀 使用消息队列处理异步任务
2. 🚀 实现分布式锁
3. 🚀 配置熔断和限流
4. 🚀 Docker/K8s 部署

---

**GinForge - 让开发更加简单** 🚀

**遇到问题？** 查看 [故障排查文档](./web/admin/TROUBLESHOOTING.md) 或提交 [Issue](https://github.com/xiaozhe2018/GinForge/issues)

