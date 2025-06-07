package article

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"rijig/model"
	"rijig/utils"

	"github.com/google/uuid"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, request RequestArticleDTO, coverImage *multipart.FileHeader) (*ArticleResponseDTO, error)
	GetAllArticles(ctx context.Context, page, limit int) ([]ArticleResponseDTO, int64, error)
	GetArticleByID(ctx context.Context, id string) (*ArticleResponseDTO, error)
	UpdateArticle(ctx context.Context, id string, request RequestArticleDTO, coverImage *multipart.FileHeader) (*ArticleResponseDTO, error)
	DeleteArticle(ctx context.Context, id string) error
}

type articleService struct {
	articleRepo ArticleRepository
}

func NewArticleService(articleRepo ArticleRepository) ArticleService {
	return &articleService{articleRepo}
}

func (s *articleService) transformToDTO(article *model.Article) (*ArticleResponseDTO, error) {
	publishedAt, err := utils.FormatDateToIndonesianFormat(article.PublishedAt)
	if err != nil {
		publishedAt = ""
	}

	updatedAt, err := utils.FormatDateToIndonesianFormat(article.UpdatedAt)
	if err != nil {
		updatedAt = ""
	}

	return &ArticleResponseDTO{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: publishedAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (s *articleService) transformToDTOs(articles []model.Article) ([]ArticleResponseDTO, error) {
	var articleDTOs []ArticleResponseDTO

	for _, article := range articles {
		dto, err := s.transformToDTO(&article)
		if err != nil {
			return nil, fmt.Errorf("failed to transform article %s: %w", article.ID, err)
		}
		articleDTOs = append(articleDTOs, *dto)
	}

	return articleDTOs, nil
}

func (s *articleService) saveCoverArticle(coverArticle *multipart.FileHeader) (string, error) {
	if coverArticle == nil {
		return "", fmt.Errorf("cover image is required")
	}

	pathImage := "/uploads/articles/"
	coverArticleDir := "./public" + os.Getenv("BASE_URL") + pathImage

	if err := os.MkdirAll(coverArticleDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory for cover article: %w", err)
	}

	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".svg": true}
	extension := filepath.Ext(coverArticle.Filename)
	if !allowedExtensions[extension] {
		return "", fmt.Errorf("invalid file type, only .jpg, .jpeg, .png, and .svg are allowed")
	}

	coverArticleFileName := fmt.Sprintf("%s_coverarticle%s", uuid.New().String(), extension)
	coverArticlePath := filepath.Join(coverArticleDir, coverArticleFileName)

	src, err := coverArticle.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(coverArticlePath)
	if err != nil {
		return "", fmt.Errorf("failed to create cover article file: %w", err)
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("failed to save cover article: %w", err)
	}

	return fmt.Sprintf("%s%s", pathImage, coverArticleFileName), nil
}

func (s *articleService) deleteCoverArticle(imagePath string) error {
	if imagePath == "" {
		return nil
	}

	baseDir := "./public/" + os.Getenv("BASE_URL")
	absolutePath := baseDir + imagePath

	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		log.Printf("Image file not found (already deleted?): %s", absolutePath)
		return nil
	}

	if err := os.Remove(absolutePath); err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	log.Printf("Image deleted successfully: %s", absolutePath)
	return nil
}

func (s *articleService) invalidateArticleCache(articleID string) {

	articleCacheKey := fmt.Sprintf("article:%s", articleID)
	if err := utils.DeleteCache(articleCacheKey); err != nil {
		log.Printf("Error deleting article cache: %v", err)
	}

	if err := utils.ScanAndDelete("articles:*"); err != nil {
		log.Printf("Error deleting articles cache: %v", err)
	}
}

