# 代码规范

遵循统一的代码规范，提高代码质量和可维护性。

## 📐 Go 代码规范

### 1. 命名规范

#### 包名

```go
// ✅ 使用小写，简短，有意义
package user
package handler
package repository

// ❌ 避免
package User         // 不要使用大写
package userHandler  // 不要使用驼峰
```

#### 文件名

```go
// ✅ 使用下划线分隔
user_handler.go
admin_service.go
user_repository.go

// ❌ 避免
userHandler.go   // 不要使用驼峰
UserHandler.go   // 不要使用大写开头
```

#### 变量和函数

```go
// ✅ 变量和私有函数：驼峰命名
var userName string
var userCount int
func getUserByID() {}

// ✅ 公开函数和方法：大写开头
func NewUserService() *UserService {}
func (s *UserService) GetUser() {}

// ✅ 常量：大写或驼峰
const MaxRetryCount = 3
const UserStatusActive = 1

// ❌ 避免
var user_name string   // 不要使用下划线
func Get_User() {}     // 不要使用下划线
```

### 2. 代码组织

#### 导入顺序

```go
import (
    // 1. 标准库
    "context"
    "fmt"
    "time"
    
    // 2. 第三方库
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    // 3. 项目内部包
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/services/admin-api/internal/model"
)
```

#### 结构体字段顺序

```go
type User struct {
    // 1. ID 字段
    ID uint64 `json:"id" gorm:"primaryKey"`
    
    // 2. 基本字段
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"-"`
    
    // 3. 状态字段
    Status int8 `json:"status"`
    
    // 4. 时间字段
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at"`
    
    // 5. 关联字段
    Roles []Role `json:"roles" gorm:"many2many:user_roles;"`
}
```

### 3. 注释规范

#### 包注释

```go
// Package handler 提供 HTTP 请求处理器
package handler
```

#### 函数注释

```go
// GetUser 根据用户 ID 获取用户信息
// 
// 参数:
//   - userID: 用户 ID
// 
// 返回:
//   - *model.User: 用户信息
//   - error: 错误信息，nil 表示成功
func (s *UserService) GetUser(userID string) (*model.User, error) {
    // 实现...
}
```

#### Swagger 注释

```go
// GetUser godoc
// @Summary 获取用户信息
// @Description 根据用户 ID 获取用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=model.User}
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    // 实现...
}
```

## 🏗️ 项目结构规范

### Service 层结构

```go
services/admin-api/
└── internal/
    ├── model/              # 数据模型
    ├── repository/         # 数据访问层
    ├── service/            # 业务逻辑层
    ├── handler/            # HTTP 处理层
    └── router/             # 路由配置
```

### 文件命名规范

```
handler/
├── admin_auth_handler.go      # 认证处理器
├── admin_user_handler.go      # 用户处理器
└── admin_system_handler.go    # 系统处理器

service/
├── admin_service.go           # 核心服务
├── admin_system_service.go    # 系统服务
└── notification_service.go    # 通知服务
```

## ✅ 代码质量

### 1. 错误处理

```go
// ✅ 始终检查错误
user, err := s.userRepo.GetByID(userID)
if err != nil {
    s.LogError("failed to get user", err, "user_id", userID)
    return nil, err
}

// ❌ 不要忽略错误
user, _ := s.userRepo.GetByID(userID)  // 危险！
```

### 2. 使用 defer

```go
// ✅ 使用 defer 确保资源释放
file, err := os.Open("config.yaml")
if err != nil {
    return err
}
defer file.Close()  // 确保关闭

// 处理文件...
```

### 3. 避免 panic

```go
// ✅ 返回错误
func GetUser(id string) (*User, error) {
    if id == "" {
        return nil, errors.New("id cannot be empty")
    }
    // ...
}

// ❌ 避免使用 panic
func GetUser(id string) *User {
    if id == "" {
        panic("id cannot be empty")  // 不推荐
    }
    // ...
}
```

## 📝 TypeScript 规范（前端）

### 命名规范

```typescript
// 文件名：PascalCase.vue 或 camelCase.ts
UserList.vue
userApi.ts

// 变量和函数：camelCase
const userName = 'john'
function getUser() {}

// 类型和接口：PascalCase
interface User {}
type UserStatus = 'active' | 'disabled'

// 常量：UPPER_SNAKE_CASE
const MAX_UPLOAD_SIZE = 1024 * 1024
```

### 类型定义

```typescript
// ✅ 定义清晰的接口
interface User {
  id: number
  username: string
  email: string
  status: number
  created_at: string
}

// ✅ 使用类型安全
const user: User = {
  id: 1,
  username: 'john',
  // ... 必须包含所有字段
}

// ❌ 避免使用 any
function getData(): any {}  // 不推荐
```

## 🎯 Git 提交规范

### Commit Message 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

- `feat`: 新功能
- `fix`: 修复 Bug
- `docs`: 文档更新
- `style`: 代码格式（不影响代码运行）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

### 示例

```bash
# 添加新功能
git commit -m "feat(user): add user export功能"

# 修复 Bug
git commit -m "fix(auth): fix login token expiration issue"

# 更新文档
git commit -m "docs: update API documentation"

# 重构代码
git commit -m "refactor(service): simplify user service logic"
```

## 🧪 测试规范

### 单元测试

```go
// user_service_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
    // Arrange（准备）
    service := NewUserService(mockRepo, mockLogger)
    username := "testuser"
    
    // Act（执行）
    user, err := service.CreateUser(username, "test@example.com")
    
    // Assert（断言）
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, username, user.Username)
}

func TestCreateUser_AlreadyExists(t *testing.T) {
    service := NewUserService(mockRepo, mockLogger)
    
    user, err := service.CreateUser("existing", "test@example.com")
    
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "已存在")
}
```

## 📚 推荐工具

### Linter

```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 运行检查
golangci-lint run

# 自动修复
golangci-lint run --fix
```

### 格式化

```bash
# 格式化代码
go fmt ./...

# 或使用 goimports（推荐）
goimports -w .
```

## 🎯 下一步

- [错误处理](./error-handling) - 优雅的错误处理
- [安全建议](./security) - 应用安全最佳实践

---

**提示**: 良好的代码规范是团队协作的基础，请务必遵循！

