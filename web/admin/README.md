# GinForge 管理后台

基于 Vue3 + Element Plus 的现代化管理后台界面。

## 🚀 快速开始

### 环境要求
- Node.js 16+
- npm 或 yarn

### 安装依赖
```bash
cd web/admin
npm install
```

### 开发模式
```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本
```bash
npm run build
```

## 📁 项目结构

```
web/admin/
├── src/
│   ├── api/           # API接口
│   ├── layout/        # 布局组件
│   ├── router/        # 路由配置
│   ├── views/         # 页面组件
│   ├── App.vue        # 根组件
│   └── main.ts        # 入口文件
├── public/            # 静态资源
├── package.json       # 依赖配置
└── vite.config.ts     # Vite配置
```

## 🎯 功能特性

- **现代化UI**: 基于 Element Plus 的专业管理界面
- **响应式设计**: 支持桌面端和移动端
- **权限管理**: 完整的用户、角色、菜单、权限管理体系
- **系统监控**: 实时系统状态监控和日志管理
- **个人设置**: 用户个人信息和账户安全设置

## 🔧 技术栈

- **Vue 3**: 渐进式JavaScript框架
- **Element Plus**: Vue 3 UI组件库
- **Vite**: 快速构建工具
- **TypeScript**: 类型安全的JavaScript
- **Vue Router**: 官方路由管理器
- **Pinia**: 状态管理库

## 📝 开发指南

### 添加新页面
1. 在 `src/views/` 下创建页面组件
2. 在 `src/router/index.ts` 中添加路由
3. 在 `src/layout/index.vue` 中添加菜单项

### API接口
- 所有API接口定义在 `src/api/` 目录下
- 使用统一的请求拦截器处理认证和错误

### 样式规范
- 使用 Element Plus 组件库
- 遵循 BEM CSS 命名规范
- 支持响应式设计

## 🚀 部署

### 构建
```bash
npm run build
```

### 部署到Nginx
将 `dist` 目录内容复制到Nginx静态文件目录，配置反向代理到后端API。

## 📞 支持

如有问题，请查看项目文档或提交Issue。


