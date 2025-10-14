package service

import (
	"errors"
	"goweb/pkg/config"
	pkgRedis "goweb/pkg/redis"
	"goweb/services/admin-api/internal/model"

	"gorm.io/gorm"
)

// MenuService 菜单服务
type MenuService struct {
	*AdminService
}

// NewMenuService 创建菜单服务实例
func NewMenuService(db *gorm.DB, cfg *config.Config, redisClient *pkgRedis.Client) *MenuService {
	return &MenuService{
		AdminService: NewAdminService(db, cfg, redisClient),
	}
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(req *model.AdminMenuCreateRequest) error {
	// 检查菜单编码是否已存在
	_, err := s.menuRepo.GetByCode(req.Code)
	if err == nil {
		return errors.New("菜单编码已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建菜单
	menu := &model.AdminMenu{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Sort:        req.Sort,
		Visible:     req.Visible,
		Status:      req.Status,
		Description: &req.Description,
	}

	if req.Path != "" {
		menu.Path = &req.Path
	}
	if req.Component != "" {
		menu.Component = &req.Component
	}
	if req.Icon != "" {
		menu.Icon = &req.Icon
	}

	if err := s.menuRepo.Create(menu); err != nil {
		return err
	}

	return nil
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(id uint64, req *model.AdminMenuUpdateRequest) error {
	// 获取菜单
	menu, err := s.menuRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 检查菜单编码是否已被其他菜单使用
	if menu.Code != req.Code {
		_, err := s.menuRepo.GetByCode(req.Code)
		if err == nil {
			return errors.New("菜单编码已被其他菜单使用")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	// 更新菜单信息
	menu.ParentID = req.ParentID
	menu.Name = req.Name
	menu.Code = req.Code
	menu.Type = req.Type
	menu.Sort = req.Sort
	menu.Visible = req.Visible
	menu.Status = req.Status
	menu.Description = &req.Description

	if req.Path != "" {
		menu.Path = &req.Path
	} else {
		menu.Path = nil
	}
	if req.Component != "" {
		menu.Component = &req.Component
	} else {
		menu.Component = nil
	}
	if req.Icon != "" {
		menu.Icon = &req.Icon
	} else {
		menu.Icon = nil
	}

	if err := s.menuRepo.Update(menu); err != nil {
		return err
	}

	return nil
}

// GetMenus 获取菜单列表
func (s *MenuService) GetMenus(req *model.AdminMenuListRequest) (*model.AdminMenuListResponse, error) {
	menus, total, err := s.menuRepo.List(req)
	if err != nil {
		return nil, err
	}

	return &model.AdminMenuListResponse{
		List:  menus,
		Total: total,
	}, nil
}

// GetMenuTree 获取菜单树
func (s *MenuService) GetMenuTree() (*model.AdminMenuTreeResponse, error) {
	menus, err := s.menuRepo.GetTree()
	if err != nil {
		return nil, err
	}

	return &model.AdminMenuTreeResponse{
		List: menus,
	}, nil
}

// GetMenu 获取菜单详情
func (s *MenuService) GetMenu(id uint64) (*model.AdminMenu, error) {
	return s.menuRepo.GetByID(id)
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint64) error {
	return s.menuRepo.Delete(id)
}

// PermissionService 权限服务
type PermissionService struct {
	*AdminService
}

// NewPermissionService 创建权限服务实例
func NewPermissionService(db *gorm.DB, cfg *config.Config, redisClient *pkgRedis.Client) *PermissionService {
	return &PermissionService{
		AdminService: NewAdminService(db, cfg, redisClient),
	}
}

// CreatePermission 创建权限
func (s *PermissionService) CreatePermission(req *model.AdminPermissionCreateRequest) error {
	// 检查权限编码是否已存在
	_, err := s.permissionRepo.GetByCode(req.Code)
	if err == nil {
		return errors.New("权限编码已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建权限
	permission := &model.AdminPermission{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Description: &req.Description,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return err
	}

	return nil
}

// UpdatePermission 更新权限
func (s *PermissionService) UpdatePermission(id uint64, req *model.AdminPermissionUpdateRequest) error {
	// 获取权限
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 检查权限编码是否已被其他权限使用
	if permission.Code != req.Code {
		_, err := s.permissionRepo.GetByCode(req.Code)
		if err == nil {
			return errors.New("权限编码已被其他权限使用")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	// 更新权限信息
	permission.Name = req.Name
	permission.Code = req.Code
	permission.Type = req.Type
	permission.Description = &req.Description

	if err := s.permissionRepo.Update(permission); err != nil {
		return err
	}

	return nil
}

// GetPermissions 获取权限列表
func (s *PermissionService) GetPermissions(req *model.AdminPermissionListRequest) (*model.AdminPermissionListResponse, error) {
	permissions, total, err := s.permissionRepo.List(req)
	if err != nil {
		return nil, err
	}

	return &model.AdminPermissionListResponse{
		List:  permissions,
		Total: total,
	}, nil
}

// GetPermission 获取权限详情
func (s *PermissionService) GetPermission(id uint64) (*model.AdminPermission, error) {
	return s.permissionRepo.GetByID(id)
}

// DeletePermission 删除权限
func (s *PermissionService) DeletePermission(id uint64) error {
	return s.permissionRepo.Delete(id)
}
