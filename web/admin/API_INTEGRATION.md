# 前后端API对接说明

## 🎉 对接完成情况

### ✅ 已完成
1. **登录认证** - Login.vue 已对接后端API
2. **用户管理** - user.ts API已定义，对接 `/api/v1/admin/users`
3. **角色管理** - role.ts API已定义，对接 `/api/v1/admin/roles`
4. **菜单管理** - menu.ts API已定义，对接 `/api/v1/admin/menus`
5. **权限管理** - permission.ts API已定义，对接 `/api/v1/admin/permissions`
6. **系统管理** - system.ts API已定义，对接 `/api/v1/admin/system`

## 🔑 登录信息

### 默认管理员账号
- **用户名**: `admin`
- **密码**: `admin123`
- **邮箱**: `admin@ginforge.com`

### 登录流程
1. 前端调用 `/api/v1/admin/auth/login` 接口
2. 后端验证用户名密码
3. 返回JWT Token、用户信息、菜单树、权限列表
4. 前端保存到 localStorage:
   - `admin_token` - JWT Token
   - `admin_user_info` - 用户信息
   - `admin_menus` - 菜单树
   - `admin_permissions` - 权限列表

## 📡 API端点

### 认证相关
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/admin/auth/login` | 用户登录 |
| POST | `/api/v1/admin/auth/logout` | 用户登出 |
| GET | `/api/v1/admin/auth/profile` | 获取当前用户信息 |
| PUT | `/api/v1/admin/auth/profile` | 更新当前用户信息 |
| PUT | `/api/v1/admin/auth/change-password` | 修改密码 |

### 用户管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/users` | 获取用户列表 |
| POST | `/api/v1/admin/users` | 创建用户 |
| GET | `/api/v1/admin/users/:id` | 获取用户详情 |
| PUT | `/api/v1/admin/users/:id` | 更新用户 |
| PUT | `/api/v1/admin/users/:id/status` | 更新用户状态 |
| DELETE | `/api/v1/admin/users/:id` | 删除用户 |

### 角色管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/roles` | 获取角色列表 |
| POST | `/api/v1/admin/roles` | 创建角色 |
| GET | `/api/v1/admin/roles/:id` | 获取角色详情 |
| PUT | `/api/v1/admin/roles/:id` | 更新角色 |
| DELETE | `/api/v1/admin/roles/:id` | 删除角色 |

### 菜单管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/menus` | 获取菜单列表 |
| GET | `/api/v1/admin/menus/tree` | 获取菜单树 |
| POST | `/api/v1/admin/menus` | 创建菜单 |
| GET | `/api/v1/admin/menus/:id` | 获取菜单详情 |
| PUT | `/api/v1/admin/menus/:id` | 更新菜单 |
| DELETE | `/api/v1/admin/menus/:id` | 删除菜单 |

### 权限管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/permissions` | 获取权限列表 |
| POST | `/api/v1/admin/permissions` | 创建权限 |
| GET | `/api/v1/admin/permissions/:id` | 获取权限详情 |
| PUT | `/api/v1/admin/permissions/:id` | 更新权限 |
| DELETE | `/api/v1/admin/permissions/:id` | 删除权限 |

## 🚀 启动服务

### 后端服务
```bash
# 确保MySQL在运行
docker ps | grep mysql

# 启动admin-api服务
cd /Users/xiaozhe/go/goweb
go run ./services/admin-api/cmd/server
```

服务运行在: `http://localhost:8083`

### 前端服务
```bash
cd /Users/xiaozhe/go/goweb/web/admin
npm run dev
```

服务运行在: `http://localhost:3000`

## 📝 API请求格式

### 请求头
```
Content-Type: application/json
Authorization: Bearer {token}  # 登录后需要
```

### 统一响应格式
```json
{
  "code": 0,           // 0表示成功，非0表示失败
  "message": "success", // 响应消息
  "data": {},          // 响应数据
  "trace_id": "xxx"    // 请求追踪ID
}
```

## 🔧 测试API

### 测试登录
```bash
curl -X POST http://localhost:8083/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 测试获取用户列表（需要token）
```bash
TOKEN="your-jwt-token-here"
curl -X GET "http://localhost:8083/api/v1/admin/users?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"
```

## 📊 前端API封装

所有API封装在 `src/api/` 目录下：
- `index.ts` - Axios实例配置，请求/响应拦截器
- `auth.ts` - 认证相关API
- `user.ts` - 用户管理API
- `role.ts` - 角色管理API
- `menu.ts` - 菜单管理API
- `permission.ts` - 权限管理API
- `system.ts` - 系统管理API

### 使用示例
```typescript
import { login } from '@/api/auth'
import { getUserList } from '@/api/user'

// 登录
const result = await login({ username: 'admin', password: 'admin123' })
localStorage.setItem('admin_token', result.token)

// 获取用户列表
const users = await getUserList({ page: 1, page_size: 10 })
```

## ⚠️ 注意事项

1. **CORS配置**: 后端已配置CORS允许前端访问
2. **Token过期**: JWT Token有效期24小时，过期后需重新登录
3. **权限控制**: 所有需要认证的接口都需要在请求头中携带Token
4. **错误处理**: 前端已配置统一的错误处理，401会自动跳转到登录页

## 🎯 下一步工作

1. 完善Dashboard页面的数据展示
2. 测试各个管理页面的CRUD操作
3. 添加更多的交互反馈和加载状态
4. 完善错误处理和边界情况
5. 添加更多的数据校验

## 📞 问题排查

### 登录失败
- 检查用户名密码是否正确 (admin/admin123)
- 检查后端服务是否运行 (localhost:8083)
- 检查数据库连接是否正常

### API请求失败
- 检查网络连接
- 检查Token是否有效
- 查看浏览器控制台错误信息
- 查看后端日志

---

**最后更新**: 2025-01-11

