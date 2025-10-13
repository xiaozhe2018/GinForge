# 前端API对接状态

## ✅ 已完成对接（真实API）

### 1. 登录认证模块 (Login.vue)
- ✅ **登录**: `POST /api/v1/admin/auth/login`
- ✅ **登出**: `POST /api/v1/admin/auth/logout`
- ✅ **Token黑名单机制**: Redis实现

### 2. 用户管理模块 (Users.vue)
- ✅ **获取用户列表**: `GET /api/v1/admin/users`
- ✅ **创建用户**: `POST /api/v1/admin/users`
- ✅ **更新用户**: `PUT /api/v1/admin/users/:id`
- ✅ **删除用户**: `DELETE /api/v1/admin/users/:id`
- ✅ **更新状态**: `PUT /api/v1/admin/users/:id/status`

### 3. 角色管理模块 (Roles.vue)
- ✅ **获取角色列表**: `GET /api/v1/admin/roles`
- ✅ **创建角色**: `POST /api/v1/admin/roles`
- ✅ **更新角色**: `PUT /api/v1/admin/roles/:id`
- ✅ **删除角色**: `DELETE /api/v1/admin/roles/:id`
- ✅ **获取角色详情**: `GET /api/v1/admin/roles/:id`
- ✅ **权限树**: `GET /api/v1/admin/menus/tree`

### 4. 菜单管理模块 (Menus.vue)
- ✅ **获取菜单树**: `GET /api/v1/admin/menus/tree`
- ✅ **创建菜单**: `POST /api/v1/admin/menus`
- ✅ **更新菜单**: `PUT /api/v1/admin/menus/:id`
- ✅ **删除菜单**: `DELETE /api/v1/admin/menus/:id`

### 5. 权限管理模块 (Permissions.vue) 🆕
- ✅ **获取权限列表**: `GET /api/v1/admin/permissions`
- ✅ **创建权限**: `POST /api/v1/admin/permissions`
- ✅ **更新权限**: `PUT /api/v1/admin/permissions/:id`
- ✅ **删除权限**: `DELETE /api/v1/admin/permissions/:id`

### 6. 个人设置模块 (Profile.vue) 🆕
- ✅ **获取个人信息**: `GET /api/v1/admin/auth/profile`
- ✅ **更新个人信息**: `PUT /api/v1/admin/auth/profile`
- ✅ **修改密码**: `PUT /api/v1/admin/auth/change-password`

## 📋 使用模拟数据（可选）

### 7. 仪表盘模块 (Dashboard.vue)
- ⚠️ 使用模拟数据
- 可选接口: 系统统计API（后端未实现）

### 8. 系统管理模块 (System.vue)
- ⚠️ 使用模拟数据
- 可选接口: 系统配置API、系统日志API（后端未实现）

## 🎯 数据格式转换说明

### 后端 → 前端字段映射

#### 菜单(Menu)
```typescript
后端字段 → 前端字段
parent_id → parentId
created_at → createdAt
updated_at → updatedAt
status: 1/0 → status: 'show'/'hide'
visible: 1/0 → visible: 1/0
code → permission (权限标识)
```

#### 角色(Role)
```typescript
后端字段 → 前端字段
created_at → createdAt
status: 1/0 → status: 'active'/'disabled'
```

#### 用户(User)
```typescript
后端字段 → 前端字段
last_login_at → last_login_at
last_login_ip → last_login_ip
created_at → created_at
status: 1/0 → status: 1/0
```

## 📝 下一步工作

1. ✅ ~~菜单管理接入真实API~~
2. ⏳ 权限管理接入真实API
3. ⏳ 个人设置接入真实API
4. ⏳ Dashboard数据可视化（可选）
5. ⏳ 系统管理接入真实API（可选）

## 🔧 测试方法

### 测试菜单管理
```bash
# 访问菜单管理页面
http://localhost:3000/dashboard/menus

# 应该能看到真实的菜单数据（来自数据库）
# 可以进行创建、编辑、删除操作
```

### 测试角色管理
```bash
# 访问角色管理页面
http://localhost:3000/dashboard/roles

# 应该能看到真实的角色数据
# 点击"权限"按钮，可以看到菜单树
```

---

**最后更新**: 2025-10-11
**对接进度**: 6/8 模块完成 (75%)

