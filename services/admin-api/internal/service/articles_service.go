package service

import (
	"errors"
	"gorm.io/gorm"
	
	"goweb/pkg/logger"
	"goweb/services/admin-api/internal/model"
	"goweb/services/admin-api/internal/repository"
)

// ArticlesService Articles管理 Service
type ArticlesService struct {
	repo   *repository.ArticlesRepository
	logger logger.Logger
}

// NewArticlesService 创建 Service 实例
func NewArticlesService(repo *repository.ArticlesRepository, logger logger.Logger) *ArticlesService {
	return &ArticlesService{
		repo:   repo,
		logger: logger,
	}
}

// Create 创建Articles管理
func (s *ArticlesService) Create(req *model.ArticlesCreateRequest) (*model.Articles, error) {
	articles := &model.Articles{
		Title: req.Title,
		Slug: req.Slug,
		AuthorId: req.AuthorId,
		AuthorName: req.AuthorName,
		CategoryId: req.CategoryId,
		Summary: req.Summary,
		Content: req.Content,
		CoverImage: req.CoverImage,
		ViewCount: req.ViewCount,
		LikeCount: req.LikeCount,
		CommentCount: req.CommentCount,
		IsPublished: req.IsPublished,
		IsTop: req.IsTop,
		IsFeatured: req.IsFeatured,
		PublishedAt: req.PublishedAt,
		Tags: req.Tags,
		SeoTitle: req.SeoTitle,
		SeoKeywords: req.SeoKeywords,
		SeoDescription: req.SeoDescription,
		Status: req.Status,
	}
	
	if err := s.repo.Create(articles); err != nil {
		s.logger.Error("创建Articles管理失败", err)
		return nil, errors.New("创建Articles管理失败")
	}
	
	return articles, nil
}

// GetByID 根据 ID 获取Articles管理
func (s *ArticlesService) GetByID(id uint64) (*model.Articles, error) {
	articles, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Articles管理不存在")
		}
		s.logger.Error("获取Articles管理失败", err, "id", id)
		return nil, errors.New("获取Articles管理失败")
	}
	
	return articles, nil
}

// Update 更新Articles管理
func (s *ArticlesService) Update(id uint64, req *model.ArticlesUpdateRequest) error {
	// 检查Articles管理是否存在
	articles, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Articles管理不存在")
		}
		return errors.New("获取Articles管理失败")
	}
	
	// 更新字段
	articles.Title = req.Title
	articles.Slug = req.Slug
	articles.AuthorId = req.AuthorId
	articles.AuthorName = req.AuthorName
	articles.CategoryId = req.CategoryId
	articles.Summary = req.Summary
	articles.Content = req.Content
	articles.CoverImage = req.CoverImage
	articles.ViewCount = req.ViewCount
	articles.LikeCount = req.LikeCount
	articles.CommentCount = req.CommentCount
	articles.IsPublished = req.IsPublished
	articles.IsTop = req.IsTop
	articles.IsFeatured = req.IsFeatured
	articles.PublishedAt = req.PublishedAt
	articles.Tags = req.Tags
	articles.SeoTitle = req.SeoTitle
	articles.SeoKeywords = req.SeoKeywords
	articles.SeoDescription = req.SeoDescription
	articles.Status = req.Status
	
	if err := s.repo.Update(articles); err != nil {
		s.logger.Error("更新Articles管理失败", err, "id", id)
		return errors.New("更新Articles管理失败")
	}
	
	return nil
}

// Delete 删除Articles管理
func (s *ArticlesService) Delete(id uint64) error {
	// 检查Articles管理是否存在
	exists, err := s.repo.Exists(id)
	if err != nil {
		s.logger.Error("检查Articles管理是否存在失败", err, "id", id)
		return errors.New("检查Articles管理是否存在失败")
	}
	
	if !exists {
		return errors.New("Articles管理不存在")
	}
	
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("删除Articles管理失败", err, "id", id)
		return errors.New("删除Articles管理失败")
	}
	
	return nil
}

// List 获取Articles管理列表
func (s *ArticlesService) List(req *model.ArticlesListRequest) ([]*model.Articles, int64, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	
	list, total, err := s.repo.List(req)
	if err != nil {
		s.logger.Error("获取Articles管理列表失败", err)
		return nil, 0, errors.New("获取Articles管理列表失败")
	}
	
	return list, total, nil
}
