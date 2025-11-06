<div align="center">

# GinForge

**企业级 Go 微服务开发框架**

*30秒启动，一键生成CRUD，开箱即用的管理后台*

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Gin Version](https://img.shields.io/badge/Gin-1.10.0-00ADD8?style=for-the-badge)](https://gin-gonic.com)
[![Vue Version](https://img.shields.io/badge/Vue-3.4.0-4FC08D?style=for-the-badge&logo=vue.js)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](./LICENSE)

[![GitHub stars](https://img.shields.io/github/stars/xiaozhe2018/GinForge?style=social)](https://github.com/xiaozhe2018/GinForge/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/xiaozhe2018/GinForge?style=social)](https://github.com/xiaozhe2018/GinForge/network/members)

[快速开始](#快速开始) • [完整文档](./docs/INDEX.md) • [更新日志](./CHANGELOG.md)

</div>

---

## 为什么选择 GinForge？

```bash
# 一行命令，生成完整的CRUD功能（1000+行代码）
go run ./cmd/generator gen:crud --table=articles

# 自动生成：Model + Repository + Service + Handler + TypeScript API + Vue页面
```

**开发效率提升 10 倍！**

## 核心特性

- **一键生成CRUD** - 从数据库表自动生成全套代码（10分钟完成一个模块）
- **完整后台管理** - Vue3 + Element Plus 现成可用（0前端开发）
- **RBAC权限系统** - 用户-角色-权限-菜单四级控制
- **微服务架构** - 8个服务 + API网关 + Nginx（生产环境直接部署）
- **容器化部署** - Docker + K8s + Istio 配置齐全

## 技术栈

**后端：** Go 1.24 + Gin + GORM + JWT + Redis + WebSocket  
**前端：** Vue 3 + TypeScript + Element Plus + Pinia + Vite  
**数据库：** MySQL / PostgreSQL / SQLite  
**部署：** Docker + Docker Compose + Kubernetes

## 快速开始

### 开发环境

```bash
# 1. 克隆项目
git clone https://github.com/xiaozhe2018/GinForge.git
cd GinForge

# 2. 安装依赖
go mod tidy
cd web/admin && npm install && cd ../..

# 3. 配置环境变量
cp env.example .env
# 编辑 .env，配置数据库连接信息

# 4. 创建数据库
mysql -uroot -p -e "CREATE DATABASE IF NOT EXISTS gin_forge DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 5. 初始化数据库
make init

# 6. 启动后端服务
make run

# 7. 启动前端（新终端）
cd web/admin && npm run dev
```

**访问地址：**
- 管理后台：http://localhost:3000 （账号：`admin` / `admin123`）
- API文档：http://localhost:8083/swagger/index.html

### Docker 部署

```bash
git clone https://github.com/xiaozhe2018/GinForge.git
cd GinForge/deployments
docker-compose up -d
```

**访问地址：** http://localhost （账号：`admin` / `admin123`）

## 一键生成CRUD

```bash
# 生成CRUD代码
go run ./cmd/generator gen:crud --table=articles

# 自动生成7个文件，共1000+行代码：
# - Model (数据模型)
# - Repository (数据访问层)
# - Service (业务逻辑层)
# - Handler (HTTP处理层)
# - TypeScript API (前端接口)
# - Vue列表页 (带搜索/分页/排序)
# - Vue表单页 (新增/编辑)
```

## 项目结构

```
GinForge/
├── cmd/generator/          # 代码生成器
├── pkg/                    # 共享基础库
│   ├── middleware/         # JWT、限流、缓存、日志等
│   ├── db/                 # 数据库管理
│   ├── redis/              # 缓存、队列、分布式锁
│   └── config/             # 配置管理
├── services/               # 微服务 (8个)
│   ├── admin-api/          # 管理后台API + RBAC权限
│   ├── user-api/           # 用户端API
│   ├── merchant-api/        # 商户端API
│   ├── gateway/            # API网关
│   └── ...
├── web/admin/              # Vue3管理后台
├── deployments/            # 部署配置
│   ├── docker-compose.yml  # Docker编排
│   └── k8s/                # K8s + Istio
└── database/migrations/     # 数据库迁移
```

## 常用命令

```bash
# 开发
make run            # 启动所有后端服务
make init           # 初始化数据库
make swagger        # 生成API文档

# 构建
make build          # 构建所有服务
make docker         # 构建Docker镜像

# 部署
make compose        # 启动Docker环境
```

## 文档

- [完整文档](./docs/INDEX.md)
- [管理后台指南](./docs/ADMIN_GUIDE.md)
- [快速开始](./docs/QUICK_START.md)
- [框架概览](./docs/FRAMEWORK_OVERVIEW.md)

## Star History

<a href="https://star-history.com/#xiaozhe2018/GinForge&Date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=xiaozhe2018/GinForge&type=Date&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=xiaozhe2018/GinForge&type=Date" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=xiaozhe2018/GinForge&type=Date" />
  </picture>
</a>

## 贡献

欢迎提 Issue 和 PR！查看 [贡献指南](./CONTRIBUTING.md) 了解详情。

## 开源协议

MIT License

---

<div align="center">

**做开源不易，点个⭐支持一下吧，谢谢！**

Made with ❤️ by [xiaozhe2018](https://github.com/xiaozhe2018)

</div>
