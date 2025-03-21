package repositories

import (
	"fmt"

	"rijig/model"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	CreateArticle(article *model.Article) error
	FindArticleByID(id string) (*model.Article, error)
	FindAllArticles(page, limit int) ([]model.Article, int, error)
	UpdateArticle(id string, article *model.Article) error
	DeleteArticle(id string) error
}

type articleRepository struct {
	DB *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{DB: db}
}

func (r *articleRepository) CreateArticle(article *model.Article) error {
	return r.DB.Create(article).Error
}

func (r *articleRepository) FindArticleByID(id string) (*model.Article, error) {
	var article model.Article
	err := r.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("article with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch article: %v", err)
	}
	return &article, nil
}

func (r *articleRepository) FindAllArticles(page, limit int) ([]model.Article, int, error) {
	var articles []model.Article
	var total int64

	if err := r.DB.Model(&model.Article{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count articles: %v", err)
	}

	fmt.Printf("Total Articles Count: %d\n", total)

	if page > 0 && limit > 0 {
		err := r.DB.Offset((page - 1) * limit).Limit(limit).Find(&articles).Error
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch articles: %v", err)
		}
	} else {
		err := r.DB.Find(&articles).Error
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch articles: %v", err)
		}
	}

	return articles, int(total), nil
}

func (r *articleRepository) UpdateArticle(id string, article *model.Article) error {
	return r.DB.Model(&model.Article{}).Where("id = ?", id).Updates(article).Error
}

func (r *articleRepository) DeleteArticle(id string) error {
	result := r.DB.Delete(&model.Article{}, "id = ?", id)
	return result.Error
}
