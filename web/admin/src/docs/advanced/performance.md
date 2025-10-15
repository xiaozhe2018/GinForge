# 性能优化

提升 GinForge 应用性能的完整指南。

## 🎯 性能优化目标

- ⚡ 降低响应时间
- 📈 提高并发能力
- 💰 减少资源消耗
- 🚀 提升用户体验

## 📊 性能分析

### 使用 pprof 分析

```go
import _ "net/http/pprof"

func main() {
    // 启动 pprof 服务
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // 你的应用代码...
}
```

查看性能数据：

```bash
# CPU 分析
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 内存分析
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine 分析
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

## 🗄️ 数据库优化

### 1. 添加索引

```sql
-- 为常查询的字段添加索引
CREATE INDEX idx_username ON users(username);
CREATE INDEX idx_email ON users(email);
CREATE INDEX idx_status_created ON users(status, created_at);

-- 查看索引使用情况
EXPLAIN SELECT * FROM users WHERE username = 'admin';
```

### 2. 避免 N+1 查询

```go
// ❌ N+1 查询问题
var users []model.User
db.Find(&users)  // 1 次查询

for _, user := range users {
    // 每个用户查询一次订单，N 次查询
    db.Where("user_id = ?", user.ID).Find(&user.Orders)
}
// 总共 1 + N 次查询

// ✅ 使用 Preload 优化
var users []model.User
db.Preload("Orders").Find(&users)  // 只需 2 次查询
```

### 3. 只查询需要的字段

```go
// ❌ 查询所有字段
db.Find(&users)

// ✅ 只查询需要的字段
db.Select("id", "username", "email").Find(&users)
```

### 4. 使用批量操作

```go
// ❌ 循环插入
for _, user := range users {
    db.Create(&user)  // N 次数据库操作
}

// ✅ 批量插入
db.CreateInBatches(users, 100)  // 每批 100 条
```

### 5. 连接池优化

```yaml
# configs/config.yaml
database:
  max_idle_conns: 20      # 最大空闲连接（推荐：CPU核心数 * 2）
  max_open_conns: 200     # 最大打开连接
  conn_max_lifetime: 3600 # 连接最大生命周期（秒）
```

## 💾 缓存优化

### 1. 多级缓存

```go
// 应用级缓存（本地） + 分布式缓存（Redis）
func (s *UserService) GetUser(userID string) (*model.User, error) {
    // 1. 查本地缓存（最快）
    if user := s.localCache.Get(userID); user != nil {
        return user, nil
    }
    
    // 2. 查 Redis 缓存
    cacheKey := fmt.Sprintf("user:%s", userID)
    if cached, err := s.redis.Get(ctx, cacheKey); err == nil {
        var user model.User
        json.Unmarshal([]byte(cached), &user)
        
        // 回填本地缓存
        s.localCache.Set(userID, &user, 1*time.Minute)
        
        return &user, nil
    }
    
    // 3. 查数据库
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return nil, err
    }
    
    // 4. 写入缓存
    data, _ := json.Marshal(user)
    s.redis.Set(ctx, cacheKey, string(data), 10*time.Minute)  // Redis: 10分钟
    s.localCache.Set(userID, user, 1*time.Minute)             // 本地: 1分钟
    
    return user, nil
}
```

### 2. 缓存预热

```go
// 系统启动时预加载热点数据
func (s *UserService) WarmupCache() error {
    // 预加载热门用户
    hotUsers, _ := s.userRepo.GetHotUsers(100)
    
    for _, user := range hotUsers {
        cacheKey := fmt.Sprintf("user:%d", user.ID)
        data, _ := json.Marshal(user)
        s.redis.Set(context.Background(), cacheKey, string(data), 1*time.Hour)
    }
    
    log.Info("缓存预热完成", "count", len(hotUsers))
    return nil
}

