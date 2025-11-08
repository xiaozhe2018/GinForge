package service

import (
	"context"
	"encoding/json"
	"fmt"
	"goweb/pkg/logger"
	"goweb/pkg/notification"
	"goweb/pkg/redis"
	"goweb/pkg/validator"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/repository"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"gorm.io/gorm"
)

// AdminSystemService 系统管理服务
type AdminSystemService struct {
	db            *gorm.DB
	redisClient   *redis.Client
	notifyService *notification.Service
	logger        logger.Logger
	loginLogRepo  *repository.LoginLogRepository
	startTime     time.Time
	version       string
	environment   string
}

// NewAdminSystemService 创建系统管理服务
func NewAdminSystemService(db *gorm.DB, redisClient *redis.Client, notifyService *notification.Service, log logger.Logger) *AdminSystemService {
	// 获取环境变量
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	return &AdminSystemService{
		db:            db,
		redisClient:   redisClient,
		notifyService: notifyService,
		logger:        log,
		loginLogRepo:  repository.NewLoginLogRepository(db),
		startTime:     time.Now(),
		version:       "1.0.0", // 可从配置或构建信息中获取
		environment:   env,
	}
}

// GetOnlineUserCount 获取在线用户数
func (s *AdminSystemService) GetOnlineUserCount() int {
	// 从Redis获取在线用户数
	if s.redisClient != nil && s.redisClient.IsEnabled() {
		count, err := s.redisClient.SCard(context.Background(), "online_users").Result()
		if err == nil {
			return int(count)
		}
	}
	return 0
}

// GetCPUUsage 获取CPU使用率
func (s *AdminSystemService) GetCPUUsage() int {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil || len(percent) == 0 {
		s.logger.Error("Failed to get CPU usage", "error", err)
		return 0
	}
	return int(percent[0])
}

// GetMemoryUsage 获取内存使用率
func (s *AdminSystemService) GetMemoryUsage() int {
	v, err := mem.VirtualMemory()
	if err != nil {
		s.logger.Error("Failed to get memory usage", "error", err)
		return 0
	}
	return int(v.UsedPercent)
}

// GetDiskUsage 获取磁盘使用率
func (s *AdminSystemService) GetDiskUsage() int {
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:"
	}

	usage, err := disk.Usage(path)
	if err != nil {
		s.logger.Error("Failed to get disk usage", "error", err)
		return 0
	}
	return int(usage.UsedPercent)
}

// GetNetworkIn 获取网络入流量
func (s *AdminSystemService) GetNetworkIn() int64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		s.logger.Error("Failed to get network stats", "error", err)
		return 0
	}
	return int64(stats[0].BytesRecv)
}

// GetNetworkOut 获取网络出流量
func (s *AdminSystemService) GetNetworkOut() int64 {
	stats, err := net.IOCounters(false)
	if err != nil || len(stats) == 0 {
		s.logger.Error("Failed to get network stats", "error", err)
		return 0
	}
	return int64(stats[0].BytesSent)
}

// GetUptime 获取系统运行时间
func (s *AdminSystemService) GetUptime() string {
	uptime := time.Since(s.startTime)
	days := int(uptime.Hours() / 24)
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟", days, hours, minutes)
	}
	return fmt.Sprintf("%d小时%d分钟", hours, minutes)
}

// GetVersion 获取系统版本
func (s *AdminSystemService) GetVersion() string {
	return s.version
}

// GetEnvironment 获取系统环境
func (s *AdminSystemService) GetEnvironment() string {
	return s.environment
}

// GetConfigList 获取配置列表
func (s *AdminSystemService) GetConfigList(ctx context.Context, req model.AdminSystemConfigListRequest) ([]model.AdminSystemConfig, int64, error) {
	var configs []model.AdminSystemConfig
	var total int64

	query := s.db.Model(&model.AdminSystemConfig{})

	// 应用过滤条件
	if req.Group != "" {
		query = query.Where("`group` = ?", req.Group)
	}
	if req.Keyword != "" {
		query = query.Where("`key` LIKE ? OR description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("`group`, sort, id").Offset(offset).Limit(req.PageSize).Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

// GetConfig 获取单个配置
func (s *AdminSystemService) GetConfig(ctx context.Context, key string) (model.AdminSystemConfig, error) {
	var config model.AdminSystemConfig
	if err := s.db.Where("`key` = ?", key).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return config, nil
		}
		return config, err
	}
	return config, nil
}

// GetConfigValue 获取配置值
func (s *AdminSystemService) GetConfigValue(ctx context.Context, key string) (string, error) {
	var config model.AdminSystemConfig
	if err := s.db.Where("`key` = ?", key).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	if config.Value == nil {
		return "", nil
	}
	return *config.Value, nil
}

// UpdateConfig 更新配置
func (s *AdminSystemService) UpdateConfig(ctx context.Context, key string, value string) error {
	var config model.AdminSystemConfig
	result := s.db.Where("`key` = ?", key).First(&config)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 如果配置不存在，创建新配置
			description := "系统配置"
			config = model.AdminSystemConfig{
				Key:         key,
				Value:       &value,
				Type:        "string",
				Description: &description,
				Group:       "default",
				Sort:        0,
			}
			return s.db.Create(&config).Error
		}
		return result.Error
	}

	// 更新配置值
	config.Value = &value
	return s.db.Save(&config).Error
}

