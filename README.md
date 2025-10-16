<div align="center">

# GinForge

**🚀 企业级 Go 微服务开发框架**

*30秒启动，一键生成CRUD，开箱即用的管理后台*

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Gin Version](https://img.shields.io/badge/Gin-1.10.0-00ADD8?style=for-the-badge)](https://gin-gonic.com)
[![Vue Version](https://img.shields.io/badge/Vue-3.4.0-4FC08D?style=for-the-badge&logo=vue.js)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](./LICENSE)

[![GitHub stars](https://img.shields.io/github/stars/xiaozhe2018/GinForge?style=social)](https://github.com/xiaozhe2018/GinForge/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/xiaozhe2018/GinForge?style=social)](https://github.com/xiaozhe2018/GinForge/network/members)

[快速开始](#-快速开始) • [在线演示](#) • [完整文档](./docs/INDEX.md) • [更新日志](./CHANGELOG.md)

</div>

---

## ✨ 为什么选择 GinForge？

```bash
# 一行命令，生成完整的CRUD功能（1000+行代码）
go run ./cmd/generator gen:crud --table=articles

# 自动生成：
# ✅ Model + Repository + Service + Handler (后端4层架构)
# ✅ TypeScript API + Vue列表页 + Vue表单页 (前端3个文件)
# ✅ 完整的增删改查 + 搜索 + 分页 + 排序
```

**开发效率提升 10 倍！**

</div>

## 🎯 核心亮点

| 特性 | 说明 | 优势 |
|-----|------|------|
| **⚡️ 一键生成CRUD** | 从数据库表自动生成全套代码 | **10分钟完成一个模块** |
| **🎨 完整后台管理** | Vue3 + Element Plus 现成可用 | **0前端开发，直接用** |
| **🔐 RBAC权限系统** | 用户-角色-权限-菜单四级控制 | **企业级权限方案** |
| **🏗️ 微服务架构** | 8个服务 + API网关 + Nginx | **生产环境直接部署** |
| **🐳 容器化部署** | Docker + K8s + Istio 配置齐全 | **一键上云** |
| **📚 详尽文档** | 46个文档 + 16个示例 | **0学习成本** |

## 💡 适用场景

✅ **企业内部管理系统** - 用户权限、数据管理、系统配置  
✅ **SaaS多租户平台** - 用户端 + 商户端 + 管理后台  
✅ **电商平台** - 商品、订单、支付、库存管理  
✅ **内容管理系统** - 文章、评论、分类、标签  
✅ **API服务** - RESTful API + Swagger文档自动生成  

## 🔥 技术栈

**后端：** Go 1.24 + Gin + GORM + JWT + Redis + WebSocket  
**前端：** Vue 3 + TypeScript + Element Plus + Pinia + Vite  
**数据库：** MySQL / PostgreSQL / SQLite  
**部署：** Docker + Docker Compose + Kubernetes + Istio  
**监控：** Prometheus + Grafana + Zap日志

## 🚀 快速开始

### 方式一：开发环境（30秒启动）

```bash
# 1. 克隆项目
git clone https://github.com/xiaozhe2018/GinForge.git
cd GinForge

# 2. 安装依赖
go mod tidy

# 3. 启动后端（使用SQLite，无需MySQL）
make run

# 4. 启动前端（新终端）
cd web/admin && npm install && npm run dev
```

**访问地址：**
- 🎉 管理后台：http://localhost:3000 （账号：`admin` / `admin123`）
- 📚 API文档：http://localhost:8083/swagger/index.html

### 方式二：Docker一键部署（推荐生产环境）

```bash
git clone https://github.com/xiaozhe2018/GinForge.git
cd GinForge/deployments
docker-compose up -d
```

**包含服务：** 后端API + MySQL + Redis + Nginx + 前端管理后台  
**访问地址：** http://localhost （账号：`admin` / `admin123`）

## 🎬 一键生成CRUD演示

```bash
# 1. 创建数据库表（或使用现有表）
# 例如：articles 文章表

# 2. 运行生成器
go run ./cmd/generator gen:crud --table=articles

# 3. 自动生成7个文件，共1000+行代码：
#   ✅ Model (数据模型)
#   ✅ Repository (数据访问层)  
#   ✅ Service (业务逻辑层)
#   ✅ Handler (HTTP处理层)
#   ✅ TypeScript API (前端接口)
#   ✅ Vue列表页 (带搜索/分页/排序)
#   ✅ Vue表单页 (新增/编辑)

# 4. 启动服务，功能立即可用！
```

**效果：** 10分钟完成一个功能模块的开发！

## 📸 界面预览

<details>
<summary>点击查看管理后台截图</summary>

### 登录页面
![登录](https://via.placeholder.com/800x450?text=Login+Page)

### 用户管理
![用户管理](https://via.placeholder.com/800x450?text=User+Management)

### 角色权限
![角色权限](https://via.placeholder.com/800x450?text=Role+Permission)

</details>

## 📚 文档

📖 [**完整文档**](./docs/INDEX.md) | ⚡ [快速开始](./docs/QUICK_START.md) | 🔍 [功能概览](./docs/FRAMEWORK_OVERVIEW.md) | 🚀 [高级功能](./docs/ADVANCED_FEATURES.md) | 💡 [示例代码](./docs/demo/)

## 🛠️ 功能特性

<details>
<summary><b>🏗️ 微服务架构（8个服务）</b></summary>

- `admin-api` (8083) - 管理后台API，RBAC权限系统
- `user-api` (8081) - 用户端API，用户信息管理
- `merchant-api` (8082) - 商户端API，商品订单管理
- `gateway` (8080) - API网关，统一入口
- `gateway-worker` (8084) - 网关工作器，异步任务
- `websocket-gateway` (8087) - WebSocket服务，实时通信
- `file-api` (8086) - 文件服务，上传下载
- `demo` (8085) - 演示服务，示例代码

</details>

<details>
<summary><b>🔐 安全认证</b></summary>

- JWT无状态认证 + Token黑名单
- Bcrypt密码加密
- RBAC四级权限控制（用户-角色-权限-菜单）
- 操作日志审计
- 登录失败锁定
- CORS跨域控制

</details>

<details>
<summary><b>📊 数据处理</b></summary>

- GORM ORM，支持 MySQL / PostgreSQL / SQLite
- 自动表结构迁移
- 软删除 + 硬删除
- 事务支持
- 关联查询（一对一、一对多、多对多）
- 分页 + 搜索 + 排序

</details>

<details>
<summary><b>⚡ 性能优化</b></summary>

- Redis多级缓存
- 数据库连接池
- 分布式锁
- 令牌桶限流
- 熔断器保护
- 异步任务处理

</details>

<details>
<summary><b>🎨 完整后台管理（17个页面）</b></summary>

- 🔐 登录/登出
- 📊 仪表盘（数据统计）
- 👥 用户管理（列表+表单+角色分配）
- 🔑 角色管理（列表+表单+权限配置）
- 📋 菜单管理（树形结构+图标选择）
- ⚙️ 权限管理（按钮级权限）
- 🧑 个人资料（信息+密码修改）
- 🔧 系统设置（配置+日志+监控）

</details>

<details>
<summary><b>🐳 容器化部署</b></summary>

- Docker多阶段构建
- Docker Compose编排
- Kubernetes配置
- Istio服务网格
- Nginx反向代理
- 环境变量配置
- 数据持久化

</details>

## 📁 项目结构

```
GinForge/
├── cmd/generator/          # 🔧 代码生成器
├── pkg/                    # 📦 共享基础库 (82个文件)
│   ├── middleware/         #    JWT、限流、缓存、日志等
│   ├── db/                 #    数据库管理
│   ├── redis/              #    缓存、队列、分布式锁
│   ├── config/             #    配置管理
│   └── ...
├── services/               # 🚀 微服务 (8个)
│   ├── admin-api/          #    管理后台API + RBAC权限
│   ├── user-api/           #    用户端API
│   ├── merchant-api/       #    商户端API
│   ├── gateway/            #    API网关
│   └── ...
├── web/admin/              # 🎨 Vue3管理后台 (17个页面)
├── deployments/            # 🐳 部署配置
│   ├── docker-compose.yml  #    Docker编排
│   ├── nginx.conf          #    Nginx配置
│   └── k8s/                #    K8s + Istio
├── docs/                   # 📚 完整文档 (46个文件)
└── database/migrations/    # 📊 数据库迁移
```

## 🎯 核心架构

**分层设计：** Handler → Service → Repository → Model  
**依赖注入：** 统一的服务注册和管理  
**中间件链：** Recovery + RequestID + CORS + JWT + 限流 + 日志 + 缓存  
**网关架构：** Nginx → Gateway → 后端微服务

## 📖 使用示例

### 快速生成CRUD

```bash
# 1. 列出所有数据库表
go run ./cmd/generator list:tables

# 2. 生成CRUD代码
go run ./cmd/generator gen:crud --table=articles

# 3. 生成的文件会自动放到正确位置
#    - services/admin-api/internal/model/
#    - services/admin-api/internal/repository/
#    - services/admin-api/internal/service/
#    - services/admin-api/internal/handler/
#    - web/admin/src/api/
#    - web/admin/src/views/Articles/
```

### 常用代码示例

<details>
<summary>Redis缓存</summary>

```go
cache := redis.NewCacheManager(redisClient)
cache.Set(ctx, "user:1001", userData, 5*time.Minute)
err := cache.Get(ctx, "user:1001", &user)
```

</details>

<details>
<summary>分布式锁</summary>

```go
lock := redis.NewDistributedLock(redisClient)
if lock.Acquire(ctx, "order:1001", 10*time.Second) {
    defer lock.Release(ctx, "order:1001")
    // 执行互斥操作
}
```

</details>

<details>
<summary>消息队列</summary>

```go
queue := redis.NewQueue(redisClient, "order-queue")
queue.Publish(ctx, map[string]interface{}{"order_id": 123})
queue.Subscribe(ctx, func(msg interface{}) error {
    // 处理消息
    return nil
})
```

</details>

更多示例请查看 [docs/demo/](./docs/demo/)

## 🔧 常用命令

```bash
# 开发
make run            # 启动所有后端服务
make stop           # 停止所有服务
make swagger        # 生成API文档

# 构建
make build          # 构建所有服务
make docker         # 构建Docker镜像

# 部署
make compose        # 启动Docker环境
make compose-down   # 停止Docker环境

# 测试
make test           # 运行所有测试
make test-coverage  # 生成覆盖率报告
```

## 🌟 Star History

如果这个项目对您有帮助，请给一个 ⭐️ Star 支持一下！

## 🤝 贡献

欢迎贡献代码、提交Issue、完善文档！

查看 [贡献指南](./CONTRIBUTING.md) 了解如何参与项目。

## 📄 开源协议

本项目采用 [MIT License](./LICENSE) 开源协议。

## 💬 交流群

- **GitHub Issues**: [提问题/建议](https://github.com/xiaozhe2018/GinForge/issues)
- **GitHub Discussions**: [技术讨论](https://github.com/xiaozhe2018/GinForge/discussions)

## 🙏 致谢

感谢所有贡献者和使用者的支持！

---

<div align="center">

**如果觉得不错，请点个 ⭐️ Star 支持一下，谢谢！**

Made with ❤️ by [xiaozhe2018](https://github.com/xiaozhe2018)

</div>