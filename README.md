# GinForge 微服务开发框架

**原则：让开发更加简单**

基于 Go + Gin 的企业级微服务开发框架，提供完整的工程化解决方案。

## 🚀 快速开始

### 环境要求
- Go 1.20+
- Docker (可选)
- Redis (可选)

### 安装依赖
```bash
git clone <your-repo> && cd goweb
go mod tidy
```

### 启动服务
```bash
make run
```

### 访问API
- 用户端API: http://localhost:8081
- 商户端API: http://localhost:8082
- 管理后台API: http://localhost:8083
- API网关: http://localhost:8080
- **前端管理后台: http://localhost:3000** 🎉

## 📚 完整文档

**所有文档都在 [docs/](./docs/) 目录下，请查看 [文档索引](./docs/INDEX.md) 获取完整文档列表。**

### 核心文档
- [📖 框架使用指南](./docs/FRAMEWORK.md) - 详细使用指南
- [⚡ 快速开始](./docs/QUICK_START.md) - 5分钟快速入门
- [🔍 功能概览](./docs/FRAMEWORK_OVERVIEW.md) - 框架功能全面概览
- [🚀 高级功能](./docs/ADVANCED_FEATURES.md) - 高级功能详解
- [💡 使用示例](./docs/demo/) - 各种使用示例和教程

## 🛠️ 主要功能

- **微服务架构**: 多端分离，服务独立部署
- **前端管理后台**: Vue3 + Element Plus 现代化界面 🎉
- **统一配置**: 环境变量 + YAML + 默认值
- **统一日志**: Zap 结构化日志，支持链路追踪
- **统一响应**: 标准化 API 响应格式
- **中间件系统**: Recovery、RequestID、CORS、JWT、限流等
- **数据库支持**: GORM + SQLite/MySQL/PostgreSQL
- **缓存系统**: Redis + 内存缓存
- **消息队列**: Redis Stream 支持
- **API文档**: Swagger/OpenAPI 自动生成
- **代码生成**: 服务模板和脚手架
- **监控指标**: Prometheus 集成
- **健康检查**: 完整的健康检查体系
- **Docker支持**: 容器化部署
- **K8s支持**: Kubernetes 部署清单

## 🏗️ 项目结构

```
goweb/
├── pkg/                    # 共享基础库
│   ├── config/            # 配置管理
│   ├── logger/            # 日志系统
│   ├── middleware/        # 中间件
│   ├── response/          # 统一响应
│   ├── db/               # 数据库管理
│   ├── cache/            # 缓存系统
│   ├── model/            # 数据模型
│   ├── utils/            # 工具函数
│   ├── testing/          # 测试框架
│   ├── mesh/             # 服务网格
│   └── ...
├── services/              # 微服务
│   ├── user-api/         # 用户端API
│   ├── merchant-api/     # 商户端API
│   ├── admin-api/        # 管理后台API
│   ├── gateway/          # API网关
│   ├── gateway-worker/   # 网关工作器
│   └── demo/             # 示例服务
├── cmd/                   # 命令行工具
│   └── cli/              # CLI工具
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

### 创建新服务
```bash
go run ./cmd/generator -command=service -name=payment
```

### 生成Swagger文档
```bash
make swagger
```

### 运行测试
```bash
make test
```

### Docker部署
```bash
make compose
```

## 🔧 开发命令

```bash
make build      # 构建所有服务
make run        # 启动所有服务
make stop       # 停止所有服务
make restart    # 重启所有服务
make test       # 运行测试
make swagger    # 生成API文档
make clean      # 清理构建文件
```

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**GinForge 框架 - 让开发更加简单** 🚀

> 📚 **查看完整文档**: [docs/INDEX.md](./docs/INDEX.md)