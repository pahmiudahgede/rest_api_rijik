package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type ArticleService struct {
	repo *repositories.ArticleRepository
}

func NewArticleService(repo *repositories.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) GetAllArticles() ([]dto.ArticleResponse, error) {
	ctx := config.Context()
	cacheKey := "articles:all"

	cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var cachedArticles []dto.ArticleResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedArticles); err == nil {
			return cachedArticles, nil
		}
	}

	articles, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []dto.ArticleResponse
	for _, article := range articles {
		result = append(result, dto.ArticleResponse{
			ID:          article.ID,
			Title:       article.Title,
			CoverImage:  article.CoverImage,
			Author:      article.Author,
			Heading:     article.Heading,
			Content:     article.Content,
			PublishedAt: utils.FormatDateToIndonesianFormat(article.PublishedAt),
			UpdatedAt:   utils.FormatDateToIndonesianFormat(article.UpdatedAt),
		})
	}

	cacheData, _ := json.Marshal(result)
	config.RedisClient.Set(ctx, cacheKey, cacheData, time.Minute*5)

	return result, nil
}

func (s *ArticleService) GetArticleByID(id string) (*dto.ArticleResponse, error) {
	ctx := config.Context()
	cacheKey := "articles:" + id

	cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var cachedArticle dto.ArticleResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedArticle); err == nil {
			return &cachedArticle, nil
		}
	}

	article, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	result := &dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: utils.FormatDateToIndonesianFormat(article.PublishedAt),
		UpdatedAt:   utils.FormatDateToIndonesianFormat(article.UpdatedAt),
	}

	cacheData, _ := json.Marshal(result)
	config.RedisClient.Set(ctx, cacheKey, cacheData, time.Minute*5)

	return result, nil
}

func (s *ArticleService) CreateArticle(request *dto.ArticleCreateRequest) (*dto.ArticleResponse, error) {

	if request.Title == "" || request.CoverImage == "" || request.Author == "" ||
		request.Heading == "" || request.Content == "" {
		return nil, errors.New("invalid input data")
	}

	newArticle := &domain.Article{
		Title:      request.Title,
		CoverImage: request.CoverImage,
		Author:     request.Author,
		Heading:    request.Heading,
		Content:    request.Content,
	}

	err := s.repo.Create(newArticle)
	if err != nil {
		return nil, errors.New("failed to create article")
	}

	ctx := config.Context()
	config.RedisClient.Del(ctx, "articles:all")

	response := &dto.ArticleResponse{
		ID:          newArticle.ID,
		Title:       newArticle.Title,
		CoverImage:  newArticle.CoverImage,
		Author:      newArticle.Author,
		Heading:     newArticle.Heading,
		Content:     newArticle.Content,
		PublishedAt: utils.FormatDateToIndonesianFormat(newArticle.PublishedAt),
		UpdatedAt:   utils.FormatDateToIndonesianFormat(newArticle.UpdatedAt),
	}

	return response, nil
}

func (s *ArticleService) UpdateArticle(id string, request *dto.ArticleUpdateRequest) (*dto.ArticleResponse, error) {

	if err := dto.GetValidator().Struct(request); err != nil {
		return nil, errors.New("invalid input data")
	}

	article, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}

	if request.Title != nil {
		article.Title = *request.Title
	}
	if request.CoverImage != nil {
		article.CoverImage = *request.CoverImage
	}
	if request.Author != nil {
		article.Author = *request.Author
	}
	if request.Heading != nil {
		article.Heading = *request.Heading
	}
	if request.Content != nil {
		article.Content = *request.Content
	}
	article.UpdatedAt = time.Now()

	err = s.repo.Update(article)
	if err != nil {
		return nil, errors.New("failed to update article")
	}

	ctx := config.Context()
	config.RedisClient.Del(ctx, "articles:all")
	config.RedisClient.Del(ctx, "articles:"+id)

	response := &dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: utils.FormatDateToIndonesianFormat(article.PublishedAt),
		UpdatedAt:   utils.FormatDateToIndonesianFormat(article.UpdatedAt),
	}

	return response, nil
}

func (s *ArticleService) DeleteArticle(id string) error {

	article, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("article not found")
	}

	err = s.repo.Delete(article)
	if err != nil {
		return errors.New("failed to delete article")
	}

	ctx := config.Context()
	config.RedisClient.Del(ctx, "articles:all")
	config.RedisClient.Del(ctx, "articles:"+id)

	return nil
}