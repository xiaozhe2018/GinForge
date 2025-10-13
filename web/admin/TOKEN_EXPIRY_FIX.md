# Token过期自动跳转功能说明

## 🎯 问题描述
用户反馈：当token过期后，前端虽然收到401错误，但没有自动跳转到登录页面。

## ✅ 解决方案

### 1. 问题分析
后端返回401时，响应体格式为：
```json
{
  "code": 401,
  "message": "认证令牌解析失败: token has invalid claims: token is expired",
  "trace_id": "..."
}
```

但是Axios响应拦截器中，这种情况会进入成功回调（因为HTTP请求成功了），而不是错误回调。

### 2. 修复内容

#### 文件：`web/admin/src/api/index.ts`

**改进点：**
1. **在响应成功回调中检查 `code === 401`**
   ```typescript
   if (data.code === 401) {
     // 处理token过期
     if (!isRedirectingToLogin) {
       isRedirectingToLogin = true
       ElMessage.error('登录已过期，请重新登录')
       // 清除所有本地存储数据
       localStorage.removeItem('admin_token')
       localStorage.removeItem('admin_user_info')
       localStorage.removeItem('admin_permissions')
       localStorage.removeItem('admin_menus')
       // 使用Vue Router跳转
       router.push('/login')
     }
     return Promise.reject(new Error(data.message || '认证失败'))
   }
   ```

2. **添加防重复跳转机制**
   - 使用 `isRedirectingToLogin` 标志避免多个401请求导致多次跳转和弹窗
   - 在登录页挂载时重置该标志

3. **同时保留HTTP 401状态码的处理**
   - 兼容不同的401响应方式

#### 文件：`web/admin/src/views/Login.vue`

**改进点：**
- 在组件挂载时重置重定向标志，确保用户可以重新登录

### 3. 完整流程

```
Token过期 → 发起请求 → 后端返回401
                       ↓
            Axios响应拦截器检测到code=401
                       ↓
            只弹出一次提示（防重复机制）
                       ↓
            清除所有本地存储数据
                       ↓
            使用Vue Router跳转到/login
                       ↓
            登录页挂载，重置重定向标志
```

## 🧪 测试方法

### 方法1：使用过期的Token
```bash
# 使用一个已过期的token
curl 'http://localhost:8083/api/v1/admin/users?page=1&page_size=20' \
  -H 'Authorization: Bearer <expired-token>'
```

**预期结果：**
- 前端弹出提示："登录已过期，请重新登录"
- 自动跳转到登录页 `/login`
- 本地存储被清空

### 方法2：在浏览器中测试
1. 登录系统，获取token
2. 等待token过期（24小时），或者手动修改localStorage中的token为过期的token
3. 刷新页面或点击任何需要认证的功能
4. 观察是否自动跳转到登录页

### 方法3：模拟快速测试
1. 修改后端JWT过期时间为1分钟（仅测试用）
2. 登录系统
3. 等待1分钟后刷新页面
4. 观察是否自动跳转

## 📋 改进清单

- ✅ 在响应成功回调中检查 `code === 401`
- ✅ 在响应错误回调中检查 HTTP `status === 401`
- ✅ 添加防重复跳转机制
- ✅ 清除所有本地存储数据（token、user_info、permissions、menus）
- ✅ 使用Vue Router进行优雅跳转
- ✅ 在登录页重置重定向标志
- ✅ 提供友好的错误提示

## 🎉 功能验证

用户测试场景：
```
1. 用户登录 ✅
2. Token过期 ✅
3. 发起任何API请求 ✅
4. 收到401错误 ✅
5. 弹出提示"登录已过期，请重新登录" ✅
6. 自动跳转到登录页 ✅
7. 本地数据被清空 ✅
8. 用户可以重新登录 ✅
```

## 💡 技术亮点

1. **双重401检查**：同时处理响应体code和HTTP status
2. **防抖机制**：避免多个请求导致的重复跳转
3. **优雅降级**：Vue Router失败时fallback到window.location
4. **完整清理**：清除所有相关的本地存储数据
5. **状态重置**：登录页挂载时重置标志位

---

**最后更新**: 2025-10-13  
**状态**: ✅ 已完成并测试

