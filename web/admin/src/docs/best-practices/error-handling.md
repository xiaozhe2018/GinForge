# 错误处理

学习如何在 GinForge 中优雅地处理错误。

## 🎯 错误处理原则

1. **永远检查错误** - 不要忽略任何错误
2. **向上传递错误** - 让调用者决定如何处理
3. **添加上下文** - 包装错误时添加有用信息
4. **统一响应格式** - 使用标准化的错误响应
5. **记录错误日志** - 便于问题排查

## 📋 错误码体系

### 错误码定义 (`pkg/errors/codes.go`)

```go
package errors

const (
    // 1xxx - 通用错误
    Success           = 0
    InternalError     = 1000
    InvalidParameter  = 1001
    MissingParameter  = 1002
    InvalidFormat     = 1003
    
    // 2xxx - 认证错误
    Unauthorized      = 2000
    TokenExpired      = 2001
    TokenInvalid      = 2002
    PermissionDenied  = 2003
    
    // 3xxx - 资源错误
    NotFound          = 3000
    AlreadyExists     = 3001
    ResourceConflict  = 3002
    
    // 4xxx - 业务错误
    UserNotFound      = 4000
    UserExists        = 4001
    InvalidPassword   = 4002
    AccountLocked     = 4003
)
```

### 错误消息映射

```go
var errorMessages = map[int]string{
    Success:          "操作成功",
    InternalError:    "服务器内部错误",
    InvalidParameter: "参数错误",
    Unauthorized:     "未授权访问",
    NotFound:         "资源不存在",
    UserExists:       "用户已存在",
}

func GetMessage(code int) string {
    if msg, ok := errorMessages[code]; ok {
        return msg
    }
    return "未知错误"
}
```

## 🔧 错误处理实践

### 1. Service 层错误处理

```go
func (s *UserService) CreateUser(username, email string) (*model.User, error) {
    // 验证参数
    if username == "" {
        return nil, errors.New("用户名不能为空")
    }
    
    // 检查用户是否已存在
    existing, err := s.userRepo.GetByUsername(username)
    if err == nil && existing != nil {
        return nil, errors.New("用户已存在")
    }
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        // 非"记录不存在"的错误，需要返回
        s.LogError("failed to check user existence", err)
        return nil, fmt.Errorf("检查用户失败: %w", err)
    }
    
    // 创建用户
    user := &model.User{
        Username: username,
        Email:    email,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        s.LogError("failed to create user", err)
        return nil, fmt.Errorf("创建用户失败: %w", err)
    }
    
    s.LogInfo("user created successfully", "user_id", user.ID)
    return user, nil
}
```

### 2. Handler 层错误处理

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req model.CreateUserRequest
    
    // 参数绑定错误
    if err := c.ShouldBindJSON(&req); err != nil {
        h.LogError("bind request failed", err)
        h.BadRequest(c, "参数格式错误")
        return
    }
    
    // 调用 Service
    user, err := h.userService.CreateUser(req.Username, req.Email)
    if err != nil {
        h.LogError("create user failed", err)
        
        // 根据错误类型返回不同响应
        if err.Error() == "用户已存在" {
            h.Error(c, 400, "用户已存在")
            return
        }
        
        h.InternalError(c, "创建用户失败")
        return
    }
    
    h.Success(c, user)
}
```

### 3. Repository 层错误处理

```go
func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
    var user model.User
    
    err := r.db.First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("用户不存在: id=%d", id)
        }
        r.LogError("database query failed", err, "user_id", id)
        return nil, fmt.Errorf("查询用户失败: %w", err)
    }
    
    return &user, nil
}
```

## 📊 统一响应格式

### 标准响应结构

```go
type Response struct {
    Code    int         `json:"code"`              // 错误码
    Message string      `json:"message"`           // 错误消息
    Data    interface{} `json:"data,omitempty"`    // 数据
    TraceID string      `json:"trace_id"`          // 追踪 ID
}
```

### 响应示例

成功响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin"
  },
  "trace_id": "xxx-xxx-xxx"
}
```

错误响应：

```json
{
  "code": 4001,
  "message": "用户已存在",
  "trace_id": "xxx-xxx-xxx"
}
```

## 🛡️ 错误包装

### 使用 fmt.Errorf 包装错误

```go
// ✅ 包装错误，保留原始错误
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// 检查特定错误
if errors.Is(err, gorm.ErrRecordNotFound) {
    return errors.New("用户不存在")
}

// 检查错误类型
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    return validationErr.Message
}
```

## 📝 自定义错误类型

### 定义业务错误

```go
// BusinessError 业务错误
type BusinessError struct {
    Code    int
    Message string
    Details interface{}
}

func (e *BusinessError) Error() string {
    return e.Message
}

func NewBusinessError(code int, message string) *BusinessError {
    return &BusinessError{
        Code:    code,
        Message: message,
    }
}

// 使用
if user == nil {
    return NewBusinessError(4000, "用户不存在")
}
```

### 在 Handler 中处理

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.userService.GetUser(userID)
    if err != nil {
        // 检查是否为业务错误
        var bizErr *errors.BusinessError
        if errors.As(err, &bizErr) {
            h.Error(c, bizErr.Code, bizErr.Message)
            return
        }
        
        // 其他错误
        h.InternalError(c, "获取用户失败")
        return
    }
    
    h.Success(c, user)
}
```

## 💡 错误处理最佳实践

### 1. 早返回，避免嵌套

```go
// ✅ 推荐：早返回
func ProcessUser(user *User) error {
    if user == nil {
        return errors.New("user is nil")
    }
    
    if user.Username == "" {
        return errors.New("username is empty")
    }
    
    // 正常逻辑
    return nil
}

// ❌ 避免：深层嵌套
func ProcessUser(user *User) error {
    if user != nil {
        if user.Username != "" {
            // 正常逻辑（嵌套太深）
        } else {
            return errors.New("username is empty")
        }
    } else {
        return errors.New("user is nil")
    }
    return nil
}
```

### 2. 添加上下文信息

```go
// ✅ 添加有用的上下文
if err != nil {
    return fmt.Errorf("failed to create user '%s': %w", username, err)
}

// ❌ 信息不足
if err != nil {
    return err
}
```

### 3. 区分可恢复和不可恢复错误

```go
// 不可恢复：Fatal
if err := db.Connect(); err != nil {
    log.Fatal("cannot connect to database", err)
    // 程序退出
}

// 可恢复：Error
if err := sendEmail(); err != nil {
    log.Error("failed to send email", err)
    // 程序继续运行
}
```

## 📚 实际示例

查看完整错误处理：

- **错误码定义**: `pkg/errors/codes.go`
- **统一响应**: `pkg/response/response.go`
- **Service 错误处理**: `services/admin-api/internal/service/`
- **Handler 错误处理**: `services/admin-api/internal/handler/`

## 🎯 下一步

- [安全建议](./security) - 应用安全最佳实践
- [性能优化](./performance) - 系统性能优化

---

**提示**: 良好的错误处理是高质量代码的标志！

