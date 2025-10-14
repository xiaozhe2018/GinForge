# GinForge 微服务开发框架

**原则：让开发更加简单**

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![Gin Version](https://img.shields.io/badge/Gin-1.10.0-00ADD8?style=flat)
![Vue Version](https://img.shields.io/badge/Vue-3.4.0-4FC08D?style=flat&logo=vue.js)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)

一个功能完整、架构清晰、开箱即用的企业级微服务开发框架

</div>

## 📖 项目简介

GinForge 是基于 **Go + Gin + Vue3** 的企业级微服务开发框架，提供从开发到部署的完整工程化解决方案。框架集成了微服务开发中常用的技术栈和最佳实践，让开发者可以快速构建生产级的 Web 应用。

### 🌟 核心特点

- **🏗️ 微服务架构** - 多端分离（用户端/商户端/管理后台），服务独立部署
- **🎨 现代化前端** - Vue3 + TypeScript + Element Plus 管理后台
- **🔐 完善的 RBAC 权限** - 用户-角色-权限-菜单四级权限控制
- **🚀 开箱即用** - 丰富的基础库和代码生成器，快速启动项目
- **📦 工程化实践** - 统一配置、日志、错误处理、API 文档
- **☁️ 云原生支持** - Docker、Kubernetes、Istio 部署配置

## 🚀 快速开始

### 环境要求

#### 开发环境
| 软件 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.20+ | 后端开发语言 |
| Node.js | 20+ | 前端开发环境 |
| SQLite | - | 自动创建（无需安装） |

#### 生产环境
| 软件 | 版本要求 | 说明 |
|------|---------|------|
| Docker | 20+ | 容器运行环境 |
| Docker Compose | 1.29+ | 容器编排工具 |
| MySQL | 8.0+ | 生产数据库（自动部署） |
| Redis | 7.0+ | 缓存服务（自动部署） |

### ⚡ 开发环境 - 30秒快速启动

```bash
# 1. 克隆项目
git clone https://github.com/xiaozhe2018/GinForge.git
cd GinForge

# 2. 初始化项目（一次性）
./scripts/init.sh

# 3. 启动所有后端服务
./scripts/start-services.sh

# 4. 启动前端（新终端）
cd web/admin && npm run dev
```

### 🚀 生产环境 - 一键部署

```bash
# 1. 克隆项目
git clone https://github.com/xiaozhe2018/GinForge.git
cd GinForge

# 2. 配置生产环境
cd deployments
cp env.production.example .env.production
vim .env.production  # 修改密码和密钥

# 3. 一键部署（自动构建前端+启动所有服务）
./deploy.sh
```

**部署包含：**
- 🔹 7个 Go 微服务
- 🔹 Vue3 管理后台
- 🔹 MySQL 8.0 数据库
- 🔹 Redis 7.x 缓存
- 🔹 Nginx 反向代理

### 访问系统

#### 开发环境
| 服务 | 地址 | 说明 |
|------|------|------|
| **前端管理后台** 🎉 | http://localhost:3000 | 默认账号：admin/admin123 |
| 管理后台API | http://localhost:8083 | RESTful API |
| Swagger文档 | http://localhost:8083/swagger/index.html | 在线API文档 |
| 用户端API | http://localhost:8081 | 用户服务 |
| 商户端API | http://localhost:8082 | 商户服务 |
| API网关 | http://localhost:8080 | 统一网关 |

#### 生产环境
| 服务 | 地址 | 说明 |
|------|------|------|
| **前端+API** 🎉 | http://localhost | Nginx 统一入口 |
| API 网关 | http://localhost:8080 | 直接访问（调试用） |
| 健康检查 | http://localhost/healthz | 服务健康状态 |
| API 文档 | http://localhost/swagger | Swagger 文档 |

## 📚 完整文档

**所有文档都在 [docs/](./docs/) 目录下，请查看 [文档索引](./docs/INDEX.md) 获取完整文档列表。**

### 核心文档
- [📖 框架使用指南](./docs/FRAMEWORK.md) - 详细使用指南
- [⚡ 快速开始](./docs/QUICK_START.md) - 5分钟快速入门
- [🔍 功能概览](./docs/FRAMEWORK_OVERVIEW.md) - 框架功能全面概览
- [🚀 高级功能](./docs/ADVANCED_FEATURES.md) - 高级功能详解
- [💡 使用示例](./docs/demo/) - 各种使用示例和教程

## 🛠️ 功能特性

### 🎯 后端核心功能

#### 🏛️ 架构设计
- **微服务架构** - 6个独立服务（用户端、商户端、管理后台、网关、网关工作器、演示服务）
- **分层架构** - Handler → Service → Repository → Model 清晰分层
- **基类体系** - BaseController、BaseService、BaseRepository 减少重复代码
- **依赖注入** - 统一的服务注册和依赖管理

#### 🔐 安全认证
- **JWT 认证** - 基于 JWT Token 的无状态认证
- **Token 黑名单** - Redis 实现的 Token 失效机制，登出即失效
- **RBAC 权限** - 用户-角色-权限-菜单四级权限控制
- **密码加密** - Bcrypt 加密存储，防止明文泄露
- **操作审计** - 完整的操作日志记录（登录、登出、数据变更）

#### 📊 数据处理
- **数据库支持** - GORM ORM + SQLite/MySQL/PostgreSQL
- **自动迁移** - 数据库结构自动同步
- **软删除** - 数据安全删除，可恢复
- **关联查询** - 支持复杂的关联关系（一对一、一对多、多对多）
- **事务管理** - 自动事务和手动事务支持

#### ⚡ 性能优化
- **Redis 缓存** - 多级缓存策略，支持缓存预热和自动刷新
- **分布式锁** - Redis 分布式锁，防止并发问题
- **连接池** - 数据库和 Redis 连接池优化
- **限流熔断** - 令牌桶限流 + 熔断器保护
- **异步处理** - Go 协程异步处理非关键任务

#### 📝 日志监控
- **结构化日志** - Zap 高性能日志库
- **请求链路追踪** - Request ID 全链路追踪
- **健康检查** - 完整的健康检查端点
- **Prometheus 监控** - 内置监控指标采集
- **Swagger 文档** - 自动生成 OpenAPI 3.0 文档

#### 🔧 工具支持
- **CLI 工具** - 命令行工具：初始化、生成代码、部署等
- **代码生成器** - 自动生成服务模板（Handler、Service、Repository）
- **配置中心** - 支持配置热更新和多环境配置
- **测试框架** - 单元测试和集成测试工具

### 🎨 前端核心功能

#### 💎 技术栈
- **Vue 3** - Composition API + `<script setup>` 语法
- **TypeScript** - 完整的类型定义，类型安全
- **Element Plus** - 企业级 UI 组件库
- **Vue Router** - 前端路由管理
- **Pinia** - 轻量级状态管理（可选）
- **Axios** - HTTP 请求库，统一拦截器

#### 🎭 功能模块
- ✅ **登录认证** - JWT Token 认证，自动登录
- ✅ **用户管理** - 用户 CRUD、状态管理、角色分配
- ✅ **角色管理** - 角色 CRUD、权限分配、菜单分配
- ✅ **菜单管理** - 树形菜单、图标选择、路由配置
- ✅ **权限管理** - 权限 CRUD、按钮级权限控制
- ✅ **个人设置** - 个人信息修改、密码修改
- ✅ **系统管理** - 系统配置、操作日志、系统监控
- ✅ **仪表盘** - 数据统计、图表展示

#### 🎯 用户体验
- **响应式设计** - 适配桌面和移动端
- **暗黑模式** - 支持主题切换（可扩展）
- **权限菜单** - 根据用户权限动态渲染菜单
- **表单验证** - 完整的前端表单验证
- **错误处理** - 统一的错误提示和处理
- **加载状态** - Loading 和骨架屏优化体验

### 🚀 DevOps 支持

#### 🐳 容器化
- **Dockerfile** - 多阶段构建，镜像体积小
- **Docker Compose** - 一键启动完整环境
- **镜像优化** - Alpine Linux 基础镜像

#### ☸️ Kubernetes
- **Deployment** - 服务部署配置
- **Service** - 服务发现和负载均衡
- **ConfigMap** - 配置管理
- **Istio 支持** - 服务网格配置（Gateway、VirtualService、DestinationRule）

#### 📦 中间件集成
- **统一配置** - Viper 配置管理（YAML、环境变量、默认值）
- **统一日志** - Zap 结构化日志（JSON/Console 输出）
- **统一响应** - 标准化 JSON 响应格式
- **中间件链** - Recovery、RequestID、CORS、JWT、限流、日志、缓存
- **消息队列** - Redis Stream 实现的消息队列和延时队列

## 🏗️ 项目结构

```
goweb/
├── bin/                    # 编译后的二进制文件（不提交到Git）
│   ├── admin-api          # 管理后台服务
│   ├── user-api           # 用户端服务
│   ├── merchant-api       # 商户端服务
│   ├── gateway            # API网关
│   ├── gateway-worker     # 网关工作器
│   ├── demo               # 示例服务
│   ├── file-api           # 文件上传服务
│   └── ginforge           # CLI工具
├── pkg/                    # 共享基础库
│   ├── config/            # 配置管理
│   ├── logger/            # 日志系统
│   ├── middleware/        # 中间件
│   ├── response/          # 统一响应
│   ├── db/                # 数据库管理
│   ├── redis/             # Redis管理
│   ├── storage/           # 文件存储
│   ├── model/             # 数据模型
│   ├── utils/             # 工具函数
│   └── ...
├── services/              # 微服务（源代码）
│   ├── user-api/         # 用户端API
│   ├── merchant-api/     # 商户端API
│   ├── admin-api/        # 管理后台API
│   ├── gateway/          # API网关
│   ├── gateway-worker/   # 网关工作器
│   ├── file-api/         # 文件上传服务
│   └── demo/             # 示例服务
├── cmd/                   # 命令行工具
│   └── cli/              # CLI工具源代码
├── templates/             # 代码生成模板
├── deployments/           # 部署配置
│   ├── docker/           # Docker配置
│   └── k8s/              # Kubernetes配置
├── docs/                  # 📚 完整文档
│   ├── INDEX.md          # 文档索引
│   ├── demo/             # 使用示例
│   └── ...
└── configs/               # 配置文件
```

## 🎯 核心特性

### 1. 微服务架构
- 多端分离设计（用户端/商户端/管理后台）
- 服务独立部署和扩展
- 统一API网关
- 服务间通信

### 2. 工程化实践
- 统一配置管理
- 结构化日志记录
- 标准化错误处理
- 完整的监控体系

### 3. 开发效率
- 代码生成工具
- 丰富的基类体系
- 完整的示例和文档
- 统一的开发规范

### 4. 生产就绪
- 健康检查
- 监控指标
- 限流和熔断
- 容器化部署

## 📖 使用示例

### 🔰 基础使用

#### 1. 创建新服务
```bash
# 使用代码生成器创建支付服务
go run ./cmd/generator -command=service -name=payment

# 生成的文件结构：
# services/payment/
# ├── cmd/server/main.go          # 服务入口
# ├── internal/
# │   ├── handler/payment_handler.go   # HTTP 处理器
# │   ├── service/payment_service.go   # 业务逻辑
# │   └── router/router.go             # 路由配置
# └── docs/                            # API文档
```

#### 2. 创建 API 接口
```go
// services/payment/internal/handler/payment_handler.go
package handler

type PaymentHandler struct {
    *base.BaseHandler
    service *service.PaymentService
}

// CreateOrder 创建支付订单
// @Summary 创建支付订单
// @Tags 支付管理
// @Accept json
// @Produce json
// @Param request body OrderRequest true "订单信息"
// @Success 200 {object} response.Response
// @Router /api/v1/payment/orders [post]
func (h *PaymentHandler) CreateOrder(c *gin.Context) {
    var req OrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.Error(c, err)
        return
    }
    
    order, err := h.service.CreateOrder(c.Request.Context(), &req)
    if err != nil {
        h.Error(c, err)
        return
    }
    
    h.Success(c, order)
}
```

#### 3. 生成 Swagger 文档
```bash
# 生成所有服务的 API 文档
make swagger

# 或单独生成某个服务的文档
swag init -g services/admin-api/cmd/server/main.go -o services/admin-api/docs

# 访问在线文档
# http://localhost:8083/swagger/index.html
```

#### 4. 使用缓存
```go
import "goweb/pkg/redis"

// 基础缓存操作
cache := redis.NewCacheManager(redisClient)

// 设置缓存（5分钟过期）
cache.Set(ctx, "user:1001", userData, 5*time.Minute)

// 获取缓存
var user User
err := cache.Get(ctx, "user:1001", &user)

// 删除缓存
cache.Delete(ctx, "user:1001")

// 批量删除（按模式匹配）
cache.DeletePattern(ctx, "user:*")
```

#### 5. 使用消息队列
```go
import "goweb/pkg/redis"

// 创建队列
queue := redis.NewQueue(redisClient, "order-queue")

// 发送消息
queue.Publish(ctx, map[string]interface{}{
    "order_id": 12345,
    "amount": 99.99,
})

// 消费消息
queue.Subscribe(ctx, func(msg interface{}) error {
    // 处理消息
    fmt.Printf("处理订单: %v\n", msg)
    return nil
})
```

#### 6. 使用分布式锁
```go
import "goweb/pkg/redis"

lock := redis.NewDistributedLock(redisClient)

// 获取锁（10秒过期）
if lock.Acquire(ctx, "order:1001", 10*time.Second) {
    defer lock.Release(ctx, "order:1001")
    
    // 执行需要互斥的操作
    // ...
}
```

### 🧪 测试

#### 运行所有测试
```bash
make test
```

#### 运行测试并生成覆盖率报告
```bash
make test-coverage
# 查看 coverage.html
```

#### 运行特定包的测试
```bash
go test -v ./pkg/middleware/...
go test -v ./services/admin-api/internal/service/...
```

#### API 测试示例
```bash
# 测试登录
curl -X POST http://localhost:8083/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 测试带Token的请求
TOKEN="your-token-here"
curl -X GET http://localhost:8083/api/v1/admin/users?page=1&page_size=10 \
  -H "Authorization: Bearer $TOKEN"

# 测试创建用户
curl -X POST http://localhost:8083/api/v1/admin/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "user@example.com",
    "real_name": "张三",
    "password": "123456"
  }'
```

### 🐳 部署

#### 开发环境部署

```bash
# 方式一：使用脚本（推荐）
./scripts/start-services.sh

# 方式二：使用 Makefile
make run

# 停止服务
./scripts/stop-services.sh
# 或
make stop
```

**特点：**
- ✅ 快速启动，适合开发调试
- ✅ 使用 SQLite，无需外部数据库
- ✅ 热重载支持
- ✅ 日志输出到文件

#### 生产环境部署（Docker）⭐

**一键部署（推荐）：**

```bash
# 使用自动化部署脚本
./deployments/deploy.sh
```

**手动部署：**

```bash
# 1. 配置环境变量
cd deployments
cp env.production.example .env.production
vim .env.production  # 修改数据库密码、JWT密钥等

# 2. 构建前端
cd ../web/admin
npm install && npm run build

# 3. 启动所有服务（包括 MySQL + Redis + Nginx）
cd ../../deployments
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d

# 4. 查看服务状态
docker-compose -f docker-compose.prod.yml ps

# 5. 查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

**生产环境特性：**
- ✅ MySQL 8.0 数据库
- ✅ Redis 7.x 缓存
- ✅ Nginx 反向代理
- ✅ 健康检查（自动重启）
- ✅ 资源限制（CPU/内存）
- ✅ 日志轮转（自动清理）
- ✅ 数据持久化（Docker Volumes）
- ✅ 环境变量隔离
- ✅ 安全加固（强密码、网络隔离）
- ✅ 零停机更新

**访问地址：**
```
前端: http://localhost         (Nginx 统一入口)
API:  http://localhost/api      (通过 Gateway)
文档: http://localhost/swagger  (API 文档)
```

**详细文档：** [生产环境部署指南](./deployments/PRODUCTION_DEPLOYMENT.md)

#### Docker Compose 开发环境
```bash
# 构建镜像
make docker

# 使用 Docker Compose 启动（开发版本）
make compose

# 查看运行状态
docker-compose -f deployments/docker-compose.yml ps

# 停止服务
make compose-down
```

#### Kubernetes 部署
```bash
# 应用配置
kubectl apply -f deployments/k8s/

# 查看服务状态
kubectl get pods
kubectl get services

# 查看日志
kubectl logs -f <pod-name>
```

#### Istio 服务网格
```bash
# 部署 Istio 配置
kubectl apply -f deployments/k8s/istio/

# 查看流量路由
kubectl get virtualservices
kubectl get destinationrules
```

## 🔧 开发命令

```bash
# 构建相关
make build          # 构建所有服务
make clean          # 清理构建文件

# 运行相关
make run            # 启动所有后端服务
make stop           # 停止所有服务
make restart        # 重启所有服务
make status         # 查看服务状态

# 测试相关
make test              # 运行所有测试
make test-coverage     # 生成测试覆盖率报告
make test-integration  # 运行集成测试
make benchmark         # 运行性能测试

# 文档相关
make swagger        # 生成 Swagger 文档

# 前端相关
make web-install    # 安装前端依赖
make web-dev        # 启动前端开发服务器
make web-build      # 构建前端生产版本

# 部署相关
make docker         # 构建 Docker 镜像
make compose        # 启动 Docker Compose
make compose-down   # 停止 Docker Compose

# 开发环境
make dev            # 快速启动开发环境（后端+文档）
make dev-full       # 启动完整开发环境（后端+前端）
```

## 💡 最佳实践

### 1. 代码组织
- **按功能分层**：Handler → Service → Repository → Model
- **单一职责**：每个函数只做一件事
- **依赖注入**：通过构造函数注入依赖
- **接口抽象**：面向接口编程，方便测试和替换

### 2. 错误处理
```go
// 使用框架的统一错误码
if err != nil {
    return errors.NewBusinessError(errors.CodeUserNotFound, "用户不存在")
}

// 在 Handler 层统一处理错误
h.Error(c, err)  // 自动根据错误类型返回对应状态码
```

### 3. 日志记录
```go
// 使用结构化日志
logger.Info("用户登录",
    zap.String("username", username),
    zap.String("ip", ip),
    zap.String("request_id", requestID))

// 记录错误日志
logger.Error("数据库查询失败", zap.Error(err))
```

### 4. 配置管理
```go
// 配置优先级：环境变量 > YAML 文件 > 默认值
config := config.Load()

// 支持配置热更新
config.Watch(func(cfg *config.Config) {
    // 配置变更回调
})
```

### 5. 数据库操作
```go
// 使用事务
err := db.Transaction(func(tx *gorm.DB) error {
    // 在事务中执行多个操作
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    if err := tx.Create(&userProfile).Error; err != nil {
        return err
    }
    return nil
})
```

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**GinForge 框架 - 让开发更加简单** 🚀

> 📚 **查看完整文档**: [docs/INDEX.md](./docs/INDEX.md)