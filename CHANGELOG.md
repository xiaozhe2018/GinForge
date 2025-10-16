# 更新日志

所有重要的项目变更都会记录在这个文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)。

## [1.0.0] - 2025-10-16

### ✨ 新增

#### 核心功能
- 🚀 **一键CRUD生成器** - 从数据库表自动生成全套代码（Model + Repository + Service + Handler + 前端）
- 🎨 **完整后台管理系统** - Vue3 + Element Plus，17个功能页面
- 🔐 **RBAC权限系统** - 用户-角色-权限-菜单四级控制
- 🏗️ **微服务架构** - 8个独立服务 + API网关

#### 基础设施
- 📦 **82个共享包** - 中间件、数据库、缓存、配置、日志等
- 🐳 **Docker部署** - 完整的docker-compose配置 + Nginx反向代理
- ☸️ **K8s支持** - Deployment + Service + Istio服务网格配置
- 📚 **46个文档** - 从快速开始到高级功能，全覆盖

#### 开发工具
- 🔧 **代码生成器** - 7个模板，自动生成1000+行代码
- 📊 **Swagger文档** - API文档自动生成
- 🧪 **测试框架** - 单元测试和集成测试支持

### 🎨 前端
- Vue 3 + TypeScript + Element Plus
- 用户管理、角色管理、权限管理、菜单管理
- 系统配置、操作日志、个人资料
- 响应式设计，暗黑模式支持

### 🔐 安全特性
- JWT认证 + Token黑名单
- Bcrypt密码加密
- 操作日志审计
- 登录失败锁定
- CORS跨域控制

### ⚡ 性能优化
- Redis多级缓存
- 数据库连接池
- 分布式锁
- 令牌桶限流
- 熔断器保护

### 📖 文档
- 快速开始指南
- 框架使用指南
- 功能概览
- 高级功能详解
- 16个示例教程
- 部署文档

---

## 路线图

### v1.1.0 (计划中)
- [ ] CI/CD配置（GitHub Actions）
- [ ] 单元测试覆盖率提升至60%+
- [ ] 性能测试报告
- [ ] 在线演示站点

### v1.2.0 (计划中)
- [ ] 多租户支持
- [ ] 插件系统
- [ ] 更多数据库支持（MongoDB、TiDB）
- [ ] GraphQL支持

### v2.0.0 (规划中)
- [ ] gRPC支持
- [ ] Service Mesh完整方案
- [ ] 可观测性增强（Tracing、Logging、Metrics）
- [ ] 低代码平台

---

## 贡献者

感谢所有为这个项目做出贡献的人！

<a href="https://github.com/xiaozhe2018/GinForge/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=xiaozhe2018/GinForge" />
</a>

---

[Unreleased]: https://github.com/xiaozhe2018/GinForge/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/xiaozhe2018/GinForge/releases/tag/v1.0.0
