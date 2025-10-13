# GinForge 框架全面功能概览

## 🎯 框架定位
GinForge 是一个基于 Go + Gin 的**企业级微服务开发框架**，目标是"让开发更加简单"，提供完整的工程化解决方案。

## 📦 核心功能模块

### 1. 基础架构层
- **配置管理** (`pkg/config/`)
  - Viper 集成，支持 YAML + 环境变量 + 默认值
  - 多环境配置（dev/prod/test）
  - 热重载支持
  - 类型安全的配置获取

- **日志系统** (`pkg/logger/`)
  - 基于 Zap 的结构化日志
  - 多级别日志（Debug/Info/Warn/Error/Fatal）
  - 链路追踪支持
  - 字段扩展功能

- **错误处理** (`pkg/errors/`)
  - 统一错误码定义（1000-9999）
  - 业务错误码分类
  - 错误消息映射
  - 标准化错误响应

- **常量管理** (`pkg/constants/`)
  - 业务状态常量
  - 缓存键前缀
  - 文件类型定义
  - 验证规则常量

### 2. 数据访问层
- **数据库管理** (`pkg/db/`)
  - GORM 集成，支持 SQLite/MySQL/PostgreSQL
  - 连接池管理
  - 事务支持
  - 迁移管理
  - 查询日志

- **缓存系统** (`pkg/cache/`)
  - 内存缓存 + Redis 缓存
  - 统一缓存接口
  - TTL 管理
  - 命名空间支持
  - 缓存统计

- **模型定义** (`pkg/model/`)
  - 基础模型（BaseModel、SoftDeleteModel）
  - 分页支持
  - 业务模型（User、Merchant、Product、Order）
  - GORM 标签支持

### 3. 业务逻辑层
- **基类体系** (`pkg/base/`)
  - BaseService：服务层基类
  - BaseHandler：处理器基类
  - BaseRepository：仓储层基类
  - BaseController：控制器基类
  - 统一日志和错误处理

- **验证系统** (`pkg/validator/`)
  - 基于 go-playground/validator
  - 统一验证接口
  - 自定义验证规则
  - 错误信息本地化

- **工具函数** (`pkg/utils/`)
  - 加密工具（MD5/SHA1/SHA256/SHA512）
  - 字符串处理
  - 时间处理
  - UUID 生成

### 4. 网络通信层
- **中间件系统** (`pkg/middleware/`)
  - Recovery：panic 恢复
  - RequestID：请求追踪
  - AccessLogger：访问日志
  - CORS：跨域处理
  - JWT：认证中间件
  - RateLimit：限流中间件
  - Validation：参数校验
  - Cache：HTTP 缓存

- **统一响应** (`pkg/response/`)
  - 标准化 API 响应格式
  - 错误码映射
  - 链路追踪 ID
  - 类型安全

- **Gateway 客户端** (`pkg/gateway/`)
  - 服务间通信 SDK
  - HTTP 客户端封装
  - 超时和重试
  - 业务方法封装

- **服务发现** (`pkg/service/`)
  - 服务注册表
  - 服务客户端
  - 健康检查
  - 负载均衡

### 5. 监控运维层
- **监控指标** (`pkg/monitor/`)
  - Prometheus 指标收集
  - HTTP 请求指标
  - 业务指标
  - 数据库指标
  - 缓存指标
  - 自定义指标

- **健康检查** (`pkg/monitor/health.go`)
  - 数据库健康检查
  - Redis 健康检查
  - 缓存健康检查
  - 自定义健康检查
  - 整体健康状态

### 6. 存储管理
- **文件存储** (`pkg/storage/`)
  - 本地文件存储
  - 文件上传下载
  - 文件信息管理
  - 文件类型检测
  - 过期文件清理

### 7. 消息队列
- **消息队列** (`pkg/mq/`)
  - Redis 消息队列
  - 发布/订阅模式
  - 消费者组
  - 死信队列
  - 延迟消息
  - 消息重试

