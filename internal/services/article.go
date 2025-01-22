package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

var ctx = context.Background()

func CreateArticle(articleRequest *dto.ArticleRequest) (dto.ArticleResponse, error) {
	article := domain.Article{
		Title:      articleRequest.Title,
		CoverImage: articleRequest.CoverImage,
		Author:     articleRequest.Author,
		Heading:    articleRequest.Heading,
		Content:    articleRequest.Content,
	}

	err := repositories.CreateArticle(&article)
	if err != nil {
		return dto.ArticleResponse{}, err
	}

	config.RedisClient.Del(ctx, "articles")

	articleResponse := dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: article.PublishedAt,
		UpdatedAt:   article.UpdatedAt,
	}

	return articleResponse, nil
}

func GetArticles() ([]dto.ArticleResponse, error) {
	var response []dto.ArticleResponse

	cachedArticles, err := config.RedisClient.Get(ctx, "articles").Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedArticles), &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	articles, err := repositories.GetArticles()
	if err != nil {
		return nil, err
	}

	for _, article := range articles {
		response = append(response, dto.ArticleResponse{
			ID:          article.ID,
			Title:       article.Title,
			CoverImage:  article.CoverImage,
			Author:      article.Author,
			Heading:     article.Heading,
			Content:     article.Content,
			PublishedAt: article.PublishedAt,
			UpdatedAt:   article.UpdatedAt,
		})
	}

	articlesJSON, _ := json.Marshal(response)
	config.RedisClient.Set(ctx, "articles", articlesJSON, 10*time.Minute)

	return response, nil
}

func GetArticleByID(id string) (dto.ArticleResponse, error) {
	cachedArticle, err := config.RedisClient.Get(ctx, "article:"+id).Result()
	if err == nil {
		var article dto.ArticleResponse
		err := json.Unmarshal([]byte(cachedArticle), &article)
		if err != nil {
			return article, err
		}
		return article, nil
	}

	article, err := repositories.GetArticleByID(id)
	if err != nil {
		return dto.ArticleResponse{}, err
	}

	articleResponse := dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: article.PublishedAt,
		UpdatedAt:   article.UpdatedAt,
	}

	articleJSON, _ := json.Marshal(articleResponse)
	config.RedisClient.Set(ctx, "article:"+id, articleJSON, 10*time.Minute)

	return articleResponse, nil
}

func UpdateArticle(id string, articleUpdateRequest *dto.ArticleUpdateRequest) (*dto.ArticleResponse, error) {

	article, err := repositories.GetArticleByID(id)
	if err != nil {
		return nil, err
	}

	article.Title = articleUpdateRequest.Title
	article.CoverImage = articleUpdateRequest.CoverImage
	article.Author = articleUpdateRequest.Author
	article.Heading = articleUpdateRequest.Heading
	article.Content = articleUpdateRequest.Content
	article.UpdatedAt = time.Now()

	err = repositories.UpdateArticle(&article)
	if err != nil {
		return nil, err
	}

	config.RedisClient.Del(ctx, "article:"+id)
	config.RedisClient.Del(ctx, "articles")

	updatedArticleResponse := dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: article.PublishedAt,
		UpdatedAt:   article.UpdatedAt,
	}

	return &updatedArticleResponse, nil
}

func DeleteArticle(id string) error {

	err := repositories.DeleteArticle(id)
	if err != nil {
		return err
	}

	config.RedisClient.Del(ctx, "article:"+id)
	config.RedisClient.Del(ctx, "articles")

	return nil
}