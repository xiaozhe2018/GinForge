# 认证授权

GinForge 提供了完整的认证授权系统，包括 JWT 认证和 RBAC 权限控制。

## 🔐 认证系统

### JWT 认证流程

```
1. 用户登录
   ↓
2. 验证用户名密码
   ↓
3. 生成 JWT Token
   ↓
4. 返回 Token 给客户端
   ↓
5. 客户端保存 Token
   ↓
6. 后续请求携带 Token
   ↓
7. 服务器验证 Token
   ↓
8. 返回数据
```

## 📝 登录实现

### 后端登录 API

```go
// handler/admin_auth_handler.go
func (h *AdminAuthHandler) Login(c *gin.Context) {
    var req model.AdminUserLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "参数错误")
        return
    }
    
    // 获取客户端 IP
    loginIP := c.ClientIP()
    
    // 调用服务层登录
    result, err := h.userService.Login(&req, loginIP)
    if err != nil {
        response.Error(c, 401, err.Error())
        return
    }
    
    response.Success(c, result)
}
```

### 登录 Service

```go
// service/admin_service.go
func (s *UserService) Login(req *model.AdminUserLoginRequest, loginIP string) (*model.AdminUserLoginResponse, error) {
    // 1. 查询用户
    user, err := s.userRepo.GetByUsername(req.Username)
    if err != nil {
        return nil, errors.New("用户名或密码错误")
    }
    
    // 2. 验证密码
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, errors.New("用户名或密码错误")
    }
    
    // 3. 检查用户状态
    if user.Status != 1 {
        return nil, errors.New("用户已被禁用")
    }
    
    // 4. 生成 JWT Token
    claims := jwt.MapClaims{
        "user_id":  fmt.Sprintf("%d", user.ID),
        "username": user.Username,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(s.config.GetString("jwt.secret")))
    if err != nil {
        return nil, errors.New("生成令牌失败")
    }
    
    // 5. 更新登录信息
    s.userRepo.UpdateLoginInfo(user.ID, loginIP)
    
    // 6. 获取用户菜单和权限
    menus, _ := s.menuRepo.GetTreeByRoleID(user.Roles[0].ID)
    permissions, _ := s.permissionRepo.GetCodesByUserID(user.ID)
    
    return &model.AdminUserLoginResponse{
        Token:       tokenString,
        User:        user,
        Menus:       menus,
        Permissions: permissions,
    }, nil
}
```

## 🔑 JWT 中间件

### 基础 JWT 认证

```go
// 添加 JWT 中间件
auth := r.Group("/api/v1")
auth.Use(middleware.JWTAuth(jwtSecret))
{
    auth.GET("/profile", getProfile)
}
```

### 带 Redis 的 JWT 认证

支持 Token 黑名单和登录失败锁定：

```go
auth.Use(middleware.JWTAuthWithRedis(jwtSecret, redisClient))
```

**优势**：
- ✅ 支持 Token 黑名单（用户退出登录）
- ✅ 支持账户锁定（登录失败次数限制）
- ✅ 支持单点登录控制

### 在 Handler 中获取用户信息

```go
func (h *UserHandler) GetProfile(c *gin.Context) {
    // JWT 中间件会自动设置这些值
    userID := c.GetString("user_id")
    username := c.GetString("username")
    
    log.Info("user accessing profile", "user_id", userID)
    
    // 查询用户信息
    user, err := h.userService.GetUserByID(userID)
    if err != nil {
        response.InternalError(c, "获取用户信息失败")
        return
    }
    
    response.Success(c, user)
}
```

## 🎭 RBAC 权限控制

### 数据模型

```go
// User 用户
type User struct {
    ID    uint64
    Name  string
    Roles []Role `gorm:"many2many:user_roles;"`
}

// Role 角色
type Role struct {
    ID          uint64
    Name        string
    Permissions []Permission `gorm:"many2many:role_permissions;"`
}

// Permission 权限
type Permission struct {
    ID   uint64
    Name string
    Code string  // 例如："user:read", "user:create"
}
```

### 权限检查

```go
// 检查用户是否有某个权限
func (s *UserService) HasPermission(userID uint64, permissionCode string) bool {
    permissions, err := s.permissionRepo.GetCodesByUserID(userID)
    if err != nil {
        return false
    }
    
    for _, code := range permissions {
        if code == permissionCode {
            return true
        }
    }
    
    return false
}
```

### 权限中间件

```go
// 检查权限的中间件
func PermissionMiddleware(required string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        
        // 从 Context 或数据库获取用户权限
        permissions, _ := c.Get("permissions")
        permList, ok := permissions.([]string)
        if !ok {
            c.JSON(403, gin.H{"error": "权限不足"})
            c.Abort()
            return
        }
        
        // 检查权限
        hasPermission := false
        for _, perm := range permList {
            if perm == required {
                hasPermission = true
                break
            }
        }
        
        if !hasPermission {
            c.JSON(403, gin.H{"error": "权限不足"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// 使用
auth.GET("/users", 
    PermissionMiddleware("user:read"),
    getUserList,
)
```

