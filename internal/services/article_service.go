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
	GetAllArticles(page, limit int) ([]dto.ArticleResponseDTO, int, error)
	GetArticleByID(id string) (*dto.ArticleResponseDTO, error)
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

func (s *articleService) GetAllArticles(page, limit int) ([]dto.ArticleResponseDTO, int, error) {

	cacheKey := fmt.Sprintf("articles_page:%d_limit:%d", page, limit)

	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var articles []dto.ArticleResponseDTO
		if data, ok := cachedData["data"].([]interface{}); ok {
			for _, item := range data {
				articleData, ok := item.(map[string]interface{})
				if ok {
					articles = append(articles, dto.ArticleResponseDTO{
						ID:          articleData["article_id"].(string),
						Title:       articleData["title"].(string),
						CoverImage:  articleData["coverImage"].(string),
						Author:      articleData["author"].(string),
						Heading:     articleData["heading"].(string),
						Content:     articleData["content"].(string),
						PublishedAt: articleData["publishedAt"].(string),
						UpdatedAt:   articleData["updatedAt"].(string),
					})
				}
			}
			return articles, len(articles), nil
		}
	}

	articles, total, err := s.ArticleRepo.FindAllArticles(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch articles: %v", err)
	}

	var articleDTOs []dto.ArticleResponseDTO
	for _, article := range articles {
		publishedAt, _ := utils.FormatDateToIndonesianFormat(article.PublishedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(article.UpdatedAt)

		articleDTOs = append(articleDTOs, dto.ArticleResponseDTO{
			ID:          article.ID,
			Title:       article.Title,
			CoverImage:  article.CoverImage,
			Author:      article.Author,
			Heading:     article.Heading,
			Content:     article.Content,
			PublishedAt: publishedAt,
			UpdatedAt:   updatedAt,
		})
	}

	cacheData := map[string]interface{}{
		"data": articleDTOs,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching articles to Redis: %v\n", err)
	}

	return articleDTOs, total, nil
}

func (s *articleService) GetArticleByID(id string) (*dto.ArticleResponseDTO, error) {
	cacheKey := fmt.Sprintf("article:%s", id)

	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		articleData, ok := cachedData["data"].(map[string]interface{})
		if ok {
			article := dto.ArticleResponseDTO{
				ID:          articleData["article_id"].(string),
				Title:       articleData["title"].(string),
				CoverImage:  articleData["coverImage"].(string),
				Author:      articleData["author"].(string),
				Heading:     articleData["heading"].(string),
				Content:     articleData["content"].(string),
				PublishedAt: articleData["publishedAt"].(string),
				UpdatedAt:   articleData["updatedAt"].(string),
			}
			return &article, nil
		}
	}

	article, err := s.ArticleRepo.FindArticleByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch article by ID: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(article.PublishedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(article.UpdatedAt)

	articleResponseDTO := &dto.ArticleResponseDTO{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: createdAt,
		UpdatedAt:   updatedAt,
	}

	cacheData := map[string]interface{}{
		"data": articleResponseDTO,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching article to Redis: %v\n", err)
	}

	return articleResponseDTO, nil
}