// 在 main.go 中调用
func main() {
    // ... 初始化代码
    
    // 预热缓存
    userService.WarmupCache()
    
    // 启动服务
    r.Run(":8080")
}
```

### 3. 缓存策略

```go
// 不同数据的缓存时间
const (
    CacheTTLConfig     = 24 * time.Hour    // 配置：24小时
    CacheTTLUser       = 10 * time.Minute  // 用户：10分钟
    CacheTTLList       = 1 * time.Minute   // 列表：1分钟
    CacheTTLRealtimeStats = 10 * time.Second  // 实时统计：10秒
)
```

## ⚡ 代码优化

### 1. 使用 Goroutine

```go
// 并发处理多个任务
func (s *UserService) ProcessUser(userID string) error {
    var wg sync.WaitGroup
    errs := make(chan error, 3)
    
    // 任务 1：发送欢迎邮件
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := s.sendWelcomeEmail(userID); err != nil {
            errs <- err
        }
    }()
    
    // 任务 2：创建用户配置
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := s.createUserConfig(userID); err != nil {
            errs <- err
        }
    }()
    
    // 任务 3：初始化用户数据
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := s.initUserData(userID); err != nil {
            errs <- err
        }
    }()
    
    // 等待所有任务完成
    wg.Wait()
    close(errs)
    
    // 检查是否有错误
    for err := range errs {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### 2. 对象池复用

```go
import "sync"

// 使用 sync.Pool 复用对象
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func processData(data []byte) string {
    // 从池中获取 buffer
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)  // 用完放回池中
    
    buf.Reset()
    buf.Write(data)
    
    return buf.String()
}
```

### 3. 避免不必要的内存分配

```go
// ❌ 每次都创建新切片
func getUsers() []User {
    users := make([]User, 0)  // 分配内存
    // ...
    return users
}

// ✅ 预分配容量
func getUsers() []User {
    users := make([]User, 0, 100)  // 预分配 100 个容量
    // ...
    return users
}
```

## 🌐 HTTP 优化

### 1. 启用 Gzip 压缩

```go
import "github.com/gin-contrib/gzip"

r := gin.Default()
r.Use(gzip.Gzip(gzip.DefaultCompression))
```

### 2. 使用 HTTP/2

```go
// 使用 TLS 启用 HTTP/2
srv := &http.Server{
    Addr:    ":8443",
    Handler: r,
}

srv.ListenAndServeTLS("cert.pem", "key.pem")
```

### 3. 设置合理的超时

```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      r,
    ReadTimeout:  60 * time.Second,
    WriteTimeout: 60 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```

## 🔧 Redis 优化

### 1. Pipeline 批量操作

```go
// ❌ 多次往返
for i := 0; i < 100; i++ {
    redis.Set(ctx, fmt.Sprintf("key:%d", i), value, 0)
}

// ✅ 使用 Pipeline
pipe := redis.Pipeline()
for i := 0; i < 100; i++ {
    pipe.Set(ctx, fmt.Sprintf("key:%d", i), value, 0)
}
pipe.Exec(ctx)  // 一次性执行
```

### 2. 连接池配置

```yaml
redis:
  pool_size: 50           # 连接池大小
  min_idle_conns: 10      # 最小空闲连接
  max_conn_age: 1h        # 连接最大寿命
  pool_timeout: 4s        # 获取连接超时
```

## 📊 监控和调优

### 1. 添加监控指标

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)

// 在中间件中记录指标
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := c.Writer.Status()
        
        httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), fmt.Sprintf("%d", status)).Inc()
        httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
    }
}
```

### 2. 慢查询日志

```go
// GORM 慢查询日志
import "gorm.io/gorm/logger"

newLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
        SlowThreshold: 200 * time.Millisecond,  // 慢查询阈值
        LogLevel:      logger.Warn,
    },
)

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: newLogger,
})
```

## 💡 最佳实践

### 1. 合理使用缓存

```go
// 缓存读多写少的数据
- ✅ 系统配置
- ✅ 用户基本信息
- ✅ 商品信息
- ❌ 实时库存
- ❌ 订单状态
```

### 2. 异步处理

```go
// 不需要立即返回结果的操作，使用异步处理
func (h *UserHandler) Register(c *gin.Context) {
    user, err := h.userService.CreateUser(req)
    if err != nil {
        response.InternalError(c, "注册失败")
        return
    }
    
    // 异步发送欢迎邮件
    go func() {
        h.emailService.SendWelcome(user.Email)
    }()
    
    // 异步初始化用户数据
    go func() {
        h.userService.InitUserData(user.ID)
    }()
    
    // 立即返回
    response.Success(c, user)
}
```

### 3. 数据库优化清单

- ✅ 为常查询字段添加索引
- ✅ 避免 SELECT *，只查询需要的字段
- ✅ 使用 LIMIT 限制查询数量
- ✅ 使用 Preload 避免 N+1 查询
- ✅ 使用批量操作减少数据库往返
- ✅ 合理设置连接池大小
- ✅ 定期分析慢查询日志

### 4. 代码优化清单

- ✅ 使用 Goroutine 并发处理
- ✅ 使用 sync.Pool 复用对象
- ✅ 避免不必要的内存分配
- ✅ 使用 strings.Builder 拼接字符串
- ✅ 避免在循环中做重复计算
- ✅ 使用 defer 要注意性能影响

## 📈 压力测试

### 使用 wrk 压测

```bash
# 安装 wrk
brew install wrk  # macOS
apt install wrk   # Ubuntu

