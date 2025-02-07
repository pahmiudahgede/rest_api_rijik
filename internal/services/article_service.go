package services

import (
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type ArticleService interface {
	CreateArticle(articleDTO dto.RequestArticleDTO) (*dto.ArticleResponseDTO, error)
}

type articleService struct {
	ArticleRepo repositories.ArticleRepository
}

func NewArticleService(articleRepo repositories.ArticleRepository) ArticleService {
	return &articleService{ArticleRepo: articleRepo}
}

func (s *articleService) CreateArticle(articleDTO dto.RequestArticleDTO) (*dto.ArticleResponseDTO, error) {

	article := &model.Article{
		Title:      articleDTO.Title,
		CoverImage: articleDTO.CoverImage,
		Author:     articleDTO.Author,
		Heading:    articleDTO.Heading,
		Content:    articleDTO.Content,
	}

	err := s.ArticleRepo.CreateArticle(article)
	if err != nil {
		return nil, fmt.Errorf("failed to create article: %v", err)
	}

	publishedAt, _ := utils.FormatDateToIndonesianFormat(article.PublishedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(article.UpdatedAt)

	articleResponseDTO := &dto.ArticleResponseDTO{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: publishedAt,
		UpdatedAt:   updatedAt,
	}

	cacheKey := fmt.Sprintf("article:%s", article.ID)
	cacheData := map[string]interface{}{
		"data": articleResponseDTO,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching article to Redis: %v\n", err)
	}

	return articleResponseDTO, nil
}
