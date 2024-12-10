package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func CreateArticle(article *domain.Article) error {
	err := config.DB.Create(article).Error
	if err != nil {
		return err
	}

	return nil
}

func GetArticles() ([]domain.Article, error) {
	var articles []domain.Article
	err := config.DB.Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func GetArticleByID(id string) (domain.Article, error) {
	var article domain.Article
	err := config.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		return article, err
	}
	return article, nil
}

func UpdateArticle(article *domain.Article) error {

	err := config.DB.Save(article).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteArticle(id string) error {

	var article domain.Article
	err := config.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		return err
	}

	err = config.DB.Delete(&article).Error
	if err != nil {
		return err
	}
	return nil
}
