package repository

import (
	"context"

	"gorm.io/gorm"

	"goweb/pkg/base"
	"goweb/services/file-api/internal/model"
)

// FileRepository 文件仓储接口
type FileRepository interface {
	// 基础操作
	Create(ctx context.Context, file *model.FileRecord) error
	Update(ctx context.Context, file *model.FileRecord) error
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (*model.FileRecord, error)
	List(ctx context.Context, page, pageSize int) ([]*model.FileRecord, int64, error)

	// 文件记录操作
	FindByHash(ctx context.Context, hash string) (*model.FileRecord, error)
	FindByPath(ctx context.Context, relativePath string) (*model.FileRecord, error)
	FindByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*model.FileRecord, int64, error)
	FindByType(ctx context.Context, fileType string, page, pageSize int) ([]*model.FileRecord, int64, error)
	IncrementDownloadCount(ctx context.Context, id uint) error

	// 日志操作
	CreateUploadLog(ctx context.Context, log *model.FileUploadLog) error
	CreateDownloadLog(ctx context.Context, log *model.FileDownloadLog) error

	// 统计操作
	GetTotalSize(ctx context.Context) (int64, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetTypeStatistics(ctx context.Context) (map[string]int64, error)
}

// fileRepository 文件仓储实现
type fileRepository struct {
	*base.BaseRepository
}

// NewFileRepository 创建文件仓储
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{
		BaseRepository: base.NewBaseRepository(db, nil),
	}
}

// Create 创建文件记录
func (r *fileRepository) Create(ctx context.Context, file *model.FileRecord) error {
	return r.BaseRepository.Create(ctx, file)
}

// Update 更新文件记录
func (r *fileRepository) Update(ctx context.Context, file *model.FileRecord) error {
	return r.BaseRepository.Update(ctx, file)
}

// Delete 删除文件记录
func (r *fileRepository) Delete(ctx context.Context, id uint) error {
	return r.BaseRepository.Delete(ctx, &model.FileRecord{}, id)
}

// Get 获取文件记录
func (r *fileRepository) Get(ctx context.Context, id uint) (*model.FileRecord, error) {
	var file model.FileRecord
	if err := r.BaseRepository.FindByID(ctx, &file, id); err != nil {
		return nil, err
	}
	return &file, nil
}

// List 获取文件列表
func (r *fileRepository) List(ctx context.Context, page, pageSize int) ([]*model.FileRecord, int64, error) {
	var files []*model.FileRecord
	var total int64

	db := r.WithContext(ctx).Where("status = 1")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("upload_time DESC").Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// FindByHash 根据哈希查找文件
func (r *fileRepository) FindByHash(ctx context.Context, hash string) (*model.FileRecord, error) {
	var file model.FileRecord
	err := r.WithContext(ctx).Where("hash = ? AND status = 1", hash).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// FindByPath 根据路径查找文件
func (r *fileRepository) FindByPath(ctx context.Context, relativePath string) (*model.FileRecord, error) {
	var file model.FileRecord
	err := r.WithContext(ctx).Where("relative_path = ? AND status = 1", relativePath).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// FindByUserID 根据用户ID查找文件
func (r *fileRepository) FindByUserID(ctx context.Context, userID uint, page, pageSize int) ([]*model.FileRecord, int64, error) {
	var files []*model.FileRecord
	var total int64

	db := r.WithContext(ctx).Where("uploaded_by = ? AND status = 1", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("upload_time DESC").Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// FindByType 根据文件类型查找
func (r *fileRepository) FindByType(ctx context.Context, fileType string, page, pageSize int) ([]*model.FileRecord, int64, error) {
	var files []*model.FileRecord
	var total int64

	db := r.WithContext(ctx).Where("file_type = ? AND status = 1", fileType)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("upload_time DESC").Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// IncrementDownloadCount 增加下载次数
func (r *fileRepository) IncrementDownloadCount(ctx context.Context, id uint) error {
	return r.WithContext(ctx).Model(&model.FileRecord{}).
		Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + ?", 1)).Error
}

// CreateUploadLog 创建上传日志
func (r *fileRepository) CreateUploadLog(ctx context.Context, log *model.FileUploadLog) error {
	return r.BaseRepository.Create(ctx, log)
}

// CreateDownloadLog 创建下载日志
func (r *fileRepository) CreateDownloadLog(ctx context.Context, log *model.FileDownloadLog) error {
	return r.BaseRepository.Create(ctx, log)
}

// GetTotalSize 获取总大小
func (r *fileRepository) GetTotalSize(ctx context.Context) (int64, error) {
	var totalSize int64
	err := r.WithContext(ctx).Model(&model.FileRecord{}).
		Where("status = 1").
		Select("COALESCE(SUM(file_size), 0)").
		Scan(&totalSize).Error
	return totalSize, err
}

// GetTotalCount 获取总数量
func (r *fileRepository) GetTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.WithContext(ctx).Model(&model.FileRecord{}).
		Where("status = 1").
		Count(&count).Error
	return count, err
}

// GetTypeStatistics 获取类型统计
func (r *fileRepository) GetTypeStatistics(ctx context.Context) (map[string]int64, error) {
	type Result struct {
		FileType string
		Count    int64
	}

	var results []Result
	err := r.WithContext(ctx).Model(&model.FileRecord{}).
		Select("file_type, COUNT(*) as count").
		Where("status = 1").
		Group("file_type").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	statistics := make(map[string]int64)
	for _, result := range results {
		statistics[result.FileType] = result.Count
	}

	return statistics, nil
}