// SendTestEmail 发送测试邮件
func (s *AdminSystemService) SendTestEmail(ctx context.Context, to string) error {
	if s.notifyService == nil {
		return fmt.Errorf("notification service not configured")
	}

	// 创建测试邮件请求
	req := &notification.NotificationRequest{
		Type:       notification.TypeEmail,
		Recipients: []string{to},
		Subject:    "GinForge 测试邮件",
		Template:   "notification",
		Data: map[string]interface{}{
			"Title":   "邮件服务测试",
			"Name":    "管理员",
			"Message": "这是一封测试邮件，如果您收到这封邮件，说明您的邮件服务配置正确。",
			"Level":   "info",
			"Year":    time.Now().Format("2006"),
		},
	}

	// 发送邮件
	_, err := s.notifyService.Send(ctx, req)
	return err
}

// TestCacheConnection 测试缓存连接
func (s *AdminSystemService) TestCacheConnection(ctx context.Context) error {
	if s.redisClient == nil || !s.redisClient.IsEnabled() {
		return fmt.Errorf("redis client not configured")
	}

	// 测试Redis连接
	_, err := s.redisClient.Ping(ctx).Result()
	return err
}

// ClearCache 清空缓存
func (s *AdminSystemService) ClearCache(ctx context.Context) error {
	if s.redisClient == nil || !s.redisClient.IsEnabled() {
		return fmt.Errorf("redis client not configured")
	}

	// 清空所有缓存
	return s.redisClient.FlushDB(ctx).Err()
}

// GetLogList 获取日志列表
func (s *AdminSystemService) GetLogList(ctx context.Context, req model.AdminOperationLogListRequest) ([]model.AdminOperationLog, int64, error) {
	var logs []model.AdminOperationLog
	var total int64

	query := s.db.Model(&model.AdminOperationLog{})

	// 应用过滤条件
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Method != "" {
		query = query.Where("method = ?", req.Method)
	}
	if req.Path != "" {
		query = query.Where("path LIKE ?", "%"+req.Path+"%")
	}
	if req.IP != "" {
		query = query.Where("ip = ?", req.IP)
	}
	if req.StatusCode != nil {
		query = query.Where("status_code = ?", *req.StatusCode)
	}
	if req.StartTime != "" {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("created_at <= ?", req.EndTime)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// ClearLogs 清空日志
func (s *AdminSystemService) ClearLogs(ctx context.Context) error {
	return s.db.Exec("DELETE FROM admin_operation_logs").Error
}

// CheckHealth 健康检查
func (s *AdminSystemService) CheckHealth(ctx context.Context) map[string]interface{} {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"services":  make(map[string]interface{}),
	}

	services := health["services"].(map[string]interface{})

	// 检查数据库
	dbStatus := "healthy"
	dbMessage := "Database connection is healthy"
	sqlDB, err := s.db.DB()
	if err != nil {
		dbStatus = "unhealthy"
		dbMessage = "Failed to get database connection: " + err.Error()
		health["status"] = "unhealthy"
	} else if err := sqlDB.Ping(); err != nil {
		dbStatus = "unhealthy"
		dbMessage = "Database ping failed: " + err.Error()
		health["status"] = "unhealthy"
	}
	services["database"] = map[string]string{
		"status":  dbStatus,
		"message": dbMessage,
	}

	// 检查Redis
	redisStatus := "healthy"
	redisMessage := "Redis connection is healthy"
	if s.redisClient == nil || !s.redisClient.IsEnabled() {
		redisStatus = "unknown"
		redisMessage = "Redis client not configured"
	} else if _, err := s.redisClient.Ping(ctx).Result(); err != nil {
		redisStatus = "unhealthy"
		redisMessage = "Redis ping failed: " + err.Error()
		health["status"] = "unhealthy"
	}
	services["redis"] = map[string]string{
		"status":  redisStatus,
		"message": redisMessage,
	}

	// 检查磁盘空间
	diskStatus := "healthy"
	diskMessage := "Disk space is sufficient"
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:"
	}
	usage, err := disk.Usage(path)
	if err != nil {
		diskStatus = "unknown"
		diskMessage = "Failed to check disk usage: " + err.Error()
	} else if usage.UsedPercent > 90 {
		diskStatus = "warning"
		diskMessage = fmt.Sprintf("Disk usage is high: %.2f%%", usage.UsedPercent)
	}
	services["disk"] = map[string]string{
		"status":  diskStatus,
		"message": diskMessage,
	}

	return health
}

