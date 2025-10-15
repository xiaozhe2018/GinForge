# 数据库操作

GinForge 使用 GORM 作为 ORM 框架，支持 SQLite、MySQL、PostgreSQL 等多种数据库。

## 🎯 数据库支持

| 数据库 | 驱动 | 适用场景 |
|--------|------|----------|
| SQLite | `gorm.io/driver/sqlite` | 开发环境、小型应用 |
| MySQL | `gorm.io/driver/mysql` | 生产环境（推荐） |
| PostgreSQL | `gorm.io/driver/postgres` | 生产环境 |

## 🔧 配置数据库

### SQLite 配置（默认）

```yaml
# configs/config.yaml
database:
  type: "sqlite"
  database: "goweb.db"
```

### MySQL 配置

```yaml
database:
  type: "mysql"
  host: "localhost"
  port: 3306
  database: "gin_forge"
  username: "root"
  password: "123456"
  charset: "utf8mb4"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
```

## 📦 初始化数据库

### 在 main.go 中初始化

```go
package main

import (
    "goweb/pkg/config"
    "goweb/pkg/db"
    "goweb/pkg/logger"
    "goweb/pkg/model"
)

func main() {
    cfg := config.New()
    log := logger.New("admin-api", cfg.GetString("log.level"))
    
    // 初始化数据库
    database, err := db.New(cfg)
    if err != nil {
        log.Fatal("failed to initialize database", err)
    }
    
    // 自动迁移表结构
    if err := database.AutoMigrate(
        &model.AdminUser{},
        &model.AdminRole{},
        &model.AdminPermission{},
        &model.AdminMenu{},
    ); err != nil {
        log.Warn("failed to auto migrate", "error", err)
    }
    
    log.Info("database initialized successfully")
}
```

## 📊 定义模型

### 基础模型

```go
package model

import (
    "time"
)

// User 用户模型
type User struct {
    ID        uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
    Username  string     `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
    Email     string     `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
    Password  string     `json:"-" gorm:"type:varchar(255);not null"`  // - 表示不序列化
    Name      string     `json:"name" gorm:"type:varchar(50)"`
    Phone     string     `json:"phone" gorm:"type:varchar(20);index"`
    Status    int8       `json:"status" gorm:"type:tinyint(1);default:1;index"`
    CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt *time.Time `json:"deleted_at" gorm:"index"`  // 软删除
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}
```

### 关联关系

```go
// User 用户模型
type User struct {
    ID    uint64  `json:"id" gorm:"primaryKey"`
    Name  string  `json:"name"`
    
    // 一对多：一个用户有多个订单
    Orders []Order `json:"orders" gorm:"foreignKey:UserID"`
    
    // 多对多：用户与角色
    Roles []Role `json:"roles" gorm:"many2many:user_roles;"`
}

// Order 订单模型
type Order struct {
    ID     uint64 `json:"id" gorm:"primaryKey"`
    UserID uint64 `json:"user_id"`
    Amount float64 `json:"amount"`
    
    // 属于某个用户
    User User `json:"user" gorm:"foreignKey:UserID"`
}

// Role 角色模型
type Role struct {
    ID   uint64 `json:"id" gorm:"primaryKey"`
    Name string `json:"name"`
    
    // 多对多：角色与用户
    Users []User `json:"users" gorm:"many2many:user_roles;"`
}
```

## 🔍 CRUD 操作

### 1. 创建 (Create)

```go
// 创建单条记录
user := &model.User{
    Username: "john",
    Email:    "john@example.com",
    Password: "hashed_password",
    Name:     "John Doe",
}

result := db.Create(user)
if result.Error != nil {
    log.Error("create user failed", result.Error)
}

log.Info("user created", "id", user.ID)  // user.ID 会自动赋值

// 批量创建
users := []model.User{
    {Username: "user1", Email: "user1@example.com"},
    {Username: "user2", Email: "user2@example.com"},
}
db.Create(&users)
```

### 2. 查询 (Read)

```go
// 根据主键查询
var user model.User
db.First(&user, 1)  // SELECT * FROM users WHERE id = 1

// 根据条件查询
db.Where("username = ?", "admin").First(&user)
db.Where("status = ? AND created_at > ?", 1, yesterday).Find(&users)

// 查询所有
var users []model.User
db.Find(&users)

// 条件查询
db.Where("status = ?", 1).
   Where("created_at > ?", time.Now().AddDate(0, -1, 0)).
   Order("created_at DESC").
   Limit(10).
   Find(&users)
```

### 3. 更新 (Update)

```go
// 更新单个字段
db.Model(&user).Update("name", "New Name")

// 更新多个字段
db.Model(&user).Updates(model.User{
    Name:   "New Name",
    Status: 1,
})

// 使用 map 更新
db.Model(&user).Updates(map[string]interface{}{
    "name":   "New Name",
    "status": 1,
})

// 批量更新
db.Model(&model.User{}).Where("status = ?", 0).Update("status", 1)
```

### 4. 删除 (Delete)

