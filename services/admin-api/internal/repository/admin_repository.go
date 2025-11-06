package repository

import (
	"goweb/services/admin-api/internal/model"
	"time"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户数据访问层实例
func NewUserRepository(database *gorm.DB) *UserRepository {
	return &UserRepository{
		db: database,
	}
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	// 手动加载角色，使用JOIN查询
	var roles []model.AdminRole
	err = r.db.Table("gf_admin_roles").
		Select("gf_admin_roles.*").
		Joins("JOIN gf_admin_user_roles ON gf_admin_roles.id = gf_admin_user_roles.role_id").
		Where("gf_admin_user_roles.user_id = ?", user.ID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	// 手动加载角色
	var roles []model.AdminRole
	err = r.db.Table("gf_admin_roles").
		Select("gf_admin_roles.*").
		Joins("JOIN gf_admin_user_roles ON gf_admin_roles.id = gf_admin_user_roles.role_id").
		Where("gf_admin_user_roles.user_id = ?", user.ID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint64) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	// 手动加载角色
	var roles []model.AdminRole
	err = r.db.Table("gf_admin_roles").
		Select("gf_admin_roles.*").
		Joins("JOIN gf_admin_user_roles ON gf_admin_roles.id = gf_admin_user_roles.role_id").
		Where("gf_admin_user_roles.user_id = ?", user.ID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *model.AdminUser) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.AdminUser) error {
	return r.db.Save(user).Error
}

// UpdateLoginInfo 更新用户登录信息
func (r *UserRepository) UpdateLoginInfo(userID uint64, loginIP string) error {
	now := time.Now()
	return r.db.Model(&model.AdminUser{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"last_login_ip": loginIP,
		}).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id uint64) error {
	return r.db.Delete(&model.AdminUser{}, id).Error
}

// List 获取用户列表
func (r *UserRepository) List(req *model.AdminUserListRequest) ([]model.AdminUser, int64, error) {
	var users []model.AdminUser
	var total int64

	query := r.db.Model(&model.AdminUser{})

	// 搜索条件
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR name LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.RoleID != nil {
		query = query.Joins("JOIN gf_admin_user_roles ON gf_admin_users.id = gf_admin_user_roles.user_id").
			Where("gf_admin_user_roles.role_id = ?", *req.RoleID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 手动加载每个用户的角色
	for i := range users {
		var roles []model.AdminRole
		err := r.db.Table("gf_admin_roles").
			Select("gf_admin_roles.*").
			Joins("JOIN gf_admin_user_roles ON gf_admin_roles.id = gf_admin_user_roles.role_id").
			Where("gf_admin_user_roles.user_id = ?", users[i].ID).
			Find(&roles).Error
		if err != nil {
			return nil, 0, err
		}
		users[i].Roles = roles
	}

	return users, total, nil
}

// UpdateStatus 更新用户状态
func (r *UserRepository) UpdateStatus(id uint64, status int8) error {
	return r.db.Model(&model.AdminUser{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateRoles 更新用户角色
func (r *UserRepository) UpdateRoles(userID uint64, roleIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除现有角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&model.AdminUserRole{}).Error; err != nil {
			return err
		}

		// 添加新角色关联
		if len(roleIDs) > 0 {
			var userRoles []model.AdminUserRole
			for _, roleID := range roleIDs {
				userRoles = append(userRoles, model.AdminUserRole{
					UserID: userID,
					RoleID: roleID,
				})
			}
			if err := tx.Create(&userRoles).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// RoleRepository 角色数据访问层
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色数据访问层实例
func NewRoleRepository(database *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: database,
	}
}

// GetByID 根据ID获取角色
func (r *RoleRepository) GetByID(id uint64) (*model.AdminRole, error) {
	var role model.AdminRole
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}

	// 手动加载权限
	var permissions []model.AdminPermission
	err = r.db.Table("gf_admin_permissions").
		Select("gf_admin_permissions.*").
		Joins("JOIN gf_admin_role_permissions ON gf_admin_permissions.id = gf_admin_role_permissions.permission_id").
		Where("gf_admin_role_permissions.role_id = ?", role.ID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions

	// 手动加载菜单
	var menus []model.AdminMenu
	err = r.db.Table("gf_admin_menus").
		Select("gf_admin_menus.*").
		Joins("JOIN gf_admin_role_menus ON gf_admin_menus.id = gf_admin_role_menus.menu_id").
		Where("gf_admin_role_menus.role_id = ?", role.ID).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	role.Menus = menus

	return &role, nil
}

// Create 创建角色
func (r *RoleRepository) Create(role *model.AdminRole) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(role *model.AdminRole) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *RoleRepository) Delete(id uint64) error {
	return r.db.Delete(&model.AdminRole{}, id).Error
}

// List 获取角色列表
func (r *RoleRepository) List(req *model.AdminRoleListRequest) ([]model.AdminRole, int64, error) {
	var roles []model.AdminRole
	var total int64

	query := r.db.Model(&model.AdminRole{})

	// 搜索条件
	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("sort ASC, created_at DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	// 手动加载每个角色的权限和菜单
	for i := range roles {
		// 加载权限
		var permissions []model.AdminPermission
		err := r.db.Table("gf_admin_permissions").
			Select("gf_admin_permissions.*").
			Joins("JOIN gf_admin_role_permissions ON gf_admin_permissions.id = gf_admin_role_permissions.permission_id").
			Where("gf_admin_role_permissions.role_id = ?", roles[i].ID).
			Find(&permissions).Error
		if err != nil {
			return nil, 0, err
		}
		roles[i].Permissions = permissions

		// 加载菜单
		var menus []model.AdminMenu
		err = r.db.Table("gf_admin_menus").
			Select("gf_admin_menus.*").
			Joins("JOIN gf_admin_role_menus ON gf_admin_menus.id = gf_admin_role_menus.menu_id").
			Where("gf_admin_role_menus.role_id = ?", roles[i].ID).
			Find(&menus).Error
		if err != nil {
			return nil, 0, err
		}
		roles[i].Menus = menus
	}

	return roles, total, nil
}

// UpdatePermissions 更新角色权限
func (r *RoleRepository) UpdatePermissions(roleID uint64, permissionIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除现有权限关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.AdminRolePermission{}).Error; err != nil {
			return err
		}

		// 添加新权限关联
		if len(permissionIDs) > 0 {
			var rolePermissions []model.AdminRolePermission
			for _, permissionID := range permissionIDs {
				rolePermissions = append(rolePermissions, model.AdminRolePermission{
					RoleID:       roleID,
					PermissionID: permissionID,
				})
			}
			if err := tx.Create(&rolePermissions).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateMenus 更新角色菜单
func (r *RoleRepository) UpdateMenus(roleID uint64, menuIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除现有菜单关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.AdminRoleMenu{}).Error; err != nil {
			return err
		}

		// 添加新菜单关联
		if len(menuIDs) > 0 {
			var roleMenus []model.AdminRoleMenu
			for _, menuID := range menuIDs {
				roleMenus = append(roleMenus, model.AdminRoleMenu{
					RoleID: roleID,
					MenuID: menuID,
				})
			}
			if err := tx.Create(&roleMenus).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetByCode 根据编码获取角色
func (r *RoleRepository) GetByCode(code string) (*model.AdminRole, error) {
	var role model.AdminRole
	err := r.db.Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
