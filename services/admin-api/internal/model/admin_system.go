package model

import (
	"time"
)

// AdminSystemConfig 系统配置模型
type AdminSystemConfig struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Key         string    `json:"key" gorm:"type:varchar(100);uniqueIndex;not null;comment:配置键"`
	Value       *string   `json:"value" gorm:"type:text;comment:配置值"`
	Type        string    `json:"type" gorm:"type:varchar(20);default:string;comment:配置类型:string,number,boolean,json"`
	Description *string   `json:"description" gorm:"type:varchar(255);comment:配置描述"`
	Group       string    `json:"group" gorm:"type:varchar(50);default:default;index;comment:配置分组"`
	Sort        int       `json:"sort" gorm:"type:int(11);default:0;index;comment:排序"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 返回表名
func (AdminSystemConfig) TableName() string {
	return "gf_admin_system_configs"
}

// AdminOperationLog 操作日志模型
type AdminOperationLog struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       *uint64   `json:"user_id" gorm:"type:bigint(20) unsigned;index;comment:操作用户ID"`
	Username     *string   `json:"username" gorm:"type:varchar(50);comment:操作用户名"`
	Method       string    `json:"method" gorm:"type:varchar(10);not null;index;comment:请求方法"`
	Path         string    `json:"path" gorm:"type:varchar(255);not null;index;comment:请求路径"`
	IP           *string   `json:"ip" gorm:"type:varchar(45);index;comment:请求IP"`
	UserAgent    *string   `json:"user_agent" gorm:"type:text;comment:用户代理"`
	RequestData  *string   `json:"request_data" gorm:"type:json;comment:请求数据"`
	ResponseData *string   `json:"response_data" gorm:"type:json;comment:响应数据"`
	StatusCode   int       `json:"status_code" gorm:"type:int(11);default:200;index;comment:响应状态码"`
	Duration     int       `json:"duration" gorm:"type:int(11);default:0;comment:请求耗时(毫秒)"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime;index"`
}

// TableName 返回表名
func (AdminOperationLog) TableName() string {
	return "gf_admin_operation_logs"
}

// AdminSystemInfo 系统信息响应
type AdminSystemInfo struct {
	OnlineUsers int    `json:"online_users"`
	CPUUsage    int    `json:"cpu_usage"`
	MemoryUsage int    `json:"memory_usage"`
	DiskUsage   int    `json:"disk_usage"`
	NetworkIn   int64  `json:"network_in"`
	NetworkOut  int64  `json:"network_out"`
	Uptime      string `json:"uptime"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

// AdminSystemConfigUpdateRequest 系统配置更新请求
type AdminSystemConfigUpdateRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// AdminSystemConfigListRequest 系统配置列表请求
type AdminSystemConfigListRequest struct {
	Page     int    `form:"page"`                        // 页码，不传时使用默认值1
	PageSize int    `form:"page_size" binding:"max=100"` // 每页数量，不传时使用默认值10
	Group    string `form:"group"`                       // 配置分组
	Keyword  string `form:"keyword"`                     // 搜索关键词
}

// AdminSystemConfigListResponse 系统配置列表响应
type AdminSystemConfigListResponse struct {
	List  []AdminSystemConfig `json:"list"`
	Total int64               `json:"total"`
}

// AdminOperationLogListRequest 操作日志列表请求
type AdminOperationLogListRequest struct {
	Page       int     `form:"page"`                        // 页码，不传时使用默认值1
	PageSize   int     `form:"page_size" binding:"max=100"` // 每页数量，不传时使用默认值10
	UserID     *uint64 `form:"user_id"`                     // 用户ID
	Username   string  `form:"username"`                    // 用户名
	Method     string  `form:"method"`                      // 请求方法
	Path       string  `form:"path"`                        // 请求路径
	IP         string  `form:"ip"`                          // IP地址
	StatusCode *int    `form:"status_code"`                 // 状态码
	StartTime  string  `form:"start_time"`                  // 开始时间
	EndTime    string  `form:"end_time"`                    // 结束时间
}

// AdminOperationLogListResponse 操作日志列表响应
type AdminOperationLogListResponse struct {
	List  []AdminOperationLog `json:"list"`
	Total int64               `json:"total"`
}
