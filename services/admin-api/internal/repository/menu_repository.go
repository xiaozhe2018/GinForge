package repository

import (
	"goweb/services/admin-api/internal/model"

	"gorm.io/gorm"
)

// MenuRepository 菜单数据访问层
type MenuRepository struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单数据访问层实例
func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{
		db: db,
	}
}

// Create 创建菜单
func (r *MenuRepository) Create(menu *model.AdminMenu) error {
	return r.db.Create(menu).Error
}

// GetByID 根据ID获取菜单
func (r *MenuRepository) GetByID(id uint64) (*model.AdminMenu, error) {
	var menu model.AdminMenu
	err := r.db.Preload("Parent").Preload("Children").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// List 获取菜单列表
func (r *MenuRepository) List(req *model.AdminMenuListRequest) ([]model.AdminMenu, int64, error) {
	var menus []model.AdminMenu
	var total int64

	query := r.db.Model(&model.AdminMenu{})

	// 搜索条件
	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.ParentID != nil {
		query = query.Where("parent_id = ?", *req.ParentID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Parent").Preload("Children").Offset(offset).Limit(req.PageSize).Order("sort ASC, created_at DESC").Find(&menus).Error; err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

// GetTree 获取菜单树
func (r *MenuRepository) GetTree() ([]model.AdminMenu, error) {
	var menus []model.AdminMenu
	err := r.db.Where("status = ?", 1).Order("sort ASC, created_at ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	return r.buildMenuTree(menus, 0), nil
}

// GetTreeByRoleID 根据角色ID获取菜单树
func (r *MenuRepository) GetTreeByRoleID(roleID uint64) ([]model.AdminMenu, error) {
	var menus []model.AdminMenu
	err := r.db.Joins("JOIN gf_admin_role_menus ON gf_admin_menus.id = gf_admin_role_menus.menu_id").
		Where("gf_admin_role_menus.role_id = ? AND gf_admin_menus.status = ?", roleID, 1).
		Order("gf_admin_menus.sort ASC, gf_admin_menus.created_at ASC").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	return r.buildMenuTree(menus, 0), nil
}

// buildMenuTree 构建菜单树
func (r *MenuRepository) buildMenuTree(menus []model.AdminMenu, parentID uint64) []model.AdminMenu {
	var tree []model.AdminMenu
	for _, menu := range menus {
		if menu.ParentID == parentID {
			children := r.buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				menu.Children = children
			}
			tree = append(tree, menu)
		}
	}
	return tree
}

// Update 更新菜单
func (r *MenuRepository) Update(menu *model.AdminMenu) error {
	return r.db.Save(menu).Error
}

// Delete 删除菜单
func (r *MenuRepository) Delete(id uint64) error {
	return r.db.Delete(&model.AdminMenu{}, id).Error
}

// GetByCode 根据编码获取菜单
func (r *MenuRepository) GetByCode(code string) (*model.AdminMenu, error) {
	var menu model.AdminMenu
	err := r.db.Where("code = ?", code).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}
