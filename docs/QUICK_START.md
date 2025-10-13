# GinForge 快速入门指南

## 🚀 5分钟快速开始

### 1. 环境准备

```bash
# 确保 Go 1.20+ 已安装
go version

# 克隆项目
git clone <your-repo> && cd goweb

# 安装依赖
go mod tidy
```

### 2. 启动 Redis（可选）

```bash
# 使用 Docker 启动 Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 或者使用本地 Redis
redis-server
```

### 3. 启动服务

```bash
# 方式1：使用 Makefile（推荐）
make run

# 方式2：手动启动
go run ./services/user-api/cmd/server &
go run ./services/merchant-api/cmd/server &
go run ./services/admin-api/cmd/server &
go run ./services/gateway/cmd/server &
go run ./services/gateway-worker/cmd/server &
```

### 4. 验证服务

```bash
# 检查 API 网关
curl http://localhost:8080/healthz

# 检查用户端 API
curl http://localhost:8081/healthz

# 检查商户端 API
curl http://localhost:8082/healthz

# 检查管理后台 API
curl http://localhost:8083/healthz

# 检查 Gateway Worker
curl http://localhost:8084/healthz
```

## 📚 学习路径

### 新手路径（30分钟）

1. **了解框架**（5分钟）
   - 阅读 [FRAMEWORK.md](./FRAMEWORK.md) 的"介绍"部分
   - 了解框架特性和架构

2. **配置系统**（10分钟）
   - 查看 [demo/config.md](./demo/config.md)
   - 了解如何配置服务

3. **第一个 API**（15分钟）
   - 查看 [demo/router_response.md](./demo/router_response.md)
   - 学习如何创建 API 接口

### 进阶路径（1小时）

1. **中间件使用**（15分钟）
   - 查看 [demo/middleware.md](./demo/middleware.md)
   - 学习中间件的使用

2. **数据库操作**（15分钟）
   - 查看 [demo/db.md](./demo/db.md)
   - 学习数据库的使用

3. **缓存系统**（15分钟）
   - 查看 [demo/cache.md](./demo/cache.md)
   - 学习缓存的使用

4. **Swagger 文档**（15分钟）
   - 查看 [demo/swagger.md](./demo/swagger.md)
   - 学习 API 文档生成

### 高级路径（2小时）

1. **消息队列**（30分钟）
   - 查看 [demo/queue_usage.md](./demo/queue_usage.md)
   - 学习消息队列的使用

2. **延时队列**（30分钟）
   - 查看 [demo/delayed_queue_usage.md](./demo/delayed_queue_usage.md)
   - 学习延时消息的处理

3. **Gateway Worker**（30分钟）
   - 查看 [demo/gateway_worker_usage.md](./demo/gateway_worker_usage.md)
   - 学习工作服务的部署

4. **高级功能**（30分钟）
   - 查看 [demo/advanced_features.md](./demo/advanced_features.md)
   - 学习监控、文件存储等功能

## 🛠️ 常用命令

### 开发命令

```bash
# 构建所有服务
make build

# 启动所有服务
make run

# 运行测试
make test

# 清理构建文件
make clean

# 生成 Swagger 文档
make swagger
```

### 服务管理

```bash
# 启动单个服务
go run ./services/user-api/cmd/server

# 构建单个服务
go build -o bin/user-api ./services/user-api/cmd/server

# 查看服务日志
tail -f logs/app.log
```

## 🔧 配置说明

### 基础配置

```yaml
# configs/config.yaml
app:
  name: "GinForge Framework"
  version: "0.1.0"
  env: "development"
  port: 8080

# 服务端口配置
services:
  user_api: 8081
  merchant_api: 8082
  admin_api: 8083
  gateway: 8080
  gateway_worker: 8084
  demo: 8085
```

### 环境变量

```bash
# .env 文件
APP_PORT=8080
APP_ENV=development
LOG_LEVEL=debug
REDIS_HOST=localhost
REDIS_PORT=6379
```

## 📖 示例代码

### 创建新服务

```bash
# 使用生成器创建新服务
go run ./cmd/generator -command=service -name=my-service

# 启动新服务
go run ./services/my-service/cmd/server
```

### 发送消息

```go
// 发送普通消息
queue.Publish(ctx, "user.notification", map[string]interface{}{
    "user_id": "123",
    "message": "欢迎使用 GinForge！",
})

// 发送延时消息
queue.PublishWithDelay(ctx, "order.reminder", map[string]interface{}{
    "order_id": "456",
    "message": "请及时完成支付",
}, 24*time.Hour)
```

### 使用缓存

```go
// 设置缓存
cache.Set(ctx, "user:123", userData, 5*time.Minute)

// 获取缓存
var userData User
cache.Get(ctx, "user:123", &userData)
```

## 🐛 常见问题

### Q: 服务启动失败？
A: 检查端口是否被占用，查看日志错误信息

### Q: Redis 连接失败？
A: 确保 Redis 服务正在运行，检查配置中的 Redis 地址

### Q: 消息队列不工作？
A: 确保 Gateway Worker 服务正在运行，检查 Redis 连接

### Q: Swagger 文档无法访问？
A: 先运行 `make swagger` 生成文档，然后访问 `/swagger/index.html`

## 📞 获取帮助

1. **查看文档**：浏览 `docs/` 目录下的详细文档
2. **运行示例**：参考 `docs/demo/` 目录下的示例代码
3. **检查日志**：查看服务日志了解错误信息
4. **社区支持**：提交 Issue 或参与讨论

## 🎯 下一步

完成快速入门后，建议：

1. 阅读完整的 [FRAMEWORK.md](./FRAMEWORK.md) 文档
2. 浏览 [demo/](./demo/) 目录下的所有示例
3. 尝试创建自己的服务
4. 探索高级功能和最佳实践

---

**开始你的 GinForge 之旅吧！** 🚀
