# GinForge 管理后台快速启动指南

## 🚀 5分钟快速体验

### 1. 安装依赖

```bash
# 进入前端目录
cd web/admin

# 安装依赖
npm install
```

### 2. 启动开发服务器

```bash
# 方式一：使用 npm
npm run dev

# 方式二：使用 make 命令（在项目根目录）
make web-dev
```

### 3. 访问管理后台

打开浏览器访问：http://localhost:3000

**默认登录账号：**
- 用户名：`admin`
- 密码：`123456`

## 🎯 功能演示

### 仪表盘
- 系统统计概览
- 最近订单列表
- 系统信息展示

### 用户管理
- 用户列表查看
- 用户状态管理
- 用户信息编辑

### 商户管理
- 商户列表查看
- 商户状态管理
- 商户信息维护

### 商品管理
- 商品列表查看
- 商品状态管理
- 商品分类筛选

### 订单管理
- 订单列表查看
- 订单状态跟踪
- 订单操作管理

### 系统设置
- 基本设置配置
- 邮件服务配置
- 安全策略设置
- 缓存配置管理
- 日志级别设置

## 🔧 开发指南

### 项目结构

```
web/admin/
├── src/
│   ├── api/           # API 接口层
│   ├── layout/        # 布局组件
│   ├── router/        # 路由配置
│   ├── views/         # 页面组件
│   ├── App.vue        # 根组件
│   └── main.ts        # 入口文件
├── public/            # 静态资源
├── package.json       # 项目配置
├── vite.config.ts     # Vite 配置
└── tsconfig.json      # TypeScript 配置
```

### 添加新页面

1. **创建页面组件**
   ```bash
   # 在 src/views/ 目录下创建新组件
   touch src/views/NewPage.vue
   ```

2. **配置路由**
   ```typescript
   // 在 src/router/index.ts 中添加路由
   {
     path: 'new-page',
     name: 'NewPage',
     component: () => import('@/views/NewPage.vue'),
     meta: { title: '新页面' }
   }
   ```

3. **添加菜单项**
   ```vue
   <!-- 在 src/layout/index.vue 中添加菜单 -->
   <el-menu-item index="/dashboard/new-page">
     <el-icon><NewIcon /></el-icon>
     <template #title>新页面</template>
   </el-menu-item>
   ```

### API 接口开发

1. **创建 API 文件**
   ```bash
   # 在 src/api/ 目录下创建新模块
   touch src/api/newModule.ts
   ```

2. **定义接口**
   ```typescript
   import api from './index'
   
   export interface NewData {
     id: string
     name: string
   }
   
   export const getNewData = () => {
     return api.get<NewData[]>('/v1/new-data')
   }
   ```

3. **在组件中使用**
   ```vue
   <script setup lang="ts">
   import { getNewData } from '@/api/newModule'
   
   const loadData = async () => {
     const data = await getNewData()
     // 处理数据
   }
   </script>
   ```

## 🎨 自定义主题

### 修改主题色

在 `src/main.ts` 中配置：

```typescript
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

app.use(ElementPlus, {
  locale: zhCn,
  // 自定义主题色
  theme: {
    primary: '#409EFF'
  }
})
```

### 自定义样式

在组件中使用 CSS 变量：

```vue
<style scoped>
.custom-button {
  --el-button-bg-color: #your-color;
  --el-button-border-color: #your-color;
}
</style>
```

## 📦 构建部署

### 开发构建

```bash
npm run build
```

### 生产部署

1. **构建项目**
   ```bash
   npm run build
   ```

2. **部署到服务器**
   ```bash
   # 将 dist/ 目录上传到服务器
   scp -r dist/* user@server:/var/www/admin/
   ```

3. **Nginx 配置**
   ```nginx
   server {
       listen 80;
       server_name admin.example.com;
       
       location / {
           root /var/www/admin;
           try_files $uri $uri/ /index.html;
       }
       
       location /api {
           proxy_pass http://localhost:8080;
       }
   }
   ```

## 🐛 常见问题

### Q: 页面空白或加载失败？
A: 检查后端服务是否启动，确保 API 接口可访问。

### Q: 登录后跳转失败？
A: 检查 localStorage 中是否有 `admin_token`，确保路由守卫配置正确。

### Q: 样式显示异常？
A: 确保 Element Plus 样式正确引入，检查是否有样式冲突。

### Q: API 请求失败？
A: 检查 `vite.config.ts` 中的代理配置，确保后端服务地址正确。

## 📞 技术支持