# 压测示例
wrk -t4 -c200 -d30s http://localhost:8083/api/v1/admin/system/health

# 参数说明：
# -t4: 4个线程
# -c200: 200个并发连接
# -d30s: 持续30秒
```

### 分析压测结果

```
Running 30s test @ http://localhost:8083/api/v1/admin/system/health
  4 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    10.23ms    5.67ms  89.12ms   75.23%
    Req/Sec     5.12k     1.23k    8.45k    68.92%
  613456 requests in 30.00s, 145.67MB read
Requests/sec:  20448.53
Transfer/sec:      4.86MB

关键指标：
- QPS: 20448 请求/秒（很好！）
- 平均延迟: 10.23ms（优秀！）
- 最大延迟: 89.12ms（可接受）
```

## 🎯 性能优化案例

### 案例 1：列表查询优化

优化前：

```go
func (s *UserService) GetUsers() ([]model.User, error) {
    var users []model.User
    db.Find(&users)  // 查询所有字段
    
    // 查询每个用户的角色（N+1 查询）
    for i := range users {
        db.Model(&users[i]).Association("Roles").Find(&users[i].Roles)
    }
    
    return users, nil
}

// 性能：200ms
```

优化后：

```go
func (s *UserService) GetUsers() ([]model.User, error) {
    var users []model.User
    
    // 只查询需要的字段 + Preload 关联
    db.Select("id", "username", "email", "status", "created_at").
       Preload("Roles", func(db *gorm.DB) *gorm.DB {
           return db.Select("id", "name")
       }).
       Where("status = ?", 1).
       Limit(20).
       Find(&users)
    
    return users, nil
}

// 性能：15ms（提升 93%！）
```

### 案例 2：配置缓存优化

优化前：

```go
func (s *SystemService) GetConfig(key string) (string, error) {
    // 每次都查数据库
    return s.configRepo.GetByKey(key)
}

// 每次请求都查库，性能差
```

优化后：

```go
func (s *SystemService) GetConfig(key string) (string, error) {
    cacheKey := fmt.Sprintf("config:%s", key)
    
    // 查缓存
    if cached, err := s.redis.Get(ctx, cacheKey); err == nil {
        return cached, nil
    }
    
    // 查数据库
    value, err := s.configRepo.GetByKey(key)
    if err != nil {
        return "", err
    }
    
    // 写缓存（24小时）
    s.redis.Set(ctx, cacheKey, value, 24*time.Hour)
    
    return value, nil
}

// 缓存命中时只需 <1ms
```

## 🚀 部署优化

### 1. 生产模式

```bash
# 设置生产模式
export GIN_MODE=release

# 或在代码中设置
gin.SetMode(gin.ReleaseMode)
```

### 2. 编译优化

```bash
# 减小二进制文件大小
go build -ldflags="-w -s" -o bin/admin-api ./services/admin-api/cmd/server/main.go

# -w: 去掉调试信息
# -s: 去掉符号表
```

### 3. 使用 CDN

```nginx
# 静态资源使用 CDN
location /static/ {
    alias /var/www/static/;
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

## 📚 参考资源

- [数据库优化](../core-concepts/database)
- [缓存系统](../features/cache)
- [部署指南](../deployment/production)

## 🎯 性能优化检查清单

### 应用层

- [ ] 启用生产模式（`GIN_MODE=release`）
- [ ] 使用 Goroutine 并发处理
- [ ] 合理使用缓存
- [ ] 异步处理耗时操作

### 数据库层

- [ ] 添加必要的索引
- [ ] 避免 N+1 查询
- [ ] 只查询需要的字段
- [ ] 使用批量操作
- [ ] 优化连接池配置

### 缓存层

- [ ] 启用 Redis 缓存
- [ ] 设置合理的 TTL
- [ ] 实现缓存预热
- [ ] 使用多级缓存

### 网络层

- [ ] 启用 Gzip 压缩
- [ ] 使用 HTTP/2
- [ ] 配置 CDN
- [ ] 合理设置超时时间

### 监控层

- [ ] 添加性能监控
- [ ] 配置慢查询日志
- [ ] 定期压力测试
- [ ] 建立性能基线

---

**提示**: 性能优化是持续的过程，需要根据实际业务场景不断调整！

