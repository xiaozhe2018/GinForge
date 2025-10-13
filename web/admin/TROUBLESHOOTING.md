# 🔍 前端调试指南

## Network面板看不到API请求？

### 检查步骤

#### 1️⃣ 打开浏览器开发者工具
- 按 `F12` 或 `Cmd+Option+I` (Mac)
- 点击 `Network` 标签

#### 2️⃣ 确保Network面板正在记录
- Network面板左上角应该有一个红色的录制按钮 🔴
- 确保它是激活状态（红色）
- 如果是灰色，点击一下激活录制

#### 3️⃣ 清除过滤器
- 检查Network面板上方的过滤器
- 确保选择 `All` 或 `Fetch/XHR`
- 清空搜索框

#### 4️⃣ 执行登录操作
```
1. 刷新页面 (Cmd+R 或 F5)
2. 在登录页面输入:
   用户名: admin
   密码: admin123
3. 点击"登录"按钮
4. 观察Network面板
```

### 应该看到的请求

#### 登录请求
```
Method: POST
URL: http://localhost:3000/api/v1/admin/auth/login
Status: 200
Response: JSON with token, user, menus, permissions
```

#### 请求详情
- **Request Headers**:
  - Content-Type: application/json
- **Request Body**:
  ```json
  {
    "username": "admin",
    "password": "admin123"
  }
  ```
- **Response**:
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "token": "eyJhbGci...",
      "user": {...},
      "menus": [...],
      "permissions": [...]
    }
  }
  ```

## 📊 Vite代理配置

前端配置了代理，所有 `/api` 请求会被转发到后端：

```typescript
// vite.config.ts
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8083',  // 后端地址
      changeOrigin: true,
      rewrite: (path) => path
    }
  }
}
```

这意味着：
- 前端请求: `http://localhost:3000/api/v1/admin/auth/login`
- 实际转发到: `http://localhost:8083/api/v1/admin/auth/login`

## 🔧 常见问题

### 1. Network面板完全空白
**原因**: 可能没有激活录制  
**解决**: 点击左上角的录制按钮 🔴

### 2. 只看到页面资源，没有API请求
**原因**: 还没有执行登录操作  
**解决**: 输入用户名密码，点击登录按钮

### 3. 看到请求但是失败 (红色)
**原因**: 后端服务没有运行或端口不对  
**解决**: 
```bash
# 检查后端服务
ps aux | grep "admin-api"

# 如果没有运行，启动它
go run ./services/admin-api/cmd/server
```

### 4. 请求显示CORS错误
**原因**: 跨域配置问题  
**解决**: 后端已配置CORS，确认后端服务正在运行

### 5. 请求超时
**原因**: 数据库连接问题  
**解决**:
```bash
# 检查MySQL
docker ps | grep mysql

# 检查配置
cat configs/config.yaml | grep -A 10 database
```

## 🧪 测试API直接访问

如果Network面板看不到请求，可以直接测试后端API：

### 方法1: 使用curl
```bash
curl -X POST http://localhost:8083/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 方法2: 使用浏览器
直接访问: `http://localhost:8083/healthz`  
应该看到: `{"status":"ok","service":"admin-api"}`

### 方法3: 使用Postman
1. Method: POST
2. URL: `http://localhost:8083/api/v1/admin/auth/login`
3. Headers: `Content-Type: application/json`
4. Body:
   ```json
   {
     "username": "admin",
     "password": "admin123"
   }
   ```

## 🎯 完整调试流程

### Step 1: 确认服务状态
```bash
# 后端服务
ps aux | grep "admin-api"
curl http://localhost:8083/healthz

# 前端服务  
ps aux | grep "vite"
curl http://localhost:3000
```

### Step 2: 打开浏览器调试
1. 访问: `http://localhost:3000`
2. 按 F12 打开开发者工具
3. 切换到 Console 标签，查看是否有错误
4. 切换到 Network 标签，激活录制
5. 确保过滤器设置为 "All" 或 "Fetch/XHR"

### Step 3: 执行登录
1. 在登录页面输入用户名和密码
2. 打开 Network 标签
3. 点击"登录"按钮
4. 观察请求列表

### Step 4: 查看请求详情
点击 `login` 请求，查看：
- **Headers** 标签: 请求头和响应头
- **Payload** 标签: 请求体
- **Response** 标签: 响应内容
- **Preview** 标签: 格式化的响应

## 📸 预期截图位置

在Network面板中应该看到：

```
Name                          Status  Type    Size
----------------------------------------------
login                         200     xhr     4.7KB
```

点击后在 Preview 标签应该看到：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJ...",
    "user": { ... },
    "menus": [ ... ],
    "permissions": [ ... ]
  }
}
```

## 💡 提示

1. **刷新页面**: 每次测试前刷新页面，清空之前的请求
2. **保留日志**: 勾选 "Preserve log" 可以保留页面跳转前的请求
3. **禁用缓存**: 勾选 "Disable cache" 确保获取最新数据
4. **查看Console**: 同时查看Console标签，可能有错误信息
5. **检查源码**: 在 Sources 标签可以查看和调试源代码

---

**如果还是看不到请求，请截图给我看看！** 📸

