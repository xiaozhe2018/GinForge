# 框架介绍

欢迎使用 **GinForge** - 一个基于 Go + Gin 的企业级微服务开发框架！

## 什么是 GinForge？

GinForge 是一个开箱即用的企业级 Web 开发框架，它整合了 Go 生态中最优秀的组件和最佳实践，帮助开发者快速构建高性能、可扩展的 Web 应用和微服务。

## 核心特性

### 🚀 高性能

- 基于 Gin 框架，处理速度快
- 支持并发请求处理
- 优化的数据库查询
- 内置缓存机制

### 🏗️ 模块化设计

- 清晰的分层架构（Controller、Service、Repository）
- 可插拔的中间件系统
- 灵活的配置管理
- 松耦合的组件设计

### 🔐 安全可靠

- JWT 身份认证
- RBAC 权限控制
- SQL 注入防护
- XSS 防护
- 请求限流

### 📦 开箱即用

- 用户认证系统
- 权限管理系统
- 文件上传服务
- WebSocket 支持
- 消息队列
- 缓存系统

## 技术栈

```
后端技术：
├── Go 1.21+               # 编程语言
├── Gin                    # Web 框架
├── GORM                   # ORM 框架
├── Redis                  # 缓存和队列
├── MySQL                  # 数据库
└── JWT                    # 身份认证

前端技术：
├── Vue 3                  # 前端框架
├── TypeScript             # 类型安全
├── Element Plus           # UI 组件库
├── Pinia                  # 状态管理
└── Vite                   # 构建工具
```

## 适用场景

GinForge 适合以下场景：

- ✅ 企业级管理后台
- ✅ RESTful API 服务
- ✅ 微服务架构
- ✅ SaaS 平台
- ✅ 中后台系统
- ✅ B2B 应用

## 快速示例

下面是一个简单的 HTTP 服务示例：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "goweb/pkg/config"
    "goweb/pkg/logger"
)

func main() {
    // 加载配置
    cfg := config.New()
    
    // 初始化日志
    log := logger.New("demo", cfg.GetString("log.level"))
    
    // 创建路由
    r := gin.Default()
    
    // 定义路由
    r.GET("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello, GinForge!",
        })
    })
    
    // 启动服务
    log.Info("server starting on :8080")
    r.Run(":8080")
}
```

## 架构设计

### 分层架构

```
┌─────────────────────────────────────┐
│           Controller Layer          │  ← HTTP 处理
├─────────────────────────────────────┤
│            Service Layer            │  ← 业务逻辑
├─────────────────────────────────────┤
│          Repository Layer           │  ← 数据访问
├─────────────────────────────────────┤
│            Model Layer              │  ← 数据模型
└─────────────────────────────────────┘
```

### 请求处理流程

```
Client Request
    ↓
Middleware (Auth, Log, etc.)
    ↓
Router → Controller
    ↓
Service (Business Logic)
    ↓
Repository (Database)
    ↓
Response
```

## 为什么选择 GinForge？

| 特性 | GinForge | 其他框架 |
|------|----------|----------|
| 学习曲线 | 平缓 | 陡峭 |
| 开发效率 | 高 | 中等 |
| 性能 | 优秀 | 良好 |
| 文档完善度 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| 社区支持 | 活跃 | 一般 |
| 企业级特性 | 完整 | 需自行集成 |

## 社区和支持

- 📖 [在线文档](https://docs.ginforge.com)
- 💬 [GitHub Issues](https://github.com/ginforge/ginforge/issues)
- 📧 [邮件支持](mailto:support@ginforge.com)
- 💡 [示例项目](https://github.com/ginforge/examples)

## 下一步

准备好开始了吗？让我们继续：

- [安装指南](./installation) - 学习如何安装 GinForge
- [快速开始](./quick-start) - 5分钟创建你的第一个应用
- [项目结构](./project-structure) - 了解框架的目录结构

---

**版本**: v1.0.0  
**最后更新**: 2025-10-15

