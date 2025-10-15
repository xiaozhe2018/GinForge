# 项目结构

了解 GinForge 的目录结构，帮助你快速定位代码和资源。

## 📁 整体结构

```
GinForge/
├── cmd/                    # 命令行工具
│   ├── cli/               # CLI 命令
│   └── generator/         # 代码生成器
├── configs/               # 配置文件
├── database/              # 数据库相关
│   └── migrations/        # 数据库迁移脚本
├── deployments/           # 部署配置
│   ├── docker/           # Docker 配置
│   └── k8s/              # Kubernetes 配置
├── docs/                  # 文档
├── pkg/                   # 公共包（核心库）
├── services/              # 微服务
│   ├── admin-api/        # 管理后台 API
│   ├── user-api/         # 用户端 API
│   ├── gateway/          # API 网关
│   └── ...
├── templates/             # 模板文件
├── web/                   # 前端项目
│   └── admin/            # 管理后台前端
└── scripts/               # 脚本文件
```

## 🎯 核心目录详解

### 1. `pkg/` - 公共包（最重要）

这是框架的核心，包含所有可复用的组件：

```
pkg/
├── base/              # 基类体系 ⭐
│   ├── controller.go  # 控制器基类
│   ├── handler.go     # 处理器基类
│   ├── service.go     # 服务基类
│   └── repository.go  # 仓储基类
│
├── config/            # 配置管理 ⭐
│   └── config.go      # 配置加载和读取
│
├── db/                # 数据库 ⭐
│   ├── db.go          # 数据库连接
│   └── init.go        # 初始化脚本
│
├── redis/             # Redis 客户端 ⭐
│   ├── client.go      # Redis 封装
│   ├── cache.go       # 缓存操作
│   ├── queue.go       # 消息队列
│   ├── lock.go        # 分布式锁
│   └── delayed_worker.go  # 延时队列
│
├── middleware/        # 中间件 ⭐
│   ├── jwt.go         # JWT 认证
│   ├── logger.go      # 访问日志
│   ├── recovery.go    # Panic 恢复
│   ├── rate_limit.go  # 限流
│   ├── cache.go       # HTTP 缓存
│   └── validation.go  # 参数验证
│
├── response/          # 统一响应 ⭐
│   └── response.go    # API 响应格式
│
├── logger/            # 日志系统
│   └── logger.go      # 结构化日志
│
├── errors/            # 错误管理
│   └── codes.go       # 错误码定义
│
├── model/             # 数据模型
│   ├── base.go        # 基础模型
│   └── user.go        # 用户模型
│
├── validator/         # 验证器
│   └── validator.go   # 参数验证
│
├── utils/             # 工具函数
│   ├── crypto.go      # 加密工具
│   ├── string.go      # 字符串处理
│   └── time.go        # 时间处理
│
├── websocket/         # WebSocket
│   ├── manager.go     # 连接管理
│   ├── client.go      # 客户端
│   └── message.go     # 消息定义
│
├── upload/            # 文件上传
│   └── upload.go      # 上传处理
│
├── storage/           # 文件存储
│   ├── local/         # 本地存储
│   └── factory/       # 存储工厂
│
├── notification/      # 通知服务
│   ├── notification.go  # 通知接口
│   ├── email/         # 邮件通知
│   └── sms/           # 短信通知
│
├── monitor/           # 监控
│   ├── prometheus.go  # Prometheus 指标
│   └── health.go      # 健康检查
│
└── security/          # 安全
    └── security.go    # 安全工具
```

### 2. `services/` - 微服务

每个服务都遵循统一的结构：