## 🔄 登出功能

### 实现登出

```go
// handler/admin_auth_handler.go
func (h *AdminAuthHandler) Logout(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        response.Success(c, gin.H{"message": "退出成功"})
        return
    }
    
    // 移除 "Bearer " 前缀
    token = strings.TrimPrefix(token, "Bearer ")
    
    // 将 Token 加入黑名单（24小时）
    if h.redisClient != nil && h.redisClient.IsEnabled() {
        key := fmt.Sprintf("token:blacklist:%s", token)
        h.redisClient.Set(context.Background(), key, "1", 24*time.Hour)
    }
    
    response.Success(c, gin.H{"message": "退出成功"})
}
```

### 前端登出

```typescript
// api/auth.ts
export const logout = () => {
  return request.post('/api/v1/admin/logout')
}

// 在组件中使用
const handleLogout = async () => {
  await logout()
  localStorage.clear()  // 清除本地存储
  router.push('/login')
}
```

## 🔐 密码安全

### 密码加密

```go
import "golang.org/x/crypto/bcrypt"

// 加密密码
func HashPassword(password string) (string, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashed), err
}

// 验证密码
func CheckPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
```

### 密码策略

```go
// 在系统配置中设置密码策略
type PasswordPolicy struct {
    MinLength   int      // 最小长度
    Complexity  []string // 复杂度要求：lowercase, uppercase, numbers, symbols
}

// 验证密码复杂度
func (s *AdminSystemService) ValidatePassword(password string) error {
    minLength := s.GetPasswordMinLength()  // 从配置读取
    complexity := s.GetPasswordComplexity()
    
    return validator.ValidatePassword(password, minLength, complexity)
}
```

## 🚫 账户锁定

### 登录失败锁定

```go
func (s *UserService) Login(req *model.AdminUserLoginRequest, loginIP string) (*model.AdminUserLoginResponse, error) {
    ctx := context.Background()
    
    // 检查账户是否被锁定
    lockKey := fmt.Sprintf("login:locked:%s", req.Username)
    locked, err := s.redisClient.Exists(ctx, lockKey)
    if err == nil && locked {
        ttl, _ := s.redisClient.TTL(ctx, lockKey).Result()
        return nil, fmt.Errorf("账户已被锁定，请在 %d 分钟后重试", int(ttl.Minutes())+1)
    }
    
    // 验证密码...
    if err := bcrypt.CompareHashAndPassword(...); err != nil {
        // 记录失败次数
        s.recordLoginFailure(ctx, req.Username)
        return nil, errors.New("用户名或密码错误")
    }
    
    // 登录成功，清除失败记录
    s.clearLoginFailures(ctx, req.Username)
    
    // 生成 Token...
    return result, nil
}

// 记录登录失败
func (s *UserService) recordLoginFailure(ctx context.Context, username string) {
    key := fmt.Sprintf("login:failures:%s", username)
    
    // 增加失败次数
    failures := s.redisClient.Incr(ctx, key).Val()
    
    // 设置过期时间（15分钟）
    if failures == 1 {
        s.redisClient.Expire(ctx, key, 15*time.Minute)
    }
    
    // 失败次数超过限制，锁定账户
    maxAttempts := s.systemService.GetMaxLoginAttempts()  // 默认 5 次
    if failures >= int64(maxAttempts) {
        lockKey := fmt.Sprintf("login:locked:%s", username)
        lockDuration := s.systemService.GetLockoutDuration()  // 默认 30 分钟
        s.redisClient.Set(ctx, lockKey, "1", time.Duration(lockDuration)*time.Minute)
        s.logger.Warn("account locked due to too many failed attempts", "username", username)
    }
}
```

## 🎯 权限管理示例

### 创建角色

```go
role := &model.AdminRole{
    Name:        "管理员",
    Description: "系统管理员",
    Status:      1,
}
db.Create(&role)
```

### 分配权限给角色

```go
// 多对多关联
permissions := []model.AdminPermission{
    {ID: 1, Code: "user:read"},
    {ID: 2, Code: "user:create"},
}

// 使用 Association 添加权限
db.Model(&role).Association("Permissions").Append(&permissions)
```

### 分配角色给用户

```go
roles := []model.AdminRole{
    {ID: 1, Name: "管理员"},
}

db.Model(&user).Association("Roles").Append(&roles)
```

### 检查用户权限

```go
// 获取用户的所有权限代码
func (r *PermissionRepository) GetCodesByUserID(userID uint64) ([]string, error) {
    var codes []string
    
    err := r.db.Model(&model.AdminPermission{}).
        Select("admin_permissions.code").
        Joins("JOIN admin_role_permissions ON admin_permissions.id = admin_role_permissions.permission_id").
        Joins("JOIN admin_user_roles ON admin_role_permissions.role_id = admin_user_roles.role_id").
        Where("admin_user_roles.user_id = ?", userID).
        Pluck("admin_permissions.code", &codes).Error
    
    return codes, err
}
```

