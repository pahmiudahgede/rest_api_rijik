package repositories

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

type ArticleRepository struct{}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{}
}

func (r *ArticleRepository) GetAll() ([]domain.Article, error) {
	var articles []domain.Article
	err := config.DB.Find(&articles).Error
	if err != nil {
		return nil, errors.New("failed to fetch articles from database")
	}
	return articles, nil
}

func (r *ArticleRepository) GetByID(id string) (*domain.Article, error) {
	var article domain.Article
	err := config.DB.First(&article, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("article not found")
	}
	return &article, nil
}

func (r *ArticleRepository) Create(article *domain.Article) error {
	return config.DB.Create(article).Error
}

func (r *ArticleRepository) Update(article *domain.Article) error {
	return config.DB.Save(article).Error
}

func (r *ArticleRepository) Delete(article *domain.Article) error {
	return config.DB.Delete(article).Error
}