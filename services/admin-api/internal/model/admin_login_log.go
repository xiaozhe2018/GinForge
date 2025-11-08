package model

import (
	"time"
)

// AdminLoginLog 登录记录模型
type AdminLoginLog struct {
	ID            uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        uint64    `json:"user_id" gorm:"type:bigint unsigned;not null;index;comment:用户ID"`
	Username      string    `json:"username" gorm:"type:varchar(50);not null;index;comment:用户名"`
	LoginIP       *string   `json:"login_ip" gorm:"type:varchar(45);index;comment:登录IP"`
	UserAgent     *string   `json:"user_agent" gorm:"type:text;comment:用户代理"`
	LoginTime     time.Time `json:"login_time" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;index;comment:登录时间"`
	Status        int8      `json:"status" gorm:"type:tinyint(1);default:1;index;comment:登录状态:1-成功,0-失败"`
	FailureReason *string   `json:"failure_reason" gorm:"type:varchar(255);comment:失败原因"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 返回表名
func (AdminLoginLog) TableName() string {
	return "gf_admin_login_logs"
}

// AdminLoginLogListRequest 登录记录列表请求
type AdminLoginLogListRequest struct {
	Page     int     `form:"page"`                        // 页码
	PageSize int     `form:"page_size" binding:"max=100"` // 每页数量
	UserID   *uint64 `form:"user_id"`                     // 用户ID
	Username string  `form:"username"`                    // 用户名
	Status   *int8   `form:"status"`                      // 登录状态
	StartTime string `form:"start_time"`                  // 开始时间
	EndTime   string `form:"end_time"`                    // 结束时间
}

// AdminLoginLogListResponse 登录记录列表响应
type AdminLoginLogListResponse struct {
	List  []AdminLoginLog `json:"list"`
	Total int64           `json:"total"`
}

// RecentLoginUser 最近登录用户
type RecentLoginUser struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      *string   `json:"name"`
	Status    int8      `json:"status"`
	LoginTime time.Time `json:"login_time"`
	LoginIP   *string   `json:"login_ip"`
}

