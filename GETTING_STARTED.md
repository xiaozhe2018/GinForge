# 🚀 GinForge 快速上手指南

## 一分钟快速开始

### 1️⃣ 检查环境

```bash
# 确认MySQL正在运行
docker ps | grep mysql
```

### 2️⃣ 启动后端

```bash
# 进入项目目录
cd /Users/xiaozhe/go/goweb

# 启动管理后台API服务
go run ./services/admin-api/cmd/server
```

✅ 后端服务运行在: **http://localhost:8083**

### 3️⃣ 启动前端

```bash
# 进入前端目录
cd web/admin

# 首次启动需要安装依赖
npm install

# 启动开发服务器
npm run dev
```

✅ 前端服务运行在: **http://localhost:3000**

### 4️⃣ 登录系统

打开浏览器访问: **http://localhost:3000**

```
用户名: admin
密码: admin123
```

## 🎯 核心功能

登录后你可以：

- ✅ **仪表盘** - 查看系统概况和统计信息
- ✅ **用户管理** - 管理系统用户，分配角色
- ✅ **角色管理** - 配置角色和权限
- ✅ **菜单管理** - 配置系统菜单结构
- ✅ **权限管理** - 细粒度权限控制
- ✅ **系统管理** - 系统配置和监控
- ✅ **个人设置** - 修改个人信息和密码

## 📚 下一步

### 查看API文档
访问Swagger文档: **http://localhost:8083/swagger/index.html**

### 测试API
```bash
# 测试登录API
curl -X POST http://localhost:8083/api/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 阅读完整文档
- 项目总结: `PROJECT_STATUS.md`
- API对接: `web/admin/API_INTEGRATION.md`
- 框架文档: `docs/FRAMEWORK.md`
- 文档索引: `docs/INDEX.md`

## 🛠️ 开发新功能

### 创建新服务
```bash
go run ./cmd/generator -command=service -name=payment
```

### 生成API文档
```bash
make swagger
```

### 运行测试
```bash
make test
```

## ⚠️ 常见问题

### 前端启动失败
```bash
# 清理node_modules重新安装
cd web/admin
rm -rf node_modules package-lock.json
npm install
```

### 后端连接数据库失败
检查 `configs/config.yaml` 中的数据库配置：
```yaml
database:
  host: "localhost"
  port: 3306
  database: "gin_forge"
  username: "root"
  password: "123456"
```

### 登录时提示密码错误
数据库中的默认密码是 `admin123`，不是 `admin` 或 `123456`

## 📞 获取帮助

- 查看完整文档: `docs/INDEX.md`
- 查看示例代码: `docs/demo/`
- 查看API接口: http://localhost:8083/swagger/index.html

---

**让开发更加简单** 🚀

