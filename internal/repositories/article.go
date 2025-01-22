package repositories

import (
	"encoding/json"
	"time"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
	"golang.org/x/net/context"
)

var ctx = context.Background()

func CreateArticle(article *domain.Article) error {
	err := config.DB.Create(article).Error
	if err != nil {
		return err
	}

	config.RedisClient.Del(ctx, "articles")
	return nil
}

func GetArticles() ([]domain.Article, error) {
	var articles []domain.Article

	cachedArticles, err := config.RedisClient.Get(ctx, "articles").Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedArticles), &articles)
		if err != nil {
			return nil, err
		}
		return articles, nil
	}

	err = config.DB.Find(&articles).Error
	if err != nil {
		return nil, err
	}

	articlesJSON, _ := json.Marshal(articles)
	config.RedisClient.Set(ctx, "articles", articlesJSON, 10*time.Minute)

	return articles, nil
}

func GetArticleByID(id string) (domain.Article, error) {
	var article domain.Article

	cachedArticle, err := config.RedisClient.Get(ctx, "article:"+id).Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedArticle), &article)
		if err != nil {
			return article, err
		}
		return article, nil
	}

	err = config.DB.Where("id = ?", id).First(&article).Error
	if err != nil {
		return article, err
	}

	articleJSON, _ := json.Marshal(article)
	config.RedisClient.Set(ctx, "article:"+id, articleJSON, 10*time.Minute)

	return article, nil
}

func UpdateArticle(article *domain.Article) error {
	err := config.DB.Save(article).Error
	if err != nil {
		return err
	}

	config.RedisClient.Del(ctx, "article:"+article.ID)
	config.RedisClient.Del(ctx, "articles")

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

	config.RedisClient.Del(ctx, "article:"+id)
	config.RedisClient.Del(ctx, "articles")

	return nil
}