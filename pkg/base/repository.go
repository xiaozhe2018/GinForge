package base

import (
	"context"
	"goweb/pkg/logger"
	"goweb/pkg/model"

	"gorm.io/gorm"
)

// BaseRepository 基础仓储类
type BaseRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository(db *gorm.DB, logger logger.Logger) *BaseRepository {
	return &BaseRepository{
		db:     db,
		logger: logger,
	}
}

// SetLogger 设置日志器
func (r *BaseRepository) SetLogger(logger logger.Logger) {
	r.logger = logger
}

// GetDB 获取数据库连接
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}

// WithContext 创建带上下文的数据库连接
func (r *BaseRepository) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// Create 创建记录
func (r *BaseRepository) Create(ctx context.Context, model interface{}) error {
	return r.WithContext(ctx).Create(model).Error
}

// Update 更新记录
func (r *BaseRepository) Update(ctx context.Context, model interface{}) error {
	return r.WithContext(ctx).Save(model).Error
}

// Delete 删除记录
func (r *BaseRepository) Delete(ctx context.Context, model interface{}, id interface{}) error {
	return r.WithContext(ctx).Delete(model, id).Error
}

// FindByID 根据ID查找
func (r *BaseRepository) FindByID(ctx context.Context, model interface{}, id interface{}) error {
	return r.WithContext(ctx).First(model, id).Error
}

// FindOne 查找单条记录
func (r *BaseRepository) FindOne(ctx context.Context, model interface{}, conditions ...interface{}) error {
	return r.WithContext(ctx).Where(conditions[0], conditions[1:]...).First(model).Error
}

// FindAll 查找所有记录
func (r *BaseRepository) FindAll(ctx context.Context, models interface{}, conditions ...interface{}) error {
	query := r.WithContext(ctx)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	return query.Find(models).Error
}

// Count 统计记录数
func (r *BaseRepository) Count(ctx context.Context, model interface{}, conditions ...interface{}) (int64, error) {
	var count int64
	query := r.WithContext(ctx).Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.Count(&count).Error
	return count, err
}

// Paginate 分页查询
func (r *BaseRepository) Paginate(ctx context.Context, models interface{}, pagination *model.Pagination, conditions ...interface{}) (*model.PaginationResult, error) {
	var total int64
	query := r.WithContext(ctx)

	// 统计总数
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	if err := query.Model(models).Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := pagination.Offset()
	limit := pagination.PageSize
	if err := query.Offset(offset).Limit(limit).Find(models).Error; err != nil {
		return nil, err
	}

	return model.NewPaginationResult(models, total, pagination.Page, pagination.PageSize), nil
}

// Transaction 事务执行
func (r *BaseRepository) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.WithContext(ctx).Transaction(fn)
}

// LogInfo 记录信息日志
func (r *BaseRepository) LogInfo(msg string, fields ...interface{}) {
	if r.logger != nil {
		r.logger.Info(msg, fields...)
	}
}

// LogError 记录错误日志
func (r *BaseRepository) LogError(msg string, err error, fields ...interface{}) {
	if r.logger != nil {
		allFields := append([]interface{}{"error", err}, fields...)
		r.logger.Error(msg, allFields...)
	}
}
