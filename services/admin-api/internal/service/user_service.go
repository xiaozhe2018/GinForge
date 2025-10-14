package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/pkg/notification"
	pkgRedis "goweb/pkg/redis"
	"goweb/pkg/websocket"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/repository"
)

// UserService 用户服务
type UserService struct {
	userRepo           *repository.UserRepository
	roleRepo           *repository.RoleRepository
	menuRepo           *repository.MenuRepository
	permissionRepo     *repository.PermissionRepository
	redisClient        *pkgRedis.Client
	logger             logger.Logger
	config             *config.Config
	notificationClient *notification.Client
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB, cfg *config.Config, redisClient *pkgRedis.Client) *UserService {
	var notifyClient *notification.Client
	if redisClient != nil {
		notifyClient = notification.NewClient(redisClient)
	}
	
	return &UserService{
		userRepo:           repository.NewUserRepository(db),
		roleRepo:           repository.NewRoleRepository(db),
		menuRepo:           repository.NewMenuRepository(db),
		permissionRepo:     repository.NewPermissionRepository(db),
		redisClient:        redisClient,
		config:             cfg,
		notificationClient: notifyClient,
		logger:             logger.New("user-service", "info"),
	}
}

// SetLogger 设置日志器
func (s *UserService) SetLogger(log logger.Logger) {
	s.logger = log
}

// Login 用户登录
func (s *UserService) Login(req *model.AdminUserLoginRequest, loginIP string) (*model.AdminUserLoginResponse, error) {
	// 查询用户
	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// 获取用户角色
	roles, err := s.roleRepo.GetRolesByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	// 获取用户菜单
	menus, err := s.menuRepo.GetMenusByRoles(roles)
	if err != nil {
		return nil, err
	}

	// 获取用户权限
	permissions, err := s.permissionRepo.GetPermissionsByRoles(roles)
	if err != nil {
		return nil, err
	}

	// 更新登录信息
	user.LastLoginIP = loginIP
	user.LastLoginTime = time.Now()
	user.LoginCount++
	if err := s.userRepo.UpdateUser(user); err != nil {
		s.logger.Error("update user login info failed", err)
	}

	s.logger.Info("user login success", "user_id", user.ID, "username", user.Username, "ip", loginIP)

	// 发送登录通知到 WebSocket（异步，延迟1秒等待前端建立连接）
	if s.notificationClient != nil {
		s.logger.Info("notification client is available, will send login notification after 1s")
		go func() {
			// 延迟1秒，等待前端建立 WebSocket 连接
			time.Sleep(time.Second)
			
			ctx := context.Background()
			
			notification := &websocket.NotificationMessage{
				Title:    "用户登录提醒",
				Body:     fmt.Sprintf("%s 于 %s 登录了系统", user.Username, time.Now().Format("2006-01-02 15:04:05")),
				Icon:     "User",
				Link:     "/system/users",
				Category: "system",
			}
			
			s.logger.Info("sending login notification", "username", user.Username)
			
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
func (s *UserService) Logout(userID uint, token string) error {
	// 将token加入黑名单
	if s.redisClient != nil {
		// 获取token剩余有效期
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.config.GetString("jwt.secret")), nil
		})
		if err != nil {
			s.logger.Error("parse token failed", err)
		}

		// 计算过期时间
		var expTime time.Time
		if exp, ok := claims["exp"].(float64); ok {
			expTime = time.Unix(int64(exp), 0)
		} else {
			// 默认24小时
			expTime = time.Now().Add(24 * time.Hour)
		}

		// 将token加入黑名单
		duration := time.Until(expTime)
		if duration > 0 {
			ctx := context.Background()
			key := fmt.Sprintf("token:blacklist:%s", token)
			err := s.redisClient.Set(ctx, key, "1", duration)
			if err != nil {
				s.logger.Error("add token to blacklist failed", err)
				return err
			}
		}
	}

	return nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.AdminUser, error) {
	return s.userRepo.GetUserByID(id)
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers(query *model.AdminUserQuery) ([]*model.AdminUser, int64, error) {
	return s.userRepo.GetUsers(query)
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *model.AdminUser) error {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return errors.New("用户名已存在")
	}

	// 设置默认值
	if user.Status == 0 {
		user.Status = 1 // 默认启用
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 创建用户
	return s.userRepo.CreateUser(user)
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *model.AdminUser) error {
	// 检查用户是否存在
	existingUser, err := s.userRepo.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	// 检查用户名是否已存在（排除自己）
	if user.Username != existingUser.Username {
		otherUser, err := s.userRepo.GetUserByUsername(user.Username)
		if err == nil && otherUser != nil {
			return errors.New("用户名已存在")
		}
	}

	// 如果密码不为空，则更新密码
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	} else {
		// 保持原密码不变
		user.Password = existingUser.Password
	}

	// 更新用户
	return s.userRepo.UpdateUser(user)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.DeleteUser(id)
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(id uint, status int) error {
	return s.userRepo.UpdateUserStatus(id, status)
}

// UpdateUserProfile 更新用户个人资料
func (s *UserService) UpdateUserProfile(id uint, profile *model.AdminUserProfileUpdate) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	// 更新可修改的字段
	user.Nickname = profile.Nickname
	user.Avatar = profile.Avatar
	user.Email = profile.Email
	user.Phone = profile.Phone

	return s.userRepo.UpdateUser(user)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.userRepo.UpdateUser(user)
}

// generateToken 生成JWT令牌
func (s *UserService) generateToken(userID uint, username string) (string, error) {
	// 设置token过期时间
	expiresIn := s.config.GetDuration("jwt.expires_in")
	if expiresIn == 0 {
		expiresIn = 24 * time.Hour // 默认24小时
	}

	// 创建JWT声明
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(expiresIn).Unix(),
		"iat":      time.Now().Unix(),
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(s.config.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
