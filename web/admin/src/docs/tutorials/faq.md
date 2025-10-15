# 常见问题（FAQ）

收集 GinForge 使用过程中的常见问题和解决方案。

## 🚀 启动问题

### Q1: 服务启动失败，提示端口被占用

**错误信息**:
```
Error: listen tcp :8083: bind: address already in use
```

**解决方案**:

```bash
# 方式 1：查找并杀死占用端口的进程
lsof -ti :8083 | xargs kill -9

# 方式 2：修改端口
export APP_PORT=8084
go run ./services/admin-api/cmd/server/main.go

# 方式 3：修改配置文件
# 编辑 configs/config.yaml，修改 services.admin_api.port
```

### Q2: 数据库连接失败

**错误信息**:
```
Error: failed to connect to database
```

**解决方案**:

```bash
# 1. 检查 MySQL 是否运行
docker ps | grep mysql

# 2. 测试连接
mysql -h localhost -u root -p123456

# 3. 检查配置
# 确保 configs/config.yaml 中的数据库配置正确

# 4. 如果使用 SQLite，检查文件权限
ls -la goweb.db
chmod 666 goweb.db
```

### Q3: Redis 连接失败

**错误信息**:
```
Error: redis connection failed
```

**解决方案**:

```bash
# 1. 检查 Redis 是否运行
docker ps | grep redis

# 2. 测试连接
docker exec redis redis-cli ping
# 应该返回: PONG

# 3. 如果不需要 Redis，可以禁用
# 编辑 configs/config.yaml
redis:
  enabled: false
```

## 🔐 认证问题

### Q4: 登录失败，提示用户名或密码错误

**解决方案**:

```bash
# 1. 确认使用正确的账号
# 默认账号：
用户名: admin
密码: admin123

# 2. 检查数据库中的用户
# MySQL:
docker exec mysql mysql -uroot -p123456 gin_forge -e "SELECT username, status FROM admin_users;"

# SQLite:
sqlite3 goweb.db "SELECT username, status FROM admin_users;"

# 3. 如果没有admin用户，重新导入SQL
docker exec -i mysql mysql -uroot -p123456 gin_forge < database/migrations/001_create_admin_tables.sql
```

### Q5: Token 过期太快

**解决方案**:

```yaml
# 修改 configs/config.yaml
jwt:
  expire_hours: 24  # 改为 24 小时或更长
```

或在系统配置中修改会话超时时间。

### Q6: 登录后立即退出

**解决方案**:

```bash
# 1. 打开浏览器控制台（F12），查看错误信息

# 2. 检查 localStorage
# 在控制台执行：
localStorage.getItem('admin_token')
localStorage.getItem('admin_user_info')

# 3. 清除浏览器缓存
localStorage.clear()
# 然后重新登录

# 4. 检查后端日志，看是否有 JWT 验证错误
tail -f logs/admin-api.log
```

## 🌐 API 请求问题

### Q7: API 请求返回 404

**错误信息**:
```
GET http://localhost:8083/api/v1/admin/users 404 (Not Found)
```

**解决方案**:

```bash
# 1. 确认后端服务运行正常
curl http://localhost:8083/api/v1/admin/system/health

# 2. 查看 Swagger 文档确认正确的 API 路径
# 访问: http://localhost:8083/swagger/index.html

# 3. 检查路由配置
# 查看: services/admin-api/internal/router/router.go
```

### Q8: CORS 跨域错误

**错误信息**:
```
Access to XMLHttpRequest has been blocked by CORS policy
```

**解决方案**:

```go
// 检查后端是否添加了 CORS 中间件
// services/admin-api/internal/router/router.go

r.Use(middleware.CORS())

// 或检查 pkg/middleware/cors.go 配置
```

### Q9: 请求返回 401 Unauthorized

**解决方案**:

```typescript
// 1. 检查前端请求拦截器是否正确添加 Token
// web/admin/src/api/index.ts

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 2. 检查 Token 是否过期
// 重新登录获取新 Token

// 3. 检查后端 JWT 中间件配置
```

## 💾 数据库问题

### Q10: AutoMigrate 失败

**错误信息**:
```
Error: failed to auto migrate database
```

**解决方案**:

