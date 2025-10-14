package model

import (
	"time"
)

// AdminRole 角色模型
type AdminRole struct {
	ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"type:varchar(50);not null;index;comment:角色名称"`
	Code        string     `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:角色编码"`
	Description *string    `json:"description" gorm:"type:text;comment:角色描述"`
	Sort        int        `json:"sort" gorm:"type:int(11);default:0;index;comment:排序"`
	Status      int8       `json:"status" gorm:"type:tinyint(1);default:1;index;comment:状态:1-启用,0-禁用"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`

	// 关联 - 手动加载，不使用GORM自动关联
	Users       []AdminUser       `json:"users,omitempty" gorm:"-"`
	Permissions []AdminPermission `json:"permissions,omitempty" gorm:"-"`
	Menus       []AdminMenu       `json:"menus,omitempty" gorm:"-"`
}

// TableName 返回表名
func (AdminRole) TableName() string {
	return "admin_roles"
}

// AdminRoleCreateRequest 创建角色请求
type AdminRoleCreateRequest struct {
	Name          string   `json:"name" binding:"required,min=2,max=20"`
	Code          string   `json:"code" binding:"required,min=2,max=20"`
	Description   string   `json:"description"`
	Sort          int      `json:"sort" binding:"min=0"`
	Status        int8     `json:"status" binding:"oneof=0 1"`
	PermissionIDs []uint64 `json:"permission_ids"`
	MenuIDs       []uint64 `json:"menu_ids"`
}

// AdminRoleUpdateRequest 更新角色请求
type AdminRoleUpdateRequest struct {
	Name          string   `json:"name" binding:"required,min=2,max=20"`
	Code          string   `json:"code" binding:"required,min=2,max=20"`
	Description   string   `json:"description"`
	Sort          int      `json:"sort" binding:"min=0"`
	Status        int8     `json:"status" binding:"oneof=0 1"`
	PermissionIDs []uint64 `json:"permission_ids"`
	MenuIDs       []uint64 `json:"menu_ids"`
}

// AdminRoleListRequest 角色列表请求
type AdminRoleListRequest struct {
	Page     int    `form:"page"`                        // 页码，不传时使用默认值1
	PageSize int    `form:"page_size" binding:"max=100"` // 每页数量，不传时使用默认值10
	Keyword  string `form:"keyword"`                     // 搜索关键词
	Status   *int8  `form:"status"`                      // 角色状态
}

// AdminRoleListResponse 角色列表响应
type AdminRoleListResponse struct {
	List  []AdminRole `json:"list"`
	Total int64       `json:"total"`
}
