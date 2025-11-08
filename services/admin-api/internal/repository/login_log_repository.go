package repository

import (
	"goweb/services/admin-api/internal/model"
	"time"

	"gorm.io/gorm"
)

// LoginLogRepository 登录记录数据访问层
type LoginLogRepository struct {
	db *gorm.DB
}

// NewLoginLogRepository 创建登录记录数据访问层实例
func NewLoginLogRepository(db *gorm.DB) *LoginLogRepository {
	return &LoginLogRepository{
		db: db,
	}
}

// Create 创建登录记录
func (r *LoginLogRepository) Create(log *model.AdminLoginLog) error {
	return r.db.Create(log).Error
}

// List 获取登录记录列表
func (r *LoginLogRepository) List(req *model.AdminLoginLogListRequest) ([]model.AdminLoginLog, int64, error) {
	var logs []model.AdminLoginLog
	var total int64

	query := r.db.Model(&model.AdminLoginLog{})

	// 搜索条件
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.StartTime != "" {
		query = query.Where("login_time >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("login_time <= ?", req.EndTime)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("login_time DESC").Offset(offset).Limit(req.PageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetRecentLoginUsers 获取最近登录的用户（成功登录的）
func (r *LoginLogRepository) GetRecentLoginUsers(limit int) ([]model.RecentLoginUser, error) {
	var results []struct {
		UserID    uint64
		Username  string
		Email     string
		Name      *string
		Status    int8
		LoginTime time.Time
		LoginIP   *string
	}

	err := r.db.Table("gf_admin_login_logs").
		Select(`
			gf_admin_login_logs.user_id,
			gf_admin_login_logs.username,
			gf_admin_users.email,
			gf_admin_users.name,
			gf_admin_users.status,
			gf_admin_login_logs.login_time,
			gf_admin_login_logs.login_ip
		`).
		Joins("INNER JOIN gf_admin_users ON gf_admin_login_logs.user_id = gf_admin_users.id").
		Where("gf_admin_login_logs.status = ?", 1).
		Where("gf_admin_login_logs.user_id = (SELECT MAX(ll2.user_id) FROM gf_admin_login_logs ll2 WHERE ll2.user_id = gf_admin_login_logs.user_id AND ll2.status = 1)").
		Group("gf_admin_login_logs.user_id, gf_admin_login_logs.username, gf_admin_users.email, gf_admin_users.name, gf_admin_users.status").
		Order("MAX(gf_admin_login_logs.login_time) DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	users := make([]model.RecentLoginUser, 0, len(results))
	for _, result := range results {
		users = append(users, model.RecentLoginUser{
			ID:        result.UserID,
			Username:  result.Username,
			Email:     result.Email,
			Name:      result.Name,
			Status:    result.Status,
			LoginTime: result.LoginTime,
			LoginIP:   result.LoginIP,
		})
	}

	return users, nil
}

// GetRecentLoginUsersV2 获取最近登录的用户（优化版本：使用子查询获取每个用户最新的登录记录）
func (r *LoginLogRepository) GetRecentLoginUsersV2(limit int) ([]model.RecentLoginUser, error) {
	var results []struct {
		UserID    uint64     `gorm:"column:user_id"`
		Username  string     `gorm:"column:username"`
		Email     string     `gorm:"column:email"`
		Name      *string    `gorm:"column:name"`
		Status    int8       `gorm:"column:status"`
		LoginTime time.Time  `gorm:"column:login_time"`
		LoginIP   *string    `gorm:"column:login_ip"`
	}

	// 使用子查询获取每个用户最新的登录记录
	err := r.db.Raw(`
		SELECT 
			ll.user_id,
			ll.username,
			u.email,
			u.name,
			u.status,
			ll.login_time,
			ll.login_ip
		FROM gf_admin_login_logs ll
		INNER JOIN gf_admin_users u ON ll.user_id = u.id
		INNER JOIN (
			SELECT user_id, MAX(login_time) as max_login_time
			FROM gf_admin_login_logs
			WHERE status = 1
			GROUP BY user_id
		) latest ON ll.user_id = latest.user_id AND ll.login_time = latest.max_login_time
		WHERE ll.status = 1
		ORDER BY ll.login_time DESC
		LIMIT ?
	`, limit).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	users := make([]model.RecentLoginUser, 0, len(results))
	for _, result := range results {
		users = append(users, model.RecentLoginUser{
			ID:        result.UserID,
			Username:  result.Username,
			Email:     result.Email,
			Name:      result.Name,
			Status:    result.Status,
			LoginTime: result.LoginTime,
			LoginIP:   result.LoginIP,
		})
	}

	return users, nil
}

