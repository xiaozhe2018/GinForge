package repository

import (
	"goweb/services/admin-api/internal/model"

	"gorm.io/gorm"
)

// OperationLogRepository 操作日志数据访问层
type OperationLogRepository struct {
	db *gorm.DB
}

// NewOperationLogRepository 创建操作日志数据访问层实例
func NewOperationLogRepository(database *gorm.DB) *OperationLogRepository {
	return &OperationLogRepository{
		db: database,
	}
}

// Create 创建操作日志
func (r *OperationLogRepository) Create(log *model.OperationLog) error {
	return r.db.Create(log).Error
}

// List 获取操作日志列表
func (r *OperationLogRepository) List(page, pageSize int, userID *uint64) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64

	query := r.db.Model(&model.OperationLog{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// DeleteOldLogs 删除旧的操作日志（保留最近N天）
func (r *OperationLogRepository) DeleteOldLogs(days int) error {
	return r.db.Where("created_at < DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Delete(&model.OperationLog{}).Error
}
