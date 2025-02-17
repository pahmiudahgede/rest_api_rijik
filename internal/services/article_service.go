package services

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type ArticleService interface {
	CreateArticle(request dto.RequestArticleDTO, coverImage *multipart.FileHeader) (*dto.ArticleResponseDTO, error)
	GetAllArticles(page, limit int) ([]dto.ArticleResponseDTO, int, error)
	GetArticleByID(id string) (*dto.ArticleResponseDTO, error)
	UpdateArticle(id string, request dto.RequestArticleDTO, coverImage *multipart.FileHeader) (*dto.ArticleResponseDTO, error)
	DeleteArticle(id string) error
}

type articleService struct {
	ArticleRepo repositories.ArticleRepository
}

func NewArticleService(articleRepo repositories.ArticleRepository) ArticleService {
	return &articleService{ArticleRepo: articleRepo}
}

func (s *articleService) CreateArticle(request dto.RequestArticleDTO, coverImage *multipart.FileHeader) (*dto.ArticleResponseDTO, error) {

	coverImageDir := "./public/uploads/articles"
	if err := os.MkdirAll(coverImageDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory for cover image: %v", err)
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	extension := filepath.Ext(coverImage.Filename)
	if !allowedExtensions[extension] {
		return nil, fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed")
	}

	coverImageFileName := fmt.Sprintf("%s_cover%s", uuid.New().String(), extension)
	coverImagePath := filepath.Join(coverImageDir, coverImageFileName)

	src, err := coverImage.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(coverImagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create cover image file: %v", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return nil, fmt.Errorf("failed to save cover image: %v", err)
	}

	article := model.Article{
		Title:      request.Title,
		CoverImage: coverImagePath,
		Author:     request.Author,
		Heading:    request.Heading,
		Content:    request.Content,
	}

	if err := s.ArticleRepo.CreateArticle(&article); err != nil {
		return nil, fmt.Errorf("failed to create article: %v", err)
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

	cacheKey := fmt.Sprintf("article:%s", article.ID)
	cacheData := map[string]interface{}{
		"data": articleResponseDTO,
	}
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching article to Redis: %v\n", err)
	}

	articles, total, err := s.ArticleRepo.FindAllArticles(0, 0)
	if err != nil {
		fmt.Printf("Error fetching all articles: %v\n", err)
	}

	var articleDTOs []dto.ArticleResponseDTO
	for _, a := range articles {
		createdAt, _ := utils.FormatDateToIndonesianFormat(a.PublishedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(a.UpdatedAt)

		articleDTOs = append(articleDTOs, dto.ArticleResponseDTO{
			ID:          a.ID,
			Title:       a.Title,
			CoverImage:  a.CoverImage,
			Author:      a.Author,
			Heading:     a.Heading,
			Content:     a.Content,
			PublishedAt: createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	articlesCacheKey := "articles:all"
	cacheData = map[string]interface{}{
		"data":  articleDTOs,
		"total": total,
	}
	if err := utils.SetJSONData(articlesCacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching all articles to Redis: %v\n", err)
	}

	return articleResponseDTO, nil
}

func (s *articleService) GetAllArticles(page, limit int) ([]dto.ArticleResponseDTO, int, error) {
	var cacheKey string

	if page == 0 && limit == 0 {
		cacheKey = "articles:all"
		cachedData, err := utils.GetJSONData(cacheKey)
		if err == nil && cachedData != nil {
			if data, ok := cachedData["data"].([]interface{}); ok {
				var articles []dto.ArticleResponseDTO
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

				if total, ok := cachedData["total"].(float64); ok {
					fmt.Printf("Cached Total Articles: %f\n", total)
					return articles, int(total), nil
				} else {
					fmt.Println("Total articles not found in cache, using 0 as fallback.")
					return articles, 0, nil
				}
			}
		}
	}

	articles, total, err := s.ArticleRepo.FindAllArticles(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch articles: %v", err)
	}

	fmt.Printf("Total Articles from Database: %d\n", total)

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

	cacheKey = fmt.Sprintf("articles_page:%d_limit:%d", page, limit)
	cacheData := map[string]interface{}{
		"data":  articleDTOs,
		"total": total,
	}

	fmt.Printf("Setting cache with total: %d\n", total)
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching articles to Redis: %v\n", err)
	}

	return articleDTOs, total, nil
}

func (s *articleService) GetArticleByID(id string) (*dto.ArticleResponseDTO, error) {

	cacheKey := fmt.Sprintf("article:%s", id)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		articleResponse := &dto.ArticleResponseDTO{}
		if data, ok := cachedData["data"].(string); ok {
			if err := json.Unmarshal([]byte(data), articleResponse); err == nil {
				return articleResponse, nil
			}
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
	if err := utils.SetJSONData(cacheKey, cacheData, time.Hour*24); err != nil {
		fmt.Printf("Error caching article to Redis: %v\n", err)
	}

	return articleResponseDTO, nil
}

func (s *articleService) UpdateArticle(id string, request dto.RequestArticleDTO, coverImage *multipart.FileHeader) (*dto.ArticleResponseDTO, error) {
	article, err := s.ArticleRepo.FindArticleByID(id)
	if err != nil {
		return nil, fmt.Errorf("article not found: %v", id)
	}

	article.Title = request.Title
	article.Heading = request.Heading
	article.Content = request.Content
	article.Author = request.Author

	var coverImagePath string
	if coverImage != nil {

		coverImagePath, err = s.saveCoverImage(coverImage, article.CoverImage)
		if err != nil {
			return nil, fmt.Errorf("failed to save cover image: %v", err)
		}
		article.CoverImage = coverImagePath
	}

	err = s.ArticleRepo.UpdateArticle(id, article)
	if err != nil {
		return nil, fmt.Errorf("failed to update article: %v", err)
	}

	updatedArticle, err := s.ArticleRepo.FindArticleByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated article: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(updatedArticle.PublishedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(updatedArticle.UpdatedAt)

	articleResponseDTO := &dto.ArticleResponseDTO{
		ID:          updatedArticle.ID,
		Title:       updatedArticle.Title,
		CoverImage:  updatedArticle.CoverImage,
		Author:      updatedArticle.Author,
		Heading:     updatedArticle.Heading,
		Content:     updatedArticle.Content,
		PublishedAt: createdAt,
		UpdatedAt:   updatedAt,
	}

	articleCacheKey := fmt.Sprintf("article:%s", updatedArticle.ID)
	err = utils.SetJSONData(articleCacheKey, map[string]interface{}{"data": articleResponseDTO}, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching updated article to Redis: %v\n", err)
	}

	articlesCacheKey := "articles:all"
	err = utils.DeleteData(articlesCacheKey)
	if err != nil {
		fmt.Printf("Error deleting articles cache: %v\n", err)
	}

	articles, _, err := s.ArticleRepo.FindAllArticles(0, 0)
	if err != nil {
		fmt.Printf("Error fetching all articles: %v\n", err)
	} else {
		var articleDTOs []dto.ArticleResponseDTO
		for _, a := range articles {
			createdAt, _ := utils.FormatDateToIndonesianFormat(a.PublishedAt)
			updatedAt, _ := utils.FormatDateToIndonesianFormat(a.UpdatedAt)

			articleDTOs = append(articleDTOs, dto.ArticleResponseDTO{
				ID:          a.ID,
				Title:       a.Title,
				CoverImage:  a.CoverImage,
				Author:      a.Author,
				Heading:     a.Heading,
				Content:     a.Content,
				PublishedAt: createdAt,
				UpdatedAt:   updatedAt,
			})
		}

		cacheData := map[string]interface{}{
			"data": articleDTOs,
		}
		err = utils.SetJSONData(articlesCacheKey, cacheData, time.Hour*24)
		if err != nil {
			fmt.Printf("Error caching updated articles to Redis: %v\n", err)
		}
	}

	return articleResponseDTO, nil
}

func (s *articleService) saveCoverImage(coverImage *multipart.FileHeader, oldImagePath string) (string, error) {
	coverImageDir := "./public/uploads/articles"
	if _, err := os.Stat(coverImageDir); os.IsNotExist(err) {
		if err := os.MkdirAll(coverImageDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory for cover image: %v", err)
		}
	}

	extension := filepath.Ext(coverImage.Filename)
	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed")
	}

	coverImageFileName := fmt.Sprintf("%s_cover%s", uuid.New().String(), extension)
	coverImagePath := filepath.Join(coverImageDir, coverImageFileName)

	if oldImagePath != "" {
		err := os.Remove(oldImagePath)
		if err != nil {
			fmt.Printf("Failed to delete old cover image: %v\n", err)
		} else {
			fmt.Printf("Successfully deleted old cover image: %s\n", oldImagePath)
		}
	}

	src, err := coverImage.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(coverImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create cover image file: %v", err)
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", fmt.Errorf("failed to save cover image: %v", err)
	}

	return coverImagePath, nil
}

func (s *articleService) DeleteArticle(id string) error {
	article, err := s.ArticleRepo.FindArticleByID(id)
	if err != nil {
		return fmt.Errorf("failed to find article: %v", id)
	}

	if article.CoverImage != "" {
		err := os.Remove(article.CoverImage)
		if err != nil {
			fmt.Printf("Failed to delete cover image: %v\n", err)
		} else {
			fmt.Printf("Successfully deleted cover image: %s\n", article.CoverImage)
		}
	}

	err = s.ArticleRepo.DeleteArticle(id)
	if err != nil {
		return fmt.Errorf("failed to delete article: %v", err)
	}

	articleCacheKey := fmt.Sprintf("article:%s", id)
	err = utils.DeleteData(articleCacheKey)
	if err != nil {
		fmt.Printf("Error deleting cache for article: %v\n", err)
	}

	articlesCacheKey := "articles:all"
	err = utils.DeleteData(articlesCacheKey)
	if err != nil {
		fmt.Printf("Error deleting cache for all articles: %v\n", err)
	}

	return nil
}
