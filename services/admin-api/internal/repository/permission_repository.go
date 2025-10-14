package repository

import (
	"goweb/services/admin-api/internal/model"

	"gorm.io/gorm"
)

// PermissionRepository 权限数据访问层
type PermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限数据访问层实例
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

// Create 创建权限
func (r *PermissionRepository) Create(permission *model.AdminPermission) error {
	return r.db.Create(permission).Error
}

// GetByID 根据ID获取权限
func (r *PermissionRepository) GetByID(id uint64) (*model.AdminPermission, error) {
	var permission model.AdminPermission
	err := r.db.Preload("Roles").First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// List 获取权限列表
func (r *PermissionRepository) List(req *model.AdminPermissionListRequest) ([]model.AdminPermission, int64, error) {
	var permissions []model.AdminPermission
	var total int64

	query := r.db.Model(&model.AdminPermission{})

	// 搜索条件
	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// Update 更新权限
func (r *PermissionRepository) Update(permission *model.AdminPermission) error {
	return r.db.Save(permission).Error
}

// Delete 删除权限
func (r *PermissionRepository) Delete(id uint64) error {
	return r.db.Delete(&model.AdminPermission{}, id).Error
}

// GetByCode 根据编码获取权限
func (r *PermissionRepository) GetByCode(code string) (*model.AdminPermission, error) {
	var permission model.AdminPermission
	err := r.db.Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetByUserID 根据用户ID获取权限列表
func (r *PermissionRepository) GetByUserID(userID uint64) ([]model.AdminPermission, error) {
	var permissions []model.AdminPermission
	err := r.db.Joins("JOIN admin_role_permissions ON admin_permissions.id = admin_role_permissions.permission_id").
		Joins("JOIN admin_user_roles ON admin_role_permissions.role_id = admin_user_roles.role_id").
		Where("admin_user_roles.user_id = ?", userID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetCodesByUserID 根据用户ID获取权限编码列表
func (r *PermissionRepository) GetCodesByUserID(userID uint64) ([]string, error) {
	var codes []string
	err := r.db.Model(&model.AdminPermission{}).
		Select("admin_permissions.code").
		Joins("JOIN admin_role_permissions ON admin_permissions.id = admin_role_permissions.permission_id").
		Joins("JOIN admin_user_roles ON admin_role_permissions.role_id = admin_user_roles.role_id").
		Where("admin_user_roles.user_id = ?", userID).
		Pluck("admin_permissions.code", &codes).Error
	if err != nil {
		return nil, err
	}
	return codes, nil
}