func (s *articleService) CreateArticle(ctx context.Context, request RequestArticleDTO, coverImage *multipart.FileHeader) (*ArticleResponseDTO, error) {
	coverArticlePath, err := s.saveCoverArticle(coverImage)
	if err != nil {
		return nil, fmt.Errorf("failed to save cover image: %w", err)
	}

	article := model.Article{
		Title:      request.Title,
		CoverImage: coverArticlePath,
		Author:     request.Author,
		Heading:    request.Heading,
		Content:    request.Content,
	}

	if err := s.articleRepo.CreateArticle(ctx, &article); err != nil {

		if deleteErr := s.deleteCoverArticle(coverArticlePath); deleteErr != nil {
			log.Printf("Failed to clean up image after create failure: %v", deleteErr)
		}
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	articleDTO, err := s.transformToDTO(&article)
	if err != nil {
		return nil, fmt.Errorf("failed to transform article: %w", err)
	}

	cacheKey := fmt.Sprintf("article:%s", article.ID)
	if err := utils.SetCache(cacheKey, articleDTO, time.Hour*24); err != nil {
		log.Printf("Error caching article: %v", err)
	}

	s.invalidateArticleCache("")

	return articleDTO, nil
}

func (s *articleService) GetAllArticles(ctx context.Context, page, limit int) ([]ArticleResponseDTO, int64, error) {

	var cacheKey string
	if page <= 0 || limit <= 0 {
		cacheKey = "articles:all"
	} else {
		cacheKey = fmt.Sprintf("articles:page:%d:limit:%d", page, limit)
	}

	type CachedArticlesData struct {
		Articles []ArticleResponseDTO `json:"articles"`
		Total    int64                `json:"total"`
	}

	var cachedData CachedArticlesData
	if err := utils.GetCache(cacheKey, &cachedData); err == nil {
		return cachedData.Articles, cachedData.Total, nil
	}

	articles, total, err := s.articleRepo.FindAllArticles(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch articles: %w", err)
	}

	articleDTOs, err := s.transformToDTOs(articles)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to transform articles: %w", err)
	}

	cacheData := CachedArticlesData{
		Articles: articleDTOs,
		Total:    total,
	}
	if err := utils.SetCache(cacheKey, cacheData, time.Hour*24); err != nil {
		log.Printf("Error caching articles: %v", err)
	}

	return articleDTOs, total, nil
}

func (s *articleService) GetArticleByID(ctx context.Context, id string) (*ArticleResponseDTO, error) {
	if id == "" {
		return nil, fmt.Errorf("article ID cannot be empty")
	}

	cacheKey := fmt.Sprintf("article:%s", id)

	var cachedArticle ArticleResponseDTO
	if err := utils.GetCache(cacheKey, &cachedArticle); err == nil {
		return &cachedArticle, nil
	}

	article, err := s.articleRepo.FindArticleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch article: %w", err)
	}

	articleDTO, err := s.transformToDTO(article)
	if err != nil {
		return nil, fmt.Errorf("failed to transform article: %w", err)
	}

	if err := utils.SetCache(cacheKey, articleDTO, time.Hour*24); err != nil {
		log.Printf("Error caching article: %v", err)
	}

	return articleDTO, nil
}

func (s *articleService) UpdateArticle(ctx context.Context, id string, request RequestArticleDTO, coverImage *multipart.FileHeader) (*ArticleResponseDTO, error) {
	if id == "" {
		return nil, fmt.Errorf("article ID cannot be empty")
	}

	existingArticle, err := s.articleRepo.FindArticleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("article not found: %w", err)
	}

	oldCoverImage := existingArticle.CoverImage
	var newCoverPath string

	if coverImage != nil {
		newCoverPath, err = s.saveCoverArticle(coverImage)
		if err != nil {
			return nil, fmt.Errorf("failed to save new cover image: %w", err)
		}
	}

	updatedArticle := &model.Article{
		Title:      request.Title,
		Author:     request.Author,
		Heading:    request.Heading,
		Content:    request.Content,
		CoverImage: existingArticle.CoverImage,
	}

	if newCoverPath != "" {
		updatedArticle.CoverImage = newCoverPath
	}

	if err := s.articleRepo.UpdateArticle(ctx, id, updatedArticle); err != nil {

		if newCoverPath != "" {
			s.deleteCoverArticle(newCoverPath)
		}
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	if newCoverPath != "" && oldCoverImage != "" {
		if err := s.deleteCoverArticle(oldCoverImage); err != nil {
			log.Printf("Warning: failed to delete old cover image: %v", err)
		}
	}

	article, err := s.articleRepo.FindArticleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated article: %w", err)
	}

	articleDTO, err := s.transformToDTO(article)
	if err != nil {
		return nil, fmt.Errorf("failed to transform updated article: %w", err)
	}

	cacheKey := fmt.Sprintf("article:%s", id)
	if err := utils.SetCache(cacheKey, articleDTO, time.Hour*24); err != nil {
		log.Printf("Error caching updated article: %v", err)
	}

	s.invalidateArticleCache(id)

	return articleDTO, nil
}

func (s *articleService) DeleteArticle(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("article ID cannot be empty")
	}

	article, err := s.articleRepo.FindArticleByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find article: %w", err)
	}

	if err := s.articleRepo.DeleteArticle(ctx, id); err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	if err := s.deleteCoverArticle(article.CoverImage); err != nil {
		log.Printf("Warning: failed to delete cover image: %v", err)
	}

	s.invalidateArticleCache(id)

	return nil
}
