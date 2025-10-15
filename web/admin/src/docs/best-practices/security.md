# 安全建议

GinForge 框架的安全最佳实践和建议。

## 🔐 认证安全

### 1. 密码安全

#### 使用强密码策略

```go
// 密码复杂度要求
type PasswordPolicy struct {
    MinLength        int      // 最小长度：8-16
    RequireLowercase bool     // 需要小写字母
    RequireUppercase bool     // 需要大写字母
    RequireNumbers   bool     // 需要数字
    RequireSymbols   bool     // 需要特殊字符
}

// 验证密码复杂度
func ValidatePassword(password string, policy PasswordPolicy) error {
    if len(password) < policy.MinLength {
        return fmt.Errorf("密码长度至少为 %d 位", policy.MinLength)
    }
    
    if policy.RequireLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
        return errors.New("密码必须包含小写字母")
    }
    
    // ... 其他检查
    
    return nil
}
```

#### 使用 bcrypt 加密

```go
import "golang.org/x/crypto/bcrypt"

// ✅ 正确：使用 bcrypt
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// ❌ 错误：使用 MD5/SHA1（不安全）
hash := md5.Sum([]byte(password))  // 不要这样做！
```

### 2. JWT 安全

#### 使用强密钥

```go
// ✅ 生成强随机密钥
jwtSecret := generateRandomString(64)

// ❌ 避免弱密钥
jwtSecret := "123456"  // 太弱！
jwtSecret := "my-secret"  // 太简单！
```

#### 设置合理的过期时间

```go
// ✅ 根据场景设置
claims["exp"] = time.Now().Add(24 * time.Hour).Unix()  // 普通用户：24小时
claims["exp"] = time.Now().Add(2 * time.Hour).Unix()   // 管理员：2小时
claims["exp"] = time.Now().Add(30 * 24 * time.Hour).Unix()  // 记住我：30天

// ❌ 避免过长
claims["exp"] = time.Now().Add(365 * 24 * time.Hour).Unix()  // 太长！
```

#### Token 黑名单

```go
// 退出登录时加入黑名单
func Logout(token string, redisClient *redis.Client) error {
    key := fmt.Sprintf("token:blacklist:%s", token)
    return redisClient.Set(ctx, key, "1", 24*time.Hour)
}

// 验证时检查黑名单
func ValidateToken(token string, redisClient *redis.Client) error {
    key := fmt.Sprintf("token:blacklist:%s", token)
    exists, _ := redisClient.Exists(ctx, key)
    if exists {
        return errors.New("token已失效")
    }
    // 继续验证...
}
```

### 3. 账户安全

#### 登录失败锁定

```go
// 记录登录失败次数
func recordLoginFailure(username string, redisClient *redis.Client) {
    key := fmt.Sprintf("login:failures:%s", username)
    
    // 增加失败次数
    failures := redisClient.Incr(ctx, key).Val()
    
    // 首次失败设置过期时间
    if failures == 1 {
        redisClient.Expire(ctx, key, 15*time.Minute)
    }
    
    // 超过 5 次，锁定账户 30 分钟
    if failures >= 5 {
        lockKey := fmt.Sprintf("login:locked:%s", username)
        redisClient.Set(ctx, lockKey, "1", 30*time.Minute)
    }
}
```

## 🛡️ 输入验证

### 1. 参数验证

```go
// 使用 binding 标签
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8,max=128"`
    Phone    string `json:"phone" binding:"omitempty,len=11"`
}

// 自定义验证规则
func ValidateUsername(username string) error {
    // 只允许字母、数字、下划线
    if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
        return errors.New("用户名只能包含字母、数字和下划线")
    }
    
    // 不允许的用户名
    forbidden := []string{"admin", "root", "system"}
    if contains(forbidden, username) {
        return errors.New("该用户名不可使用")
    }
    
    return nil
}
```

### 2. SQL 注入防护

```go
// ✅ 使用参数化查询
db.Where("username = ?", username).First(&user)

// ✅ 使用命名参数
db.Where("username = @username AND status = @status", 
    sql.Named("username", username),
    sql.Named("status", status),
).First(&user)

// ❌ 避免字符串拼接
query := "SELECT * FROM users WHERE username = '" + username + "'"  // 危险！
db.Raw(query).Scan(&user)
```

### 3. XSS 防护

```go
import "html"

// 对用户输入进行转义
func SanitizeInput(input string) string {
    return html.EscapeString(input)
}

// 在保存前清理
user.Bio = SanitizeInput(req.Bio)
```

## 🌐 网络安全

### 1. HTTPS

```nginx
# 强制使用 HTTPS
server {
    listen 80;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
}
```

### 2. CORS 配置

```go
// ✅ 限制允许的来源
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // 只允许特定域名
        allowedOrigins := []string{
            "https://admin.ginforge.com",
            "https://www.ginforge.com",
        }
        
        if contains(allowedOrigins, origin) {
            c.Header("Access-Control-Allow-Origin", origin)
        }
        
        c.Next()
    }
}

// ❌ 避免允许所有来源（生产环境）
c.Header("Access-Control-Allow-Origin", "*")  // 不安全！
```

### 3. 限流保护

```go
// API 限流
r.Use(middleware.RateLimit(100))  // 每分钟 100 次

// 登录接口严格限流
r.POST("/login", middleware.RateLimit(10), loginHandler)