```go
// 删除单条记录
db.Delete(&user, 1)  // DELETE FROM users WHERE id = 1

// 软删除（如果模型有 DeletedAt 字段）
db.Delete(&user)  // UPDATE users SET deleted_at = NOW() WHERE id = 1

// 永久删除
db.Unscoped().Delete(&user)

// 批量删除
db.Where("status = ?", 0).Delete(&model.User{})
```

## 🔄 事务处理

### 基础事务

```go
// 开始事务
tx := db.Begin()

// 执行操作
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()  // 回滚
    return err
}

if err := tx.Create(&order).Error; err != nil {
    tx.Rollback()
    return err
}

// 提交事务
tx.Commit()
```

### 使用闭包事务（推荐）

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // 创建用户
    if err := tx.Create(&user).Error; err != nil {
        return err  // 自动回滚
    }
    
    // 创建订单
    if err := tx.Create(&order).Error; err != nil {
        return err  // 自动回滚
    }
    
    // 返回 nil 自动提交
    return nil
})

if err != nil {
    log.Error("transaction failed", err)
}
```

## 📈 高级查询

### 1. 关联查询

```go
// 预加载关联数据
var user model.User
db.Preload("Roles").First(&user, 1)
// 同时查询用户和角色

// 预加载多个关联
db.Preload("Roles").
   Preload("Orders").
   First(&user, 1)

// 条件预加载
db.Preload("Orders", "status = ?", 1).First(&user)
```

### 2. 联表查询

```go
var users []model.User

// 使用 Joins
db.Joins("LEFT JOIN orders ON orders.user_id = users.id").
   Where("orders.status = ?", 1).
   Find(&users)

// 手动 JOIN
db.Table("users").
   Select("users.*, orders.amount").
   Joins("LEFT JOIN orders ON orders.user_id = users.id").
   Scan(&results)
```

### 3. 聚合查询

```go
// 计数
var count int64
db.Model(&model.User{}).Where("status = ?", 1).Count(&count)

// 求和
var total float64
db.Model(&model.Order{}).Select("SUM(amount)").Row().Scan(&total)

// 分组统计
type Result struct {
    Date  string
    Count int
}
var results []Result
db.Model(&model.Order{}).
   Select("DATE(created_at) as date, COUNT(*) as count").
   Group("DATE(created_at)").
   Scan(&results)
```

### 4. 原始 SQL

```go
// 执行原始 SQL
db.Exec("UPDATE users SET status = ? WHERE id = ?", 1, 123)

// 查询原始 SQL
var users []model.User
db.Raw("SELECT * FROM users WHERE status = ? AND created_at > ?", 1, yesterday).Scan(&users)

// 使用命名参数
db.Raw("SELECT * FROM users WHERE username = @username", 
    sql.Named("username", "admin")).Scan(&user)
```

## 🏗️ Repository 模式

GinForge 推荐使用 Repository 模式封装数据库操作：

```go
// repository/user_repository.go
package repository

import (
    "gorm.io/gorm"
    "goweb/pkg/model"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

// GetByID 根据 ID 获取
func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
    var user model.User
    err := r.db.First(&user, id).Error
    return &user, err
}

// GetByUsername 根据用户名获取
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
    var user model.User
    err := r.db.Where("username = ?", username).First(&user).Error
    return &user, err
}

// List 获取用户列表（分页）
func (r *UserRepository) List(page, pageSize int) ([]model.User, int64, error) {
    var users []model.User
    var total int64
    
    // 计数
    r.db.Model(&model.User{}).Count(&total)
    
    // 分页查询
    offset := (page - 1) * pageSize
    err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error
    
    return users, total, err
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
    return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *UserRepository) Delete(id uint64) error {
    return r.db.Delete(&model.User{}, id).Error
}
```

### 在 Service 中使用

```go
// service/user_service.go
package service

type UserService struct {
    userRepo *repository.UserRepository
}

