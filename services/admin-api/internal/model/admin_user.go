package model

import (
	"time"
)

// AdminUser 管理员用户模型
type AdminUser struct {
	ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string     `json:"username" gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名"`
	Email       string     `json:"email" gorm:"type:varchar(100);uniqueIndex;not null;comment:邮箱"`
	Phone       *string    `json:"phone" gorm:"type:varchar(20);index;comment:手机号"`
	Password    string     `json:"-" gorm:"type:varchar(255);not null;comment:密码"`
	Name        *string    `json:"name" gorm:"type:varchar(50);comment:真实姓名"`
	Avatar      *string    `json:"avatar" gorm:"type:varchar(255);comment:头像URL"`
	Status      int8       `json:"status" gorm:"type:tinyint(1);default:1;index;comment:状态:1-启用,0-禁用"`
	LastLoginAt *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"`
	LastLoginIP *string    `json:"last_login_ip" gorm:"type:varchar(45);comment:最后登录IP"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime;index"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`

	// 关联 - 手动加载，不使用GORM自动关联
	Roles []AdminRole `json:"roles,omitempty" gorm:"-"`
}

// TableName 返回表名
func (AdminUser) TableName() string {
	return "gf_admin_users"
}

// AdminUserCreateRequest 创建用户请求
type AdminUserCreateRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=20"`
	Email    string   `json:"email" binding:"required,email"`
	Phone    string   `json:"phone" binding:"omitempty,len=11"`
	Password string   `json:"password" binding:"required,min=6,max=20"`
	Name     string   `json:"name" binding:"required"`
	RoleIDs  []uint64 `json:"role_ids" binding:"required"`
}

// AdminUserUpdateRequest 更新用户请求
type AdminUserUpdateRequest struct {
	Email   string   `json:"email" binding:"required,email"`
	Phone   string   `json:"phone" binding:"omitempty,len=11"`
	Name    string   `json:"name" binding:"required"`
	Status  int8     `json:"status" binding:"oneof=0 1"`
	RoleIDs []uint64 `json:"role_ids" binding:"required"`
}

// AdminUserListRequest 用户列表请求
type AdminUserListRequest struct {
	Page     int     `form:"page"`                        // 页码，不传时使用默认值1
	PageSize int     `form:"page_size" binding:"max=100"` // 每页数量，不传时使用默认值10
	Keyword  string  `form:"keyword"`                     // 搜索关键词
	Status   *int8   `form:"status"`                      // 用户状态
	RoleID   *uint64 `form:"role_id"`                     // 角色ID
}

// AdminUserListResponse 用户列表响应
type AdminUserListResponse struct {
	List  []AdminUser `json:"list"`
	Total int64       `json:"total"`
}

// AdminUserLoginRequest 用户登录请求
type AdminUserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Captcha  string `json:"captcha"`
}

// AdminUserLoginResponse 用户登录响应
type AdminUserLoginResponse struct {
	Token       string      `json:"token"`
	User        AdminUser   `json:"user"`
	Menus       []AdminMenu `json:"menus"`
	Permissions []string    `json:"permissions"`
}
