# 🚀 GinForge 快速使用指南

## 一、启动服务

### 启动后端
```bash
cd /Users/xiaozhe/go/goweb
make run
```

### 启动前端
```bash
cd /Users/xiaozhe/go/goweb/web/admin
npm run dev
```

### 停止所有服务
```bash
make stop
```

## 二、访问系统

- **前端管理后台**: http://localhost:3000
- **后端API**: http://localhost:8083
- **Swagger文档**: http://localhost:8083/swagger/index.html

### 登录信息
- 用户名: `admin`
- 密码: `admin123`

## 三、功能说明

### ✅ 已实现功能（真实API）
1. **登录/登出** - Token黑名单机制
2. **用户管理** - CRUD + 状态管理
3. **角色管理** - CRUD + 权限分配
4. **菜单管理** - 树形结构管理
5. **权限管理** - 权限CRUD
6. **个人设置** - 信息修改 + 密码修改

### ⚠️ 模拟数据（前端展示）
7. **仪表盘** - 统计图表
8. **系统管理** - 系统配置和日志

## 四、测试验证

### 测试登录登出
```bash
# 1. 登录
curl -X POST http://localhost:8083/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 2. 使用Token访问
TOKEN="your-token-here"
curl -X GET http://localhost:8083/api/v1/admin/auth/profile \
  -H "Authorization: Bearer $TOKEN"

# 3. 登出
curl -X POST http://localhost:8083/api/v1/admin/auth/logout \
  -H "Authorization: Bearer $TOKEN"

# 4. 再次使用Token（应该失败）
curl -X GET http://localhost:8083/api/v1/admin/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

### 查看操作日志
```bash
docker exec mysql mysql -uroot -p123456 gin_forge -e \
  "SELECT id, username, method, path, created_at FROM admin_operation_logs ORDER BY created_at DESC LIMIT 5;"
```

### 查看Token黑名单
```bash
docker exec redis redis-cli KEYS "token:blacklist:*"
```

## 五、常见问题

**Q: 登录失败？**  
A: 确认密码是`admin123`，检查后端服务是否运行

**Q: Network面板看不到请求？**  
A: 打开开发者工具，确保Network录制按钮是红色，过滤器选择Fetch/XHR

**Q: Token失效？**  
A: Token有效期24小时，或者已经登出（黑名单）

---

**让开发更加简单** 🎊