### 8. 分布式功能
- **分布式锁** (`pkg/lock/`)
  - Redis 分布式锁
  - 锁超时管理
  - 锁续期
  - 互斥锁封装
  - 锁状态检查

- **熔断器** (`pkg/circuit/`)
  - 熔断器模式
  - 状态管理（关闭/开启/半开）
  - 失败率统计
  - 自动恢复
  - 熔断器管理

### 9. 开发工具
- **代码生成** (`cmd/generator/`)
  - 服务模板生成
  - 代码脚手架
  - 模板系统

- **API 文档** (`pkg/swagger/`)
  - Swagger/OpenAPI 集成
  - 自动文档生成
  - 在线测试

- **测试支持**
  - 单元测试框架
  - 集成测试支持
  - Mock 工具

## 🏗️ 架构特点

### 1. 微服务架构
- 服务分离（user-api、merchant-api、admin-api、gateway）
- 服务间通信
- 统一配置管理
- 独立部署

### 2. 分层设计
- 表现层（Handler/Controller）
- 业务层（Service）
- 数据层（Repository）
- 基础设施层（Config/Logger/Cache/DB）

### 3. 工程化实践
- 统一代码风格
- 标准化目录结构
- 完整的错误处理
- 全面的日志记录
- 监控和告警

### 4. 云原生支持
- Docker 容器化
- Kubernetes 部署
- 健康检查
- 配置外部化
- 服务发现

## 🚀 使用场景

### 1. 企业级应用
- 电商平台
- 内容管理系统
- 用户管理系统
- 订单管理系统

### 2. 微服务架构
- 服务拆分
- 服务治理
- 分布式事务
- 服务监控

### 3. 高并发系统
- 限流和熔断
- 缓存优化
- 数据库优化
- 监控告警

## 📊 技术栈

### 核心框架
- **Web 框架**: Gin
- **数据库**: GORM + SQLite/MySQL/PostgreSQL
- **缓存**: Redis + 内存缓存
- **消息队列**: Redis Stream
- **监控**: Prometheus
- **日志**: Zap

### 开发工具
- **配置管理**: Viper
- **验证**: go-playground/validator
- **文档**: Swagger/OpenAPI
- **测试**: Go 内置测试框架

### 部署运维
- **容器化**: Docker
- **编排**: Docker Compose
- **容器编排**: Kubernetes
- **监控**: Prometheus + Grafana

## 🎯 框架优势

### 1. 开发效率
- 丰富的基类体系
- 代码生成工具
- 统一的开发规范
- 完整的示例和文档

### 2. 工程化程度
- 标准化目录结构
- 统一的错误处理
- 全面的日志记录
- 完整的监控体系

### 3. 可维护性
- 清晰的分层架构
- 模块化设计
- 接口抽象
- 配置外部化

### 4. 可扩展性
- 插件化架构
- 中间件系统
- 自定义组件
- 服务发现

### 5. 生产就绪
- 健康检查
- 监控指标
- 错误处理
- 性能优化

## 📈 性能特点

### 1. 高并发支持
- Gin 高性能 HTTP 框架
- 连接池管理
- 限流和熔断
- 缓存优化

### 2. 低延迟
- 内存缓存
- 连接复用
- 异步处理
- 批量操作

### 3. 高可用
- 熔断器模式
- 健康检查
- 服务降级
- 故障恢复

## 🔧 快速开始

```bash
# 1. 获取代码
git clone <your-repo> && cd goweb

# 2. 安装依赖
go mod tidy

# 3. 生成服务
go run ./cmd/generator -command=service -name=my-service

# 4. 启动服务
go run ./services/my-service/cmd/server

# 5. 查看文档
make swagger
```

## 📚 文档资源

- **框架文档**: `docs/FRAMEWORK.md`
- **使用示例**: `demo/` 目录
- **API 文档**: 各服务的 `/swagger/index.html`
- **配置说明**: `configs/` 目录
- **部署指南**: `deployments/` 目录

GinForge 框架提供了从开发到部署的完整解决方案，让开发者能够快速构建高质量的企业级微服务应用。
