# 更新日志

本项目的所有重要变更都会记录在这个文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 新增
- **文件上传微服务**
  - 独立的文件上传、下载、管理服务
  - 支持本地存储和云存储(OSS/S3/MinIO)
  - 文件元数据管理和统计
  - 文件秒传、断点续传支持
  - 完整的API文档

### 计划中
- Dashboard 统计 API 实现
- 系统配置管理 API
- 数据导出功能
- 多语言支持

## [1.0.0] - 2025-01-13

### 🎉 首次发布

#### 新增
- **微服务架构**
  - 6个独立微服务（user-api, merchant-api, admin-api, gateway, gateway-worker, demo）
  - 清晰的分层架构（Handler → Service → Repository → Model）
  - 统一的基类体系

- **核心功能**
  - JWT 认证系统（含 Token 黑名单机制）
  - RBAC 权限管理（用户-角色-权限-菜单）
  - 完整的管理后台（Vue3 + TypeScript + Element Plus）
  - 用户管理（CRUD、状态管理、角色分配）
  - 角色管理（CRUD、权限分配、菜单分配）
  - 菜单管理（树形结构、图标选择、路由配置）
  - 权限管理（CRUD、按钮级权限）
  - 个人设置（信息修改、密码修改）

- **基础库**
  - 统一配置管理（Viper）
  - 结构化日志（Zap）
  - 数据库支持（GORM + SQLite/MySQL/PostgreSQL）
  - Redis 集成（缓存、分布式锁、消息队列）
  - 中间件系统（JWT、CORS、限流、日志、恢复等）
  - 统一响应格式
  - 错误码管理

- **高级功能**
  - 分布式锁
  - 消息队列（Redis Stream）
  - 延时队列
  - 熔断器
  - 限流器
  - 健康检查
  - Prometheus 监控
  - Swagger API 文档

- **开发工具**
  - CLI 工具
  - 代码生成器
  - 测试框架
  - Makefile 自动化

- **部署支持**
  - Docker 镜像
  - Docker Compose 配置
  - Kubernetes 部署清单
  - Istio 服务网格配置

- **文档系统**
  - 完整的项目文档
  - API 使用文档
  - 部署指南
  - 开发示例
  - 故障排查指南

#### 技术栈
- **后端**: Go 1.24+, Gin 1.10, GORM
- **前端**: Vue 3.4, TypeScript 5.3, Element Plus 2.4
- **数据库**: SQLite, MySQL, PostgreSQL
- **缓存**: Redis
- **监控**: Prometheus
- **部署**: Docker, Kubernetes, Istio

#### 性能指标
- API 响应时间 < 100ms（P95）
- 支持 1000+ QPS（单实例）
- 测试覆盖率 > 60%

## [0.2.0] - 2025-01-11

### 新增
- 完成前后端完整对接
- 实现 Token 黑名单机制
- 添加操作日志记录
- 完善 API 文档

### 修复
- 修复登录后立即退出的问题
- 修复 TypeScript 类型错误
- 修复 CORS 跨域问题

### 优化
- 优化数据库查询性能
- 改进错误处理机制
- 完善日志记录

## [0.1.0] - 2025-01-05

### 新增
- 初始化项目结构
- 实现基础框架
- 添加核心中间件
- 创建示例服务

---

## 版本说明

### 语义化版本格式

格式：`主版本号.次版本号.修订号`

- **主版本号**：不兼容的 API 修改
- **次版本号**：向下兼容的功能性新增
- **修订号**：向下兼容的问题修正

### 变更类型

- **新增（Added）**：新功能
- **变更（Changed）**：已有功能的变更
- **弃用（Deprecated）**：即将移除的功能
- **移除（Removed）**：已移除的功能
- **修复（Fixed）**：任何 bug 的修复
- **安全（Security）**：安全相关的改进

---

## 升级指南

### 从 0.x 升级到 1.0

#### 数据库迁移
```bash
# 备份数据
mysqldump -u root -p gin_forge > backup.sql

# 运行迁移脚本
go run ./scripts/migrate.go
```

#### 配置文件更新
```yaml
# 旧版本
app:
  name: goweb

# 新版本
app:
  name: GinForge
```

#### API 变更
- `/api/login` → `/api/v1/admin/auth/login`
- `/api/users` → `/api/v1/admin/users`

详细升级指南请查看 [升级文档](./docs/UPGRADE.md)。

---

## 贡献者

感谢所有为 GinForge 做出贡献的开发者！

完整贡献者列表请访问：[Contributors](https://github.com/xiaozhe2018/GinForge/graphs/contributors)

---

**查看完整发布历史**：[GitHub Releases](https://github.com/xiaozhe2018/GinForge/releases)