```
services/admin-api/
├── cmd/
│   └── server/
│       └── main.go        # 服务入口
│
├── internal/              # 内部代码（不可导入）
│   ├── handler/          # HTTP 处理器
│   │   ├── admin_auth_handler.go
│   │   ├── admin_user_handler.go
│   │   └── ...
│   │
│   ├── service/          # 业务逻辑
│   │   ├── admin_service.go
│   │   ├── admin_system_service.go
│   │   └── ...
│   │
│   ├── repository/       # 数据访问
│   │   ├── admin_repository.go
│   │   └── ...
│   │
│   ├── model/            # 数据模型
│   │   ├── admin_user.go
│   │   ├── admin_role.go
│   │   └── ...
│   │
│   └── router/           # 路由配置
│       └── router.go
│
└── docs/                  # Swagger 文档
    └── swagger.yaml
```

### 3. `web/admin/` - 前端项目

```
web/admin/
├── src/
│   ├── api/              # API 接口定义
│   │   ├── auth.ts
│   │   ├── user.ts
│   │   └── system.ts
│   │
│   ├── views/            # 页面组件
│   │   ├── Login.vue     # 登录页
│   │   ├── Dashboard.vue # 仪表盘
│   │   ├── Users.vue     # 用户列表
│   │   └── Documentation/ # 文档中心
│   │
│   ├── layout/           # 布局组件
│   │   └── index.vue     # 主布局
│   │
│   ├── components/       # 公共组件
│   ├── stores/           # 状态管理
│   ├── router/           # 路由配置
│   ├── utils/            # 工具函数
│   └── styles/           # 样式文件
│
├── public/               # 静态资源
├── package.json          # 依赖配置
└── vite.config.ts        # Vite 配置
```

## 🏗️ 分层架构

GinForge 采用经典的分层架构：

```
┌─────────────────────────────────────┐
│         HTTP Layer (Gin Router)     │  ← 路由层
├─────────────────────────────────────┤
│      Middleware Layer (中间件)       │  ← 中间件层
├─────────────────────────────────────┤
│     Handler/Controller Layer        │  ← 处理器层
│     (HTTP 请求处理，参数绑定)          │
├─────────────────────────────────────┤
│         Service Layer               │  ← 业务逻辑层
│     (核心业务逻辑，事务管理)           │
├─────────────────────────────────────┤
│       Repository Layer              │  ← 数据访问层
│     (数据库操作，CRUD 封装)            │
├─────────────────────────────────────┤
│         Model Layer                 │  ← 模型层
│     (数据结构定义，GORM 标签)         │
├─────────────────────────────────────┤
│    Infrastructure Layer             │  ← 基础设施层
│  (DB, Redis, Logger, Config...)     │
└─────────────────────────────────────┘
```

## 📦 核心包说明

### `pkg/base/` - 基类体系

提供四大基类，所有业务代码都应该继承这些基类：

| 基类 | 作用 | 使用场景 |
|------|------|---------|
| `BaseHandler` | HTTP 处理器基类 | 所有 HTTP 接口 |
| `BaseService` | 业务服务基类 | 所有业务逻辑 |
| `BaseRepository` | 数据仓储基类 | 所有数据库操作 |
| `BaseController` | 控制器基类 | RESTful 控制器 |

**优势**：
- ✅ 统一的日志记录
- ✅ 统一的错误处理
- ✅ 统一的响应格式
- ✅ 减少重复代码

### `pkg/middleware/` - 中间件

| 中间件 | 功能 | 说明 |
|--------|------|------|
| `Recovery` | Panic 恢复 | 防止服务崩溃 |
| `RequestID` | 请求追踪 | 链路追踪 |
| `AccessLogger` | 访问日志 | 记录所有请求 |
| `JWT` | JWT 认证 | Token 验证 |
| `CORS` | 跨域处理 | 前后端分离 |
| `RateLimit` | 限流 | 防止滥用 |
| `Cache` | HTTP 缓存 | 性能优化 |
| `OperationLog` | 操作日志 | 审计追踪 |

### `pkg/config/` - 配置管理

配置优先级（从高到低）：

```
环境变量 > .env 文件 > config.yaml > 默认值
```

示例：

