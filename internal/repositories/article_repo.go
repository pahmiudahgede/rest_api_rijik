package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	CreateArticle(article *model.Article) error
	FindArticleByID(id string) (*model.Article, error)
	FindAllArticles(page, limit int) ([]model.Article, int, error)
}

type articleRepository struct {
	DB *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{DB: db}
}

func (r *articleRepository) CreateArticle(article *model.Article) error {
	err := r.DB.Create(article).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepository) FindArticleByID(id string) (*model.Article, error) {
	var article model.Article
	err := r.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) FindAllArticles(page, limit int) ([]model.Article, int, error) {
	var articles []model.Article
	var total int64

	err := r.DB.Model(&model.Article{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		err := r.DB.Offset((page - 1) * limit).Limit(limit).Find(&articles).Error
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.DB.Find(&articles).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return articles, int(total), nil
}