```go
// 1. 检查数据库连接是否正常

// 2. 检查模型定义是否正确
// 确保所有字段都有正确的 gorm 标签

// 3. 手动执行 SQL 迁移
// 使用 database/migrations/ 下的 SQL 文件

// 4. 查看详细错误信息
if err := db.AutoMigrate(&model.User{}); err != nil {
    log.Error("migrate failed", err)  // 查看具体错误
}
```

### Q11: 查询返回空结果

**问题**: 明明有数据，但查询返回空

**解决方案**:

```go
// 1. 检查是否使用了软删除
// 如果模型有 DeletedAt 字段，查询会自动过滤已删除记录

// 查询包括已删除的记录
db.Unscoped().Find(&users)

// 2. 检查查询条件是否正确
// 3. 查看 SQL 日志
// 在配置中启用 SQL 日志：
database:
  log_level: "info"  # 显示所有 SQL
```

## 📦 依赖问题

### Q12: go mod download 失败

**错误信息**:
```
go: module github.com/xxx: Get "https://proxy.golang.org/...": dial tcp: i/o timeout
```

**解决方案**:

```bash
# 设置 Go 代理
export GOPROXY=https://goproxy.cn,direct

# 或使用其他代理
export GOPROXY=https://goproxy.io,direct

# 然后重新下载
go mod download
go mod tidy
```

### Q13: npm install 失败

**解决方案**:

```bash
# 方式 1：清理缓存重新安装
cd web/admin
rm -rf node_modules package-lock.json
npm cache clean --force
npm install

# 方式 2：使用国内镜像
npm config set registry https://registry.npmmirror.com
npm install

# 方式 3：使用 cnpm
npm install -g cnpm --registry=https://registry.npmmirror.com
cnpm install
```

## 🎨 前端问题

### Q14: 前端页面空白

**解决方案**:

```bash
# 1. 打开浏览器控制台（F12），查看错误信息

# 2. 检查前端服务是否正常运行
# 应该看到 Vite 启动信息

# 3. 检查 API 代理配置
# web/admin/vite.config.ts

server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8083',
      changeOrigin: true
    }
  }
}

# 4. 清除浏览器缓存
Ctrl+Shift+Delete（Chrome）

# 5. 重新启动前端
npm run dev
```

### Q15: Element Plus 组件不显示

**解决方案**:

```bash
# 1. 确认 Element Plus 已安装
npm list element-plus

# 2. 检查自动导入配置
# vite.config.ts 应该有 AutoImport 和 Components 插件

# 3. 重新安装依赖
npm install element-plus
```

## 🔧 开发问题

### Q16: 如何调试后端代码？

**方案 1：使用 Delve**

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试服务
dlv debug ./services/admin-api/cmd/server/main.go

# 设置断点
(dlv) break main.main
(dlv) break handler.GetUser
(dlv) continue
```

**方案 2：使用 VS Code**

创建 `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug admin-api",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/services/admin-api/cmd/server"
    }
  ]
}
```

### Q17: 如何生成 Swagger 文档？

```bash
# 1. 安装 swag
go install github.com/swaggo/swag/cmd/swag@latest

# 2. 在代码中添加注释（参考已有代码）

# 3. 生成文档
cd services/admin-api
swag init -g cmd/server/main.go -o docs

# 4. 访问文档
# http://localhost:8083/swagger/index.html
```

## 🐳 Docker 问题

### Q18: Docker 容器无法启动

**解决方案**:

```bash
# 1. 查看容器日志
docker logs ginforge-admin

# 2. 检查端口映射
docker ps -a

# 3. 检查网络配置
docker network ls
docker network inspect ginforge-network

# 4. 重新构建镜像
docker-compose build --no-cache
docker-compose up -d
```

### Q19: 容器内无法连接 MySQL/Redis

**解决方案**:

```yaml
# 确保所有服务在同一网络
# docker-compose.yml

networks:
  ginforge-network:
    driver: bridge

services:
  mysql:
    networks:
      - ginforge-network
  
  admin-api:
    networks:
      - ginforge-network
    environment:
      DB_HOST: mysql  # 使用服务名而不是 localhost
```

## 📊 性能问题

### Q20: API 响应慢

**解决方案**:

```bash
# 1. 启用 Redis 缓存
redis:
  enabled: true