```go
// 读取配置
cfg := config.New()
port := cfg.GetInt("app.port")           // 读取整数
dbHost := cfg.GetString("database.host") // 读取字符串
isProd := cfg.IsProduction()             // 判断环境
```

### `pkg/response/` - 统一响应

所有 API 响应都使用统一格式：

```go
// 成功响应
response.Success(c, data)

// 错误响应
response.BadRequest(c, "参数错误")
response.Unauthorized(c, "未授权")
response.NotFound(c, "资源不存在")
response.InternalError(c, "服务器错误")
```

响应格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {...},
  "trace_id": "xxx-xxx-xxx"
}
```

## 🔄 请求处理流程

一个完整的 API 请求处理流程：

```
1. 客户端发送请求
   ↓
2. Gin Router 匹配路由
   ↓
3. 中间件链处理（RequestID → Logger → CORS → JWT → ...）
   ↓
4. Handler 处理请求
   ├─ 参数绑定和验证
   ├─ 调用 Service 层
   └─ 返回响应
   ↓
5. Service 执行业务逻辑
   ├─ 数据验证
   ├─ 调用 Repository
   ├─ 处理业务规则
   └─ 返回结果
   ↓
6. Repository 访问数据库
   ├─ CRUD 操作
   ├─ 事务管理
   └─ 返回数据
   ↓
7. 响应返回客户端
```

## 📚 示例代码位置

如果你想查看实际代码示例：

- **基类使用**: `docs/demo/base_classes_usage.md`
- **路由定义**: `docs/demo/router_response.md`
- **数据库操作**: `docs/demo/db.md`
- **缓存使用**: `docs/demo/cache.md`
- **中间件**: `docs/demo/middleware.md`

## 🎯 开发新功能的建议结构

当你需要添加新功能时，按照以下结构组织代码：

```
services/your-service/
└── internal/
    ├── model/
    │   └── your_model.go        # 1. 先定义数据模型
    │
    ├── repository/
    │   └── your_repository.go   # 2. 实现数据访问
    │
    ├── service/
    │   └── your_service.go      # 3. 编写业务逻辑
    │
    ├── handler/
    │   └── your_handler.go      # 4. 创建 HTTP 处理器
    │
    └── router/
        └── router.go            # 5. 注册路由
```

## 🔍 快速定位指南

| 我想... | 去哪里找 |
|---------|---------|
| 修改配置 | `configs/config.yaml` |
| 查看日志 | `logs/` 目录 |
| 添加 API | `services/*/internal/handler/` |
| 修改业务逻辑 | `services/*/internal/service/` |
| 数据库操作 | `services/*/internal/repository/` |
| 添加中间件 | `pkg/middleware/` |
| 修改响应格式 | `pkg/response/` |
| 添加工具函数 | `pkg/utils/` |
| 修改前端页面 | `web/admin/src/views/` |
| 修改前端 API | `web/admin/src/api/` |

## 💡 命名规范

### 文件命名

- Go 文件：`snake_case.go`（例如：`user_handler.go`）
- Vue 文件：`PascalCase.vue`（例如：`UserList.vue`）
- TypeScript 文件：`camelCase.ts`（例如：`userApi.ts`）

### 代码命名

```go
// 结构体：PascalCase
type UserService struct {}

// 函数/方法：PascalCase（公开），camelCase（私有）
func NewUserService() *UserService {}
func (s *UserService) GetUser() {}
func (s *UserService) validateUser() {} // 私有方法

// 变量：camelCase
var userName string

// 常量：PascalCase 或 UPPER_CASE
const UserTypeAdmin = 1
const MAX_RETRY_COUNT = 3
```

## 🎯 下一步

了解项目结构后，你可以：

- [学习配置系统](../core-concepts/configuration) - 掌握配置管理
- [创建你的第一个 API](../core-concepts/routing) - 开始开发
- [使用基类体系](../api-reference/base-classes) - 提高开发效率

---

**提示**: 建议在 IDE 中打开项目，边看文档边浏览代码，效果更佳！

