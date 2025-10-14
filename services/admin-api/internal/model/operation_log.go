package model

import (
	"time"
)

// OperationLog 操作日志模型
type OperationLog struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       *uint64   `json:"user_id" gorm:"index"`
	Username     *string   `json:"username"`
	Method       string    `json:"method" gorm:"size:10;not null"`
	Path         string    `json:"path" gorm:"size:255;not null"`
	IP           *string   `json:"ip" gorm:"size:45;index"`
	UserAgent    *string   `json:"user_agent" gorm:"type:text"`
	RequestData  *string   `json:"request_data" gorm:"type:json"`
	ResponseData *string   `json:"response_data" gorm:"type:json"`
	StatusCode   int       `json:"status_code" gorm:"default:200;index"`
	Duration     int       `json:"duration" gorm:"default:0;comment:请求耗时(毫秒)"`
	CreatedAt    time.Time `json:"created_at" gorm:"index"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "admin_operation_logs"
}