// IP 级别限流
func IPRateLimit() gin.HandlerFunc {
    limiters := make(map[string]*rate.Limiter)
    
    return func(c *gin.Context) {
        ip := c.ClientIP()
        
        limiter, exists := limiters[ip]
        if !exists {
            limiter = rate.NewLimiter(rate.Every(time.Minute), 100)
            limiters[ip] = limiter
        }
        
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "请求过于频繁"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## 🔒 数据安全

### 1. 敏感数据加密

```go
// 加密敏感字段
type User struct {
    ID       uint64
    Phone    string  // 存储加密后的手机号
    IDCard   string  // 存储加密后的身份证号
}

// 加密
encrypted, err := encrypt(plaintext, key)

// 解密
plaintext, err := decrypt(encrypted, key)
```

### 2. 日志脱敏

```go
// ✅ 不要记录敏感信息
log.Info("user login", "username", username)  // OK

// ❌ 避免记录密码等敏感信息
log.Info("user login", "username", username, "password", password)  // 危险！

// ✅ 脱敏处理
func maskPhone(phone string) string {
    if len(phone) < 11 {
        return phone
    }
    return phone[:3] + "****" + phone[7:]
}

log.Info("user info", "phone", maskPhone(user.Phone))
// 输出：138****5678
```

## 🎯 权限控制

### 1. 最小权限原则

```go
// ✅ 只授予必要的权限
role := &model.Role{
    Name: "Editor",
    Permissions: []string{
        "article:read",
        "article:create",
        "article:update",
        // 不包含 article:delete
    },
}

// ❌ 避免授予过多权限
role.Permissions = []string{"*"}  // 危险！
```

### 2. 垂直权限检查

```go
// 检查用户是否有权限操作该资源
func (s *ArticleService) UpdateArticle(userID uint64, articleID uint64, data map[string]interface{}) error {
    // 获取文章
    article, err := s.articleRepo.GetByID(articleID)
    if err != nil {
        return err
    }
    
    // 检查是否是作者或管理员
    if article.AuthorID != userID && !s.isAdmin(userID) {
        return errors.New("无权修改该文章")
    }
    
    // 更新文章
    return s.articleRepo.Update(article)
}
```

## ⚠️ 常见安全问题

### 1. SQL 注入

```go
// ❌ 危险
username := c.Query("username")
db.Raw("SELECT * FROM users WHERE username = '" + username + "'").Scan(&user)

// ✅ 安全
db.Where("username = ?", username).First(&user)
```

### 2. 目录遍历

```go
// ❌ 危险
filepath := c.Query("file")
c.File(filepath)  // 用户可能传入 "../../etc/passwd"

// ✅ 安全
func SafeFilePath(userInput string) (string, error) {
    // 清理路径
    cleanPath := filepath.Clean(userInput)
    
    // 检查是否包含 ..
    if strings.Contains(cleanPath, "..") {
        return "", errors.New("invalid file path")
    }
    
    // 限制在特定目录
    baseDir := "/opt/uploads"
    fullPath := filepath.Join(baseDir, cleanPath)
    
    // 确保在允许的目录内
    if !strings.HasPrefix(fullPath, baseDir) {
        return "", errors.New("access denied")
    }
    
    return fullPath, nil
}
```

### 3. CSRF 防护

```go
// 生成 CSRF Token
func GenerateCSRFToken() string {
    return utils.RandomString(32)
}

// 验证 CSRF Token
func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method != "GET" {
            token := c.GetHeader("X-CSRF-Token")
            sessionToken := c.GetString("csrf_token")
            
            if token == "" || token != sessionToken {
                c.JSON(403, gin.H{"error": "CSRF token invalid"})
                c.Abort()
                return
            }
        }
        c.Next()
    }
}
```

## 📋 安全检查清单

### 代码层面

- [ ] 所有密码使用 bcrypt 加密
- [ ] JWT Secret 使用强随机字符串
- [ ] 使用参数化查询防止 SQL 注入
- [ ] 验证所有用户输入
- [ ] 脱敏处理敏感信息日志
- [ ] 设置 Session 超时
- [ ] 实现登录失败锁定

### 配置层面

- [ ] 生产环境关闭 Debug 模式
- [ ] 使用环境变量存储敏感配置
- [ ] 限制 CORS 允许的域名
- [ ] 配置请求限流
- [ ] 启用 HTTPS
- [ ] 设置安全的 HTTP 头

### 部署层面

- [ ] 使用防火墙限制端口访问
- [ ] 数据库不对外暴露
- [ ] Redis 设置密码
- [ ] 定期更新依赖包
- [ ] 定期备份数据
- [ ] 配置入侵检测

## 🔍 安全审计

### 操作日志

```go
// 记录所有敏感操作
func (s *UserService) DeleteUser(operatorID, userID uint64) error {
    // 记录审计日志
    s.auditLog.Record(&AuditLog{
        OperatorID: operatorID,
        Action:     "delete_user",
        Target:     fmt.Sprintf("user:%d", userID),
        IP:         s.clientIP,
        Timestamp:  time.Now(),
    })
    
    // 执行删除
    return s.userRepo.Delete(userID)
}
```

### 查询审计日志

```sql
-- 查询某用户的所有操作
SELECT * FROM audit_logs WHERE operator_id = 123 ORDER BY created_at DESC;

-- 查询敏感操作
SELECT * FROM audit_logs WHERE action IN ('delete_user', 'update_permission') ORDER BY created_at DESC;
```

## 📚 相关文档

- [认证授权](../features/authentication) - 认证系统详解
- [生产部署](../deployment/production) - 生产环境配置
- [错误处理](./error-handling) - 错误处理最佳实践

## 🎯 安全资源

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go 安全编码规范](https://github.com/Checkmarx/Go-SCP)
- [JWT 最佳实践](https://tools.ietf.org/html/rfc8725)

---

**提示**: 安全无小事，务必重视每一个安全细节！

