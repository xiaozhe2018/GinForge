package model

import (
	"time"
)

// AdminMenu 菜单模型
type AdminMenu struct {
	ID          uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentID    uint64     `json:"parent_id" gorm:"type:bigint(20) unsigned;default:0;index;comment:父菜单ID"`
	Name        string     `json:"name" gorm:"type:varchar(50);not null;index;comment:菜单名称"`
	Code        string     `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:菜单编码"`
	Type        string     `json:"type" gorm:"type:varchar(20);default:menu;index;comment:菜单类型:directory-目录,menu-菜单,button-按钮"`
	Path        *string    `json:"path" gorm:"type:varchar(200);comment:路由路径"`
	Component   *string    `json:"component" gorm:"type:varchar(200);comment:组件路径"`
	Icon        *string    `json:"icon" gorm:"type:varchar(50);comment:菜单图标"`
	Sort        int        `json:"sort" gorm:"type:int(11);default:0;index;comment:排序"`
	Visible     int8       `json:"visible" gorm:"type:tinyint(1);default:1;comment:是否显示:1-显示,0-隐藏"`
	Status      int8       `json:"status" gorm:"type:tinyint(1);default:1;index;comment:状态:1-启用,0-禁用"`
	Description *string    `json:"description" gorm:"type:text;comment:菜单描述"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`

	// 关联
	Parent   *AdminMenu  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []AdminMenu `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Roles    []AdminRole `json:"roles,omitempty" gorm:"many2many:admin_role_menus;"`
}

// TableName 返回表名
func (AdminMenu) TableName() string {
	return "admin_menus"
}

// AdminMenuCreateRequest 创建菜单请求
type AdminMenuCreateRequest struct {
	ParentID    uint64 `json:"parent_id" binding:"min=0"`
	Name        string `json:"name" binding:"required,min=2,max=20"`
	Code        string `json:"code" binding:"required,min=2,max=20"`
	Type        string `json:"type" binding:"required,oneof=directory menu button"`
	Path        string `json:"path"`
	Component   string `json:"component"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort" binding:"min=0"`
	Visible     int8   `json:"visible" binding:"oneof=0 1"`
	Status      int8   `json:"status" binding:"oneof=0 1"`
	Description string `json:"description"`
}

// AdminMenuUpdateRequest 更新菜单请求
type AdminMenuUpdateRequest struct {
	ParentID    uint64 `json:"parent_id" binding:"min=0"`
	Name        string `json:"name" binding:"required,min=2,max=20"`
	Code        string `json:"code" binding:"required,min=2,max=20"`
	Type        string `json:"type" binding:"required,oneof=directory menu button"`
	Path        string `json:"path"`
	Component   string `json:"component"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort" binding:"min=0"`
	Visible     int8   `json:"visible" binding:"oneof=0 1"`
	Status      int8   `json:"status" binding:"oneof=0 1"`
	Description string `json:"description"`
}

// AdminMenuListRequest 菜单列表请求
type AdminMenuListRequest struct {
	Page     int     `form:"page" binding:"min=1"`
	PageSize int     `form:"page_size" binding:"min=1,max=100"`
	Keyword  string  `form:"keyword"`
	Type     string  `form:"type"`
	Status   *int8   `form:"status"`
	ParentID *uint64 `form:"parent_id"`
}

// AdminMenuListResponse 菜单列表响应
type AdminMenuListResponse struct {
	List  []AdminMenu `json:"list"`
	Total int64       `json:"total"`
}

// AdminMenuTreeResponse 菜单树响应
type AdminMenuTreeResponse struct {
	List []AdminMenu `json:"list"`
}