# 2. 优化数据库查询
# - 添加索引
# - 避免 N+1 查询
# - 使用 Preload 预加载关联数据

# 3. 增加数据库连接池
database:
  max_open_conns: 100
  max_idle_conns: 20

# 4. 使用缓存中间件
r.GET("/config", middleware.Cache(10*time.Minute, redis), getConfig)

# 5. 查看慢查询
# MySQL:
docker exec mysql mysql -uroot -p123456 -e "SHOW VARIABLES LIKE 'slow_query%';"
```

## 💡 开发技巧

### Q21: 如何快速创建新功能？

**步骤**:

1. 定义数据模型（`internal/model/`）
2. 创建 Repository（`internal/repository/`）
3. 创建 Service（`internal/service/`）
4. 创建 Handler（`internal/handler/`）
5. 注册路由（`internal/router/router.go`）

参考：[创建完整的业务模块](../tutorials/create-module)

### Q22: 如何添加自定义中间件？

```go
// pkg/middleware/custom.go
package middleware

import "github.com/gin-gonic/gin"

func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 请求前处理
        startTime := time.Now()
        
        c.Next()  // 继续处理
        
        // 请求后处理
        duration := time.Since(startTime)
        log.Printf("Request took %v", duration)
    }
}

// 在路由中使用
r.Use(middleware.CustomMiddleware())
```

### Q23: 如何实现文件上传？

参考：[文件上传](../features/file-upload)

### Q24: 如何实现 WebSocket？

参考：[WebSocket 实时通信](../features/websocket)

## 📚 学习建议

### Q25: 新手应该如何学习？

**推荐学习路径**:

1. **第 1 天**：
   - [框架介绍](../getting-started/introduction)
   - [快速开始](../getting-started/quick-start)
   - 运行起来，熟悉界面

2. **第 2-3 天**:
   - [项目结构](../getting-started/project-structure)
   - [配置系统](../core-concepts/configuration)
   - [路由管理](../core-concepts/routing)

3. **第 4-5 天**:
   - [中间件](../core-concepts/middleware)
   - [数据库操作](../core-concepts/database)
   - [基础类](../api-reference/base-classes)

4. **第 1-2 周**:
   - [创建完整模块](../tutorials/create-module)
   - [认证授权](../features/authentication)
   - [文件上传](../features/file-upload)

5. **第 3-4 周**:
   - [消息队列](../advanced/message-queue)
   - [缓存系统](../features/cache)
   - [WebSocket](../features/websocket)

### Q26: 文档看不懂怎么办？

**建议**:

1. ✅ 先跑起来，边用边学
2. ✅ 查看 Swagger API 文档
3. ✅ 阅读 `docs/demo/` 下的示例代码
4. ✅ 参考实际项目代码
5. ✅ 在 GitHub 提 Issue 提问

## 🔧 部署问题

### Q27: 如何部署到生产环境？

参考：[生产部署](../deployment/production)

### Q28: 如何使用 Docker 部署？

参考：[Docker 部署](../deployment/docker)

### Q29: 如何配置 HTTPS？

**使用 Nginx 反向代理**:

```nginx
server {
    listen 443 ssl http2;
    server_name api.ginforge.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    location / {
        proxy_pass http://127.0.0.1:8083;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 🎯 其他问题

### Q30: 如何贡献代码？

1. Fork 项目
2. 创建功能分支
3. 编写代码和测试
4. 提交 Pull Request
5. 等待 Review

### Q31: 如何报告 Bug？

访问 GitHub Issues，提供以下信息：

- 系统环境（OS、Go 版本等）
- 错误信息和日志
- 复现步骤
- 期望行为

### Q32: 如何获取技术支持？

- 📖 查看完整文档
- 💬 GitHub Discussions
- 🐛 GitHub Issues
- 📧 邮件联系

---

## 📝 没找到答案？

如果以上问题都不能解决你的问题：

1. 📖 查看 [完整文档](../getting-started/introduction)
2. 🔍 搜索 GitHub Issues
3. 💬 在 Discussions 提问
4. 🐛 提交新的 Issue

---

**提示**: 大部分问题都可以通过查看日志和错误信息来定位！

