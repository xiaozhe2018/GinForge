package model

import (
	"time"
)

// AdminPermission 权限模型
type AdminPermission struct {
	ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"type:varchar(50);not null;index;comment:权限名称"`
	Code        string     `json:"code" gorm:"type:varchar(100);uniqueIndex;not null;comment:权限编码"`
	Type        string     `json:"type" gorm:"type:varchar(20);default:menu;index;comment:权限类型:menu-菜单,button-按钮,api-接口"`
	Description *string    `json:"description" gorm:"type:text;comment:权限描述"`
	Status      int8       `json:"status" gorm:"type:tinyint(1);not null;default:1;index;comment:状态:0=禁用,1=启用"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`

	// 关联
	Roles []AdminRole `json:"roles,omitempty" gorm:"many2many:gf_admin_role_permissions;"`
}

// TableName 返回表名
func (AdminPermission) TableName() string {
	return "gf_admin_permissions"
}

// AdminPermissionCreateRequest 创建权限请求
type AdminPermissionCreateRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=20"`
	Code        string `json:"code" binding:"required,min=2,max=50"`
	Type        string `json:"type" binding:"required,oneof=menu button api"`
	Description string `json:"description"`
	Status      int8   `json:"status" binding:"oneof=0 1"` // 状态：0=禁用，1=启用
}

// AdminPermissionUpdateRequest 更新权限请求
type AdminPermissionUpdateRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=20"`
	Code        string `json:"code" binding:"required,min=2,max=50"`
	Type        string `json:"type" binding:"required,oneof=menu button api"`
	Description string `json:"description"`
	Status      int8   `json:"status" binding:"oneof=0 1"` // 状态：0=禁用，1=启用
}

// AdminPermissionListRequest 权限列表请求
type AdminPermissionListRequest struct {
	Page     int    `form:"page"`                        // 页码，不传时使用默认值1
	PageSize int    `form:"page_size" binding:"max=100"` // 每页数量，不传时使用默认值10
	Keyword  string `form:"keyword"`                     // 搜索关键词
	Type     string `form:"type"`                        // 权限类型
}

// AdminPermissionListResponse 权限列表响应
type AdminPermissionListResponse struct {
	List  []AdminPermission `json:"list"`
	Total int64             `json:"total"`
}