func (s *UserService) CreateUser(username, email string) (*model.User, error) {
    // 检查用户是否已存在
    _, err := s.userRepo.GetByUsername(username)
    if err == nil {
        return nil, errors.New("用户已存在")
    }
    
    // 创建用户
    user := &model.User{
        Username: username,
        Email:    email,
        Status:   1,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

## 💾 数据迁移

### 自动迁移

```go
// 自动创建/更新表结构
db.AutoMigrate(
    &model.User{},
    &model.Role{},
    &model.Permission{},
)
```

### SQL 迁移脚本

```sql
-- database/migrations/001_create_users_table.sql
CREATE TABLE IF NOT EXISTS `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(50) NOT NULL,
    `email` varchar(100) NOT NULL,
    `password` varchar(255) NOT NULL,
    `name` varchar(50) DEFAULT NULL,
    `phone` varchar(20) DEFAULT NULL,
    `status` tinyint(1) NOT NULL DEFAULT 1,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    UNIQUE KEY `idx_email` (`email`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

执行迁移：

```bash
# 导入 SQL 文件
mysql -h localhost -u root -p123456 gin_forge < database/migrations/001_create_users_table.sql
```

## 🔍 查询技巧

### 1. 链式查询

```go
// 构建复杂查询
query := db.Model(&model.User{})

if keyword != "" {
    query = query.Where("username LIKE ? OR email LIKE ?", 
        "%"+keyword+"%", "%"+keyword+"%")
}

if status != 0 {
    query = query.Where("status = ?", status)
}

if startDate != nil {
    query = query.Where("created_at >= ?", startDate)
}

// 执行查询
var users []model.User
query.Order("created_at DESC").Limit(10).Find(&users)
```

### 2. Scopes（查询作用域）

```go
// 定义 Scope
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", 1)
}

func RecentUsers(days int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        since := time.Now().AddDate(0, 0, -days)
        return db.Where("created_at > ?", since)
    }
}

// 使用 Scope
var users []model.User
db.Scopes(ActiveUsers, RecentUsers(7)).Find(&users)
// 查询最近 7 天的活跃用户
```

### 3. 子查询

```go
// 使用子查询
subQuery := db.Model(&model.Order{}).
    Select("user_id").
    Where("total_amount > ?", 1000)

var users []model.User
db.Where("id IN (?)", subQuery).Find(&users)
// 查询订单金额大于 1000 的用户
```

## 📊 分页查询

### 通用分页

```go
type PaginationRequest struct {
    Page     int `form:"page" binding:"min=1"`
    PageSize int `form:"page_size" binding:"min=1,max=100"`
}

func ListUsers(db *gorm.DB, req *PaginationRequest) ([]model.User, int64, error) {
    var users []model.User
    var total int64
    
    // 计算总数
    if err := db.Model(&model.User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // 分页查询
    offset := (req.Page - 1) * req.PageSize
    err := db.Offset(offset).
        Limit(req.PageSize).
        Order("created_at DESC").
        Find(&users).Error
    
    return users, total, err
}
```

### 使用 BaseRepository

```go
// 使用框架提供的基类
type UserRepository struct {
    *base.BaseRepository
}

func (r *UserRepository) List(req *model.PaginationRequest) (*model.PaginationResult, error) {
    var users []model.User
    return r.Paginate(context.Background(), &users, req)
}
```

## 🔒 Hook（钩子）

### 生命周期钩子

```go
// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 密码加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

// AfterCreate 创建后钩子
func (u *User) AfterCreate(tx *gorm.DB) error {
    // 发送欢迎邮件
    log.Info("user created", "id", u.ID, "username", u.Username)
    return nil
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    // 验证逻辑
    if u.Status < 0 || u.Status > 2 {
        return errors.New("invalid status")
    }
    return nil
}

// AfterFind 查询后钩子
func (u *User) AfterFind(tx *gorm.DB) error {
    // 数据处理
    return nil
}
```

## 🎯 性能优化

### 1. 索引优化

```go
type User struct {
    Username string `gorm:"index"`                    // 普通索引
    Email    string `gorm:"uniqueIndex"`              // 唯一索引
    Status   int    `gorm:"index:idx_status_created"` // 组合索引
    CreatedAt time.Time `gorm:"index:idx_status_created"`
}
```

### 2. 只查询需要的字段

```go
// ❌ 不推荐：查询所有字段
db.Find(&users)

// ✅ 推荐：只查询需要的字段
db.Select("id", "username", "email").Find(&users)
```

### 3. 批量操作

```go
// 批量创建（性能更好）
db.CreateInBatches(users, 100)  // 每批 100 条

// 批量更新
db.Model(&model.User{}).
   Where("status = ?", 0).
   Updates(map[string]interface{}{"status": 1})
```

### 4. 预编译查询

```go
// 准备查询语句
stmt := db.Session(&gorm.Session{PrepareStmt: true})

// 多次执行
for i := 0; i < 10; i++ {
    stmt.Where("id = ?", i).First(&user)
}
```

## 🛡️ 数据安全

### 1. SQL 注入防护

```go
// ✅ 正确：使用参数化查询
db.Where("username = ?", username).First(&user)

// ❌ 错误：直接拼接 SQL（有注入风险）
db.Where("username = '" + username + "'").First(&user)
```

### 2. 软删除

```go
// 模型定义
type User struct {
    ID        uint64
    DeletedAt *time.Time `gorm:"index"`  // 软删除字段
}

// 删除（软删除）
db.Delete(&user)  // UPDATE users SET deleted_at = NOW()

// 查询会自动过滤已删除的记录
db.Find(&users)  // WHERE deleted_at IS NULL

// 包含已删除的记录
db.Unscoped().Find(&users)

// 永久删除
db.Unscoped().Delete(&user)
```

## 📚 实际示例

查看完整的数据库操作示例：

- **Repository 示例**: `services/admin-api/internal/repository/admin_repository.go`
- **模型定义**: `services/admin-api/internal/model/admin_user.go`
- **数据库配置**: `pkg/db/db.go`
- **使用示例**: `docs/demo/db.md`

## 🎯 下一步

- [认证授权](../features/authentication) - 学习用户认证
- [缓存系统](../features/cache) - 使用 Redis 缓存
- [基础类使用](../api-reference/base-classes) - 使用 BaseRepository

---

**提示**: GORM 功能强大，建议查阅 [GORM 官方文档](https://gorm.io/zh_CN/docs/) 了解更多高级特性。

