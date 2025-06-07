package article

import (
	"context"
	"errors"
	"fmt"

	"rijig/model"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article *model.Article) error
	FindArticleByID(ctx context.Context, id string) (*model.Article, error)
	FindAllArticles(ctx context.Context, page, limit int) ([]model.Article, int64, error)
	UpdateArticle(ctx context.Context, id string, article *model.Article) error
	DeleteArticle(ctx context.Context, id string) error
	ArticleExists(ctx context.Context, id string) (bool, error)
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) CreateArticle(ctx context.Context, article *model.Article) error {
	if article == nil {
		return errors.New("article cannot be nil")
	}

	if err := r.db.WithContext(ctx).Create(article).Error; err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}

func (r *articleRepository) FindArticleByID(ctx context.Context, id string) (*model.Article, error) {
	if id == "" {
		return nil, errors.New("article ID cannot be empty")
	}

	var article model.Article
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("article with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch article: %w", err)
	}
	return &article, nil
}

func (r *articleRepository) FindAllArticles(ctx context.Context, page, limit int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	if page < 0 || limit < 0 {
		return nil, 0, errors.New("page and limit must be non-negative")
	}

	if err := r.db.WithContext(ctx).Model(&model.Article{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	query := r.db.WithContext(ctx).Model(&model.Article{})

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch articles: %w", err)
	}

	return articles, total, nil
}

func (r *articleRepository) UpdateArticle(ctx context.Context, id string, article *model.Article) error {
	if id == "" {
		return errors.New("article ID cannot be empty")
	}
	if article == nil {
		return errors.New("article cannot be nil")
	}

	exists, err := r.ArticleExists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check article existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("article with ID %s not found", id)
	}

	result := r.db.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).Updates(article)
	if result.Error != nil {
		return fmt.Errorf("failed to update article: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected when updating article with ID %s", id)
	}

	return nil
}

func (r *articleRepository) DeleteArticle(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("article ID cannot be empty")
	}

	exists, err := r.ArticleExists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check article existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("article with ID %s not found", id)
	}

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Article{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete article: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected when deleting article with ID %s", id)
	}

	return nil
}

func (r *articleRepository) ArticleExists(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, errors.New("article ID cannot be empty")
	}

	var count int64
	err := r.db.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check article existence: %w", err)
	}

	return count > 0, nil
}