## 🎨 前端集成

### 登录请求

```typescript
// api/auth.ts
export const login = (data: LoginParams) => {
  return request.post<LoginResponse>('/api/v1/admin/login', data)
}

// 使用
const handleLogin = async () => {
  const result = await login({
    username: 'admin',
    password: 'admin123'
  })
  
  // 保存 Token
  localStorage.setItem('admin_token', result.token)
  localStorage.setItem('admin_user_info', JSON.stringify(result.user))
  localStorage.setItem('admin_permissions', JSON.stringify(result.permissions))
  
  // 跳转到仪表盘
  router.push('/dashboard')
}
```

### 请求拦截器

```typescript
// api/index.ts
import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// 请求拦截器：自动添加 Token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('admin_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器：处理 401
request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      // Token 过期，跳转到登录页
      localStorage.clear()
      router.push('/login')
    }
    return Promise.reject(error)
  }
)
```

### 路由守卫

```typescript
// router/index.ts
router.beforeEach((to, from, next) => {
  // 检查是否需要认证
  if (to.meta.requiresAuth !== false) {
    const token = localStorage.getItem('admin_token')
    if (!token) {
      next('/login')
      return
    }
  }
  
  // 检查权限
  if (to.meta.permission) {
    const permissions = JSON.parse(localStorage.getItem('admin_permissions') || '[]')
    if (!permissions.includes(to.meta.permission)) {
      ElMessage.error('没有权限访问此页面')
      next('/dashboard')
      return
    }
  }
  
  next()
})
```

## 🔒 会话管理

### 会话超时

```go
// 动态会话超时（从系统配置读取）
sessionTimeout := s.systemService.GetSessionTimeout(ctx)  // 默认 24 小时

claims := jwt.MapClaims{
    "user_id":  userID,
    "username": username,
    "exp":      time.Now().Add(time.Duration(sessionTimeout) * time.Minute).Unix(),
}
```

### 刷新 Token

```go
func (h *AuthHandler) RefreshToken(c *gin.Context) {
    // 从当前 Token 获取用户信息
    userID := c.GetString("user_id")
    username := c.GetString("username")
    
    // 生成新 Token
    claims := jwt.MapClaims{
        "user_id":  userID,
        "username": username,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString([]byte(jwtSecret))
    
    response.Success(c, gin.H{"token": tokenString})
}
```

## 🎯 权限指令（前端）

### Vue 自定义指令

```typescript
// directives/permission.ts
export const permission = {
  mounted(el: HTMLElement, binding: any) {
    const { value } = binding
    const permissions = JSON.parse(localStorage.getItem('admin_permissions') || '[]')
    
    if (!permissions.includes(value)) {
      // 移除没有权限的元素
      el.parentNode?.removeChild(el)
    }
  }
}

// 注册指令
app.directive('permission', permission)
```

### 使用权限指令

```vue
<template>
  <!-- 只有拥有 user:create 权限的用户才能看到这个按钮 -->
  <el-button v-permission="'user:create'" @click="handleCreate">
    创建用户
  </el-button>
</template>
```

## 📊 完整示例

### 后端完整流程

1. **定义路由**：`services/admin-api/internal/router/router.go`
2. **登录 Handler**：`services/admin-api/internal/handler/admin_auth_handler.go`
3. **登录 Service**：`services/admin-api/internal/service/admin_service.go`
4. **JWT 中间件**：`pkg/middleware/jwt.go`

### 前端完整流程

1. **登录 API**：`web/admin/src/api/auth.ts`
2. **登录页面**：`web/admin/src/views/Login.vue`
3. **路由守卫**：`web/admin/src/router/index.ts`
4. **请求拦截**：`web/admin/src/api/index.ts`

## 💡 最佳实践

### 1. 密码安全

✅ 使用 bcrypt 加密密码  
✅ 设置最小密码长度和复杂度  
✅ 定期提示用户更新密码  
✅ 禁止弱密码和常见密码  

### 2. Token 安全

✅ 使用 HTTPS 传输  
✅ Token 设置合理的过期时间  
✅ 支持 Token 刷新机制  
✅ 退出登录时加入黑名单  

### 3. 权限设计

✅ 最小权限原则  
✅ 权限粒度要合理  
✅ 支持动态权限更新  
✅ 记录权限变更日志  

### 4. 安全防护

✅ 登录失败次数限制  
✅ 账户异常登录提醒  
✅ IP 白名单（可选）  
✅ 双因素认证（可选）  

## 🎯 下一步

- [文件上传](./file-upload) - 学习文件上传功能
- [缓存系统](./cache) - 使用 Redis 缓存
- [WebSocket](./websocket) - 实时通信

---

**提示**: 安全是系统的生命线，务必认真对待认证和授权的实现！