如有问题，请查看：
- [Vue 3 官方文档](https://vuejs.org/)
- [Element Plus 文档](https://element-plus.org/)
- [Vite 文档](https://vitejs.dev/)

---

**让开发更加简单！** 🚀


## 🚀 5分钟快速体验

### 1. 安装依赖

```bash
# 进入前端目录
cd web/admin

# 安装依赖
npm install
```

### 2. 启动开发服务器

```bash
# 方式一：使用 npm
npm run dev

# 方式二：使用 make 命令（在项目根目录）
make web-dev
```

### 3. 访问管理后台

打开浏览器访问：http://localhost:3000

**默认登录账号：**
- 用户名：`admin`
- 密码：`123456`

## 🎯 功能演示

### 仪表盘
- 系统统计概览
- 最近订单列表
- 系统信息展示

### 用户管理
- 用户列表查看
- 用户状态管理
- 用户信息编辑

### 商户管理
- 商户列表查看
- 商户状态管理
- 商户信息维护

### 商品管理
- 商品列表查看
- 商品状态管理
- 商品分类筛选

### 订单管理
- 订单列表查看
- 订单状态跟踪
- 订单操作管理

### 系统设置
- 基本设置配置
- 邮件服务配置
- 安全策略设置
- 缓存配置管理
- 日志级别设置

## 🔧 开发指南

### 项目结构

```
web/admin/
├── src/
│   ├── api/           # API 接口层
│   ├── layout/        # 布局组件
│   ├── router/        # 路由配置
│   ├── views/         # 页面组件
│   ├── App.vue        # 根组件
│   └── main.ts        # 入口文件
├── public/            # 静态资源
├── package.json       # 项目配置
├── vite.config.ts     # Vite 配置
└── tsconfig.json      # TypeScript 配置
```

### 添加新页面

1. **创建页面组件**
   ```bash
   # 在 src/views/ 目录下创建新组件
   touch src/views/NewPage.vue
   ```

2. **配置路由**
   ```typescript
   // 在 src/router/index.ts 中添加路由
   {
     path: 'new-page',
     name: 'NewPage',
     component: () => import('@/views/NewPage.vue'),
     meta: { title: '新页面' }
   }
   ```

3. **添加菜单项**
   ```vue
   <!-- 在 src/layout/index.vue 中添加菜单 -->
   <el-menu-item index="/dashboard/new-page">
     <el-icon><NewIcon /></el-icon>
     <template #title>新页面</template>
   </el-menu-item>
   ```

### API 接口开发

1. **创建 API 文件**
   ```bash
   # 在 src/api/ 目录下创建新模块
   touch src/api/newModule.ts
   ```

2. **定义接口**
   ```typescript
   import api from './index'
   
   export interface NewData {
     id: string
     name: string
   }
   
   export const getNewData = () => {
     return api.get<NewData[]>('/v1/new-data')
   }
   ```

3. **在组件中使用**
   ```vue
   <script setup lang="ts">
   import { getNewData } from '@/api/newModule'
   
   const loadData = async () => {
     const data = await getNewData()
     // 处理数据
   }
   </script>
   ```

## 🎨 自定义主题

### 修改主题色

在 `src/main.ts` 中配置：

```typescript
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

app.use(ElementPlus, {
  locale: zhCn,
  // 自定义主题色
  theme: {
    primary: '#409EFF'
  }
})
```

### 自定义样式

在组件中使用 CSS 变量：

```vue
<style scoped>
.custom-button {
  --el-button-bg-color: #your-color;
  --el-button-border-color: #your-color;
}
</style>
```

## 📦 构建部署

### 开发构建

```bash
npm run build
```

### 生产部署

1. **构建项目**
   ```bash
   npm run build
   ```

2. **部署到服务器**
   ```bash
   # 将 dist/ 目录上传到服务器
   scp -r dist/* user@server:/var/www/admin/
   ```

3. **Nginx 配置**
   ```nginx
   server {
       listen 80;
       server_name admin.example.com;
       
       location / {
           root /var/www/admin;
           try_files $uri $uri/ /index.html;
       }
       
       location /api {
           proxy_pass http://localhost:8080;
       }
   }
   ```

## 🐛 常见问题

### Q: 页面空白或加载失败？
A: 检查后端服务是否启动，确保 API 接口可访问。

### Q: 登录后跳转失败？
A: 检查 localStorage 中是否有 `admin_token`，确保路由守卫配置正确。

### Q: 样式显示异常？
A: 确保 Element Plus 样式正确引入，检查是否有样式冲突。

### Q: API 请求失败？
A: 检查 `vite.config.ts` 中的代理配置，确保后端服务地址正确。

## 📞 技术支持

如有问题，请查看：
- [Vue 3 官方文档](https://vuejs.org/)
- [Element Plus 文档](https://element-plus.org/)
- [Vite 文档](https://vitejs.dev/)

---

**让开发更加简单！** 🚀






