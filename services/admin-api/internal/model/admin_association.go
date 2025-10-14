package model

import (
	"time"
)

// AdminUserRole 用户角色关联表
type AdminUserRole struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `json:"user_id" gorm:"type:bigint(20) unsigned;not null;uniqueIndex:uk_user_role;index;comment:用户ID"`
	RoleID    uint64    `json:"role_id" gorm:"type:bigint(20) unsigned;not null;uniqueIndex:uk_user_role;index;comment:角色ID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 返回表名
func (AdminUserRole) TableName() string {
	return "admin_user_roles"
}

// AdminRolePermission 角色权限关联表
type AdminRolePermission struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID       uint64    `json:"role_id" gorm:"type:bigint(20) unsigned;not null;uniqueIndex:uk_role_permission;index;comment:角色ID"`
	PermissionID uint64    `json:"permission_id" gorm:"type:bigint(20) unsigned;not null;uniqueIndex:uk_role_permission;index;comment:权限ID"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 返回表名
func (AdminRolePermission) TableName() string {
	return "admin_role_permissions"
}

// AdminRoleMenu 角色菜单关联表
type AdminRoleMenu struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID    uint64    `json:"role_id" gorm:"type:bigint(20) unsigned;not null;uniqueIndex:uk_role_menu;index;comment:角色ID"`
	MenuID    uint64    `json:"menu_id" gorm:"type:bigint(20) unsigned;not null;uniqueIndex:uk_role_menu;index;comment:菜单ID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 返回表名
func (AdminRoleMenu) TableName() string {
	return "admin_role_menus"
}