// GetRecentLoginUsers 获取最近登录的用户
func (s *AdminSystemService) GetRecentLoginUsers(ctx context.Context, limit int) ([]model.RecentLoginUser, error) {
	if s.loginLogRepo == nil {
		return []model.RecentLoginUser{}, nil
	}
	return s.loginLogRepo.GetRecentLoginUsersV2(limit)
}

// GetPasswordMinLength 获取密码最小长度配置
func (s *AdminSystemService) GetPasswordMinLength(ctx context.Context) int {
	value, err := s.GetConfigValue(ctx, "security.min_password_length")
	if err != nil || value == "" {
		return 8 // 默认值
	}
	
	minLength, err := strconv.Atoi(value)
	if err != nil || minLength < 6 {
		return 8 // 默认值
	}
	
	return minLength
}

// GetPasswordComplexity 获取密码复杂度配置
func (s *AdminSystemService) GetPasswordComplexity(ctx context.Context) validator.PasswordComplexity {
	value, err := s.GetConfigValue(ctx, "security.password_complexity")
	if err != nil || value == "" {
		// 默认复杂度要求
		return validator.PasswordComplexity{
			RequireLowercase: true,
			RequireNumbers:   true,
		}
	}
	
	// 解析JSON数组
	var requirements []string
	if err := json.Unmarshal([]byte(value), &requirements); err != nil {
		// 解析失败，返回默认值
		return validator.PasswordComplexity{
			RequireLowercase: true,
			RequireNumbers:   true,
		}
	}
	
	return validator.ParseComplexity(requirements)
}

// ValidatePassword 验证密码是否符合安全策略
func (s *AdminSystemService) ValidatePassword(ctx context.Context, password string) error {
	minLength := s.GetPasswordMinLength(ctx)
	complexity := s.GetPasswordComplexity(ctx)
	
	return validator.ValidatePassword(password, minLength, complexity)
}

// GetMaxLoginAttempts 获取最大登录尝试次数
func (s *AdminSystemService) GetMaxLoginAttempts(ctx context.Context) int {
	value, err := s.GetConfigValue(ctx, "security.max_login_attempts")
	if err != nil || value == "" {
		return 5 // 默认值
	}
	
	attempts, err := strconv.Atoi(value)
	if err != nil || attempts < 1 {
		return 5
	}
	
	return attempts
}

// GetLockoutDuration 获取账户锁定时长（分钟）
func (s *AdminSystemService) GetLockoutDuration(ctx context.Context) int {
	value, err := s.GetConfigValue(ctx, "security.lockout_duration")
	if err != nil || value == "" {
		return 15 // 默认15分钟
	}
	
	duration, err := strconv.Atoi(value)
	if err != nil || duration < 1 {
		return 15
	}
	
	return duration
}

// GetSessionTimeout 获取会话超时时间（分钟）
func (s *AdminSystemService) GetSessionTimeout(ctx context.Context) int {
	value, err := s.GetConfigValue(ctx, "security.session_timeout")
	if err != nil || value == "" {
		return 120 // 默认2小时
	}
	
	timeout, err := strconv.Atoi(value)
	if err != nil || timeout < 1 {
		return 120
	}
	
	return timeout
}

// GetSystemBasicInfo 获取系统基本信息
func (s *AdminSystemService) GetSystemBasicInfo(ctx context.Context) map[string]string {
	info := make(map[string]string)
	
	// 读取基本配置
	configs := []string{
		"system.name",
		"system.version",
		"system.description",
		"system.logo",
		"system.default_language",
	}
	
	for _, key := range configs {
		value, err := s.GetConfigValue(ctx, key)
		if err == nil && value != "" {
			info[key] = value
		}
	}
	
	// 设置默认值
	if info["system.name"] == "" {
		info["system.name"] = "GinForge 管理后台"
	}
	if info["system.version"] == "" {
		info["system.version"] = s.version
	}
	if info["system.description"] == "" {
		info["system.description"] = "基于 Go + Gin 的企业级微服务开发框架"
	}
	if info["system.logo"] == "" {
		info["system.logo"] = "/logo.svg"
	}
	if info["system.default_language"] == "" {
		info["system.default_language"] = "zh-CN"
	}
	
	return info
}