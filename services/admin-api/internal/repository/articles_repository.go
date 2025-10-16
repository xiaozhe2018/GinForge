package repository

import (
	"gorm.io/gorm"
	"goweb/services/admin-api/internal/model"
)

// ArticlesRepository Articles管理 Repository
type ArticlesRepository struct {
	db *gorm.DB
}

// NewArticlesRepository 创建 Repository 实例
func NewArticlesRepository(db *gorm.DB) *ArticlesRepository {
	return &ArticlesRepository{
		db: db,
	}
}

// Create 创建Articles管理
func (r *ArticlesRepository) Create(articles *model.Articles) error {
	return r.db.Create(articles).Error
}

// GetByID 根据 ID 获取Articles管理
func (r *ArticlesRepository) GetByID(id uint64) (*model.Articles, error) {
	var articles model.Articles
	err := r.db.First(&articles, id).Error
	if err != nil {
		return nil, err
	}
	return &articles, nil
}

// Update 更新Articles管理
func (r *ArticlesRepository) Update(articles *model.Articles) error {
	return r.db.Save(articles).Error
}

// Delete 删除Articles管理
func (r *ArticlesRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Articles{}, id).Error
}

// List 获取Articles管理列表
func (r *ArticlesRepository) List(req *model.ArticlesListRequest) ([]*model.Articles, int64, error) {
	var list []*model.Articles
	var total int64
	
	db := r.db.Model(&model.Articles{})
	// 搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		db = db.Where("id LIKE ? OR title LIKE ? OR content LIKE ? OR summary LIKE ? OR author_name LIKE ? OR cover_image LIKE ? OR tags LIKE ?", keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}
	
	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 排序
	if req.SortBy != "" {
		order := req.SortBy
		if req.SortOrder == "desc" {
			order += " DESC"
		}
		db = db.Order(order)
	} else {
		db = db.Order("id DESC")
	}
	// 分页
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		db = db.Offset(offset).Limit(req.PageSize)
	}
	
	err := db.Find(&list).Error
	return list, total, err
}

// Exists 检查Articles管理是否存在
func (r *ArticlesRepository) Exists(id uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Articles{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Restore 恢复已删除的Articles管理
func (r *ArticlesRepository) Restore(id uint64) error {
	return r.db.Model(&model.Articles{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error
}

// ForceDelete 永久删除Articles管理
func (r *ArticlesRepository) ForceDelete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Articles{}, id).Error
}
