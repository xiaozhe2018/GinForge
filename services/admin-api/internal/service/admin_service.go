package service

import (
	"context"
	"errors"
	"fmt"
	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/notification"
	pkgRedis "goweb/pkg/redis"
	"goweb/pkg/websocket"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminService 管理后台服务
type AdminService struct {
	userRepo           *repository.UserRepository
	roleRepo           *repository.RoleRepository
	menuRepo           *repository.MenuRepository
	permissionRepo     *repository.PermissionRepository
	logRepo            *repository.OperationLogRepository
	redisClient        *pkgRedis.Client
	logger             logger.Logger
	config             *config.Config
	notificationClient *notification.Client
}

// NewAdminService 创建管理后台服务实例
func NewAdminService(db *gorm.DB, cfg *config.Config, redisClient *pkgRedis.Client) *AdminService {
	var notifyClient *notification.Client
	if redisClient != nil {
		notifyClient = notification.NewClient(redisClient)
	}
	
	return &AdminService{
		userRepo:           repository.NewUserRepository(db),
		roleRepo:           repository.NewRoleRepository(db),
		menuRepo:           repository.NewMenuRepository(db),
		permissionRepo:     repository.NewPermissionRepository(db),
		logRepo:            repository.NewOperationLogRepository(db),
		redisClient:        redisClient,
		config:             cfg,
		notificationClient: notifyClient,
	}
}

// SetLogger 设置日志器
func (s *AdminService) SetLogger(logger logger.Logger) {
	s.logger = logger
}

// UserService 用户服务
type UserService struct {
	*AdminService
	systemService *AdminSystemService
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB, cfg *config.Config, redisClient *pkgRedis.Client) *UserService {
	return &UserService{
		AdminService: NewAdminService(db, cfg, redisClient),
	}
}

// SetSystemService 设置系统服务
func (s *UserService) SetSystemService(systemService *AdminSystemService) {
	s.systemService = systemService
}

// Login 用户登录
func (s *UserService) Login(req *model.AdminUserLoginRequest, loginIP string) (*model.AdminUserLoginResponse, error) {
	ctx := context.Background()
	
	// 检查账户是否被锁定（使用Redis存储锁定信息）
	if s.systemService != nil && s.redisClient != nil && s.redisClient.IsEnabled() {
		lockKey := fmt.Sprintf("login:locked:%s", req.Username)
		locked, err := s.redisClient.Exists(ctx, lockKey)
		if err == nil && locked {
			// 获取剩余锁定时间
			ttl, _ := s.redisClient.TTL(ctx, lockKey).Result()
			return nil, fmt.Errorf("账户已被锁定，请在 %d 分钟后重试", int(ttl.Minutes())+1)
		}
	}
	
	// 根据用户名获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录失败尝试
			s.recordLoginFailure(ctx, req.Username)
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// 记录失败尝试
		s.recordLoginFailure(ctx, req.Username)
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}
	
	// 登录成功，清除失败记录
	s.clearLoginFailures(ctx, req.Username)

	// 更新用户登录信息（异步执行，不影响登录流程）
	go func() {
		if err := s.userRepo.UpdateLoginInfo(user.ID, loginIP); err != nil {
			s.logger.Error("update login info error", err, "user_id", user.ID)
		}
	}()

	// 记录登录日志
	go func() {
		username := user.Username
		userID := user.ID
		logEntry := &model.OperationLog{
			UserID:     &userID,
			Username:   &username,
			Method:     "POST",
			Path:       "/api/v1/admin/auth/login",
			IP:         &loginIP,
			StatusCode: 200,
			CreatedAt:  time.Now(),
		}
		if err := s.logRepo.Create(logEntry); err != nil {
			s.logger.Error("create login log error", err)
		}
	}()

	// 获取用户菜单
	var menus []model.AdminMenu
	if len(user.Roles) > 0 {
		menus, err = s.menuRepo.GetTreeByRoleID(user.Roles[0].ID)
		if err != nil {
			s.logger.Error("get user menus error", err)
		}
	}

	// 获取用户权限
	permissions, err := s.permissionRepo.GetCodesByUserID(user.ID)
	if err != nil {
		s.logger.Error("get user permissions error", err)
	}

	// 生成JWT Token
	// 获取会话超时时间（从系统配置读取）
	sessionTimeout := 24 * 60 // 默认24小时（分钟）
	if s.systemService != nil {
		sessionTimeout = s.systemService.GetSessionTimeout(ctx)
	}
	
	claims := jwt.MapClaims{
		"user_id":  fmt.Sprintf("%d", user.ID),
		"username": user.Username,
		"exp":      time.Now().Add(time.Duration(sessionTimeout) * time.Minute).Unix(),
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenObj.SignedString([]byte(s.config.GetString("jwt.secret")))
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	s.logger.Info("user login success", "user_id", user.ID, "username", user.Username, "ip", loginIP)

	// 发送登录通知到 WebSocket（异步，延迟1秒等待前端建立连接）
	if s.notificationClient != nil {
		s.logger.Info("notification client is available, will send login notification after 1s")
		go func() {
			// 延迟1秒，等待前端建立 WebSocket 连接
			time.Sleep(time.Second)
			
			ctx := context.Background()
			loginTime := time.Now().Format("2006-01-02 15:04:05")
			
			notification := &websocket.NotificationMessage{
				Title: "用户登录提醒",
				Body:  fmt.Sprintf("%s 于 %s 登录了系统", user.Username, loginTime),
				Icon:  "User",
				Link:  "/system/users",
			}
			
			s.logger.Info("sending login notification", "username", user.Username, "notification", notification)
			
			// 广播给所有在线管理员
			if err := s.notificationClient.BroadcastNotification(ctx, notification); err != nil {
				s.logger.Error("failed to send login notification", err)
			} else {
				s.logger.Info("login notification sent successfully", "username", user.Username)
			}
		}()
	} else {
		s.logger.Warn("notification client is nil, cannot send login notification")
	}

	return &model.AdminUserLoginResponse{
		Token:       token,
		User:        *user,
		Menus:       menus,
		Permissions: permissions,
	}, nil
}

// Logout 用户登出
func (s *UserService) Logout(userID string, username string, token string, ip string) error {
	s.logger.Info("user logout", "user_id", userID, "username", username, "ip", ip)

	// 将token加入黑名单（过期时间24小时，与JWT过期时间一致）
	if s.redisClient != nil {
		ctx := context.Background()
		blacklistKey := fmt.Sprintf("token:blacklist:%s", token)

		// 设置到Redis，24小时后自动过期
		if err := s.redisClient.Set(ctx, blacklistKey, "1", 24*time.Hour); err != nil {
			s.logger.Error("add token to blacklist error", err)
			// 不影响登出流程，继续执行
		} else {
			s.logger.Info("token added to blacklist", "user_id", userID)
		}
	}

	// 记录登出日志到操作日志表
	go func() {
		// 转换userID为uint64
		var uid uint64
		fmt.Sscanf(userID, "%d", &uid)

		logEntry := &model.OperationLog{
			UserID:     &uid,
			Username:   &username,
			Method:     "POST",
			Path:       "/api/v1/admin/auth/logout",
			IP:         &ip,
			StatusCode: 200,
			CreatedAt:  time.Now(),
		}
		if err := s.logRepo.Create(logEntry); err != nil {
			s.logger.Error("create logout log error", err)
		}
	}()

	return nil
}

// IsTokenBlacklisted 检查token是否在黑名单中
func (s *UserService) IsTokenBlacklisted(token string) bool {
	if s.redisClient == nil {
		return false
	}

	ctx := context.Background()
	blacklistKey := fmt.Sprintf("token:blacklist:%s", token)

	exists, err := s.redisClient.Exists(ctx, blacklistKey)
	if err != nil {
		s.logger.Error("check token blacklist error", err)
		return false
	}

	return exists
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *model.AdminUserCreateRequest) error {
	// 检查用户名是否已存在
	_, err := s.userRepo.GetByUsername(req.Username)
	if err == nil {
		return errors.New("用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 检查邮箱是否已存在
	_, err = s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return errors.New("邮箱已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 验证密码是否符合安全策略（从数据库配置读取）
	if s.systemService != nil {
		ctx := context.Background()
		if err := s.systemService.ValidatePassword(ctx, req.Password); err != nil {
			return err
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &model.AdminUser{
		Username: req.Username,
		Email:    req.Email,
		Phone:    &req.Phone,
		Password: string(hashedPassword),
		Name:     &req.Name,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// 分配角色
	if err := s.userRepo.UpdateRoles(user.ID, req.RoleIDs); err != nil {
		return err
	}

	return nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id uint64, req *model.AdminUserUpdateRequest) error {
	// 获取用户
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 检查邮箱是否已被其他用户使用
	if user.Email != req.Email {
		_, err := s.userRepo.GetByEmail(req.Email)
		if err == nil {
			return errors.New("邮箱已被其他用户使用")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	// 更新用户信息
	user.Email = req.Email
	user.Phone = &req.Phone
	user.Name = &req.Name
	user.Status = req.Status

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// 更新角色
	if err := s.userRepo.UpdateRoles(user.ID, req.RoleIDs); err != nil {
		return err
	}

	return nil
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers(req *model.AdminUserListRequest) (*model.AdminUserListResponse, error) {
	users, total, err := s.userRepo.List(req)
	if err != nil {
		return nil, err
	}

	return &model.AdminUserListResponse{
		List:  users,
		Total: total,
	}, nil
}

// GetUser 获取用户详情
func (s *UserService) GetUser(id uint64) (*model.AdminUser, error) {
	return s.userRepo.GetByID(id)
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(id uint64, status int8) error {
	return s.userRepo.UpdateStatus(id, status)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint64) error {
	return s.userRepo.Delete(id)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint64, req *model.ChangePasswordRequest) error {
	// 获取用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 验证新密码是否符合安全策略（从数据库配置读取）
	if s.systemService != nil {
		ctx := context.Background()
		if err := s.systemService.ValidatePassword(ctx, req.NewPassword); err != nil {
			return err
		}
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

// RoleService 角色服务
type RoleService struct {
	*AdminService
}

// NewRoleService 创建角色服务实例
func NewRoleService(db *gorm.DB, cfg *config.Config, redisClient *pkgRedis.Client) *RoleService {
	return &RoleService{
		AdminService: NewAdminService(db, cfg, redisClient),
	}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(req *model.AdminRoleCreateRequest) error {
	// 检查角色编码是否已存在
	_, err := s.roleRepo.GetByCode(req.Code)
	if err == nil {
		return errors.New("角色编码已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建角色
	role := &model.AdminRole{
		Name:        req.Name,
		Code:        req.Code,
		Description: &req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return err
	}

	// 分配权限
	if err := s.roleRepo.UpdatePermissions(role.ID, req.PermissionIDs); err != nil {
		return err
	}

	// 分配菜单
	if err := s.roleRepo.UpdateMenus(role.ID, req.MenuIDs); err != nil {
		return err
	}

	return nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(id uint64, req *model.AdminRoleUpdateRequest) error {
	// 获取角色
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 检查角色编码是否已被其他角色使用
	if role.Code != req.Code {
		_, err := s.roleRepo.GetByCode(req.Code)
		if err == nil {
			return errors.New("角色编码已被其他角色使用")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	// 更新角色信息
	role.Name = req.Name
	role.Code = req.Code
	role.Description = &req.Description
	role.Sort = req.Sort
	role.Status = req.Status

	if err := s.roleRepo.Update(role); err != nil {
		return err
	}

	// 更新权限
	if err := s.roleRepo.UpdatePermissions(role.ID, req.PermissionIDs); err != nil {
		return err
	}

	// 更新菜单
	if err := s.roleRepo.UpdateMenus(role.ID, req.MenuIDs); err != nil {
		return err
	}

	return nil
}

// GetRoles 获取角色列表
func (s *RoleService) GetRoles(req *model.AdminRoleListRequest) (*model.AdminRoleListResponse, error) {
	roles, total, err := s.roleRepo.List(req)
	if err != nil {
		return nil, err
	}

	return &model.AdminRoleListResponse{
		List:  roles,
		Total: total,
	}, nil
}

// GetRole 获取角色详情
func (s *RoleService) GetRole(id uint64) (*model.AdminRole, error) {
	return s.roleRepo.GetByID(id)
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint64) error {
	return s.roleRepo.Delete(id)
}

// recordLoginFailure 记录登录失败
func (s *UserService) recordLoginFailure(ctx context.Context, username string) {
	if s.redisClient == nil || !s.redisClient.IsEnabled() || s.systemService == nil {
		return
	}
	
	failureKey := fmt.Sprintf("login:failures:%s", username)
	lockKey := fmt.Sprintf("login:locked:%s", username)
	
	// 增加失败次数
	failures, err := s.redisClient.Incr(ctx, failureKey).Result()
	if err != nil {
		s.logger.Error("Failed to record login failure", "error", err, "username", username)
		return
	}
	
	// 设置失败记录过期时间（登录锁定时长）
	lockoutDuration := s.systemService.GetLockoutDuration(ctx)
	s.redisClient.Expire(ctx, failureKey, time.Duration(lockoutDuration)*time.Minute).Result()
	
	// 检查是否达到最大失败次数
	maxAttempts := s.systemService.GetMaxLoginAttempts(ctx)
	if int(failures) >= maxAttempts {
		// 锁定账户
		if err := s.redisClient.Set(ctx, lockKey, "1", time.Duration(lockoutDuration)*time.Minute); err != nil {
			s.logger.Error("Failed to lock account", "error", err, "username", username)
			return
		}
		s.logger.Warn("Account locked due to too many failed login attempts", 
			"username", username, 
			"failures", failures,
			"lockout_duration", lockoutDuration)
	}
}

// clearLoginFailures 清除登录失败记录
func (s *UserService) clearLoginFailures(ctx context.Context, username string) {
	if s.redisClient == nil || !s.redisClient.IsEnabled() {
		return
	}
	
	failureKey := fmt.Sprintf("login:failures:%s", username)
	s.redisClient.Del(ctx, failureKey)
}
