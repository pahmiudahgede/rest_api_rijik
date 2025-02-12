package services

import (
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

	extension := filepath.Ext(coverImage.Filename)
	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
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

	articlesCacheKey := "articles:all"
	if err := utils.DeleteData(articlesCacheKey); err != nil {

		fmt.Printf("Error deleting articles cache: %v\n", err)
	}

	articles, _, err := s.ArticleRepo.FindAllArticles(0, 0)
	if err == nil {
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

		cacheData = map[string]interface{}{
			"data": articleDTOs,
		}
		if err := utils.SetJSONData(articlesCacheKey, cacheData, time.Hour*24); err != nil {

			fmt.Printf("Error caching all articles to Redis: %v\n", err)
		}
	} else {

		fmt.Printf("Error fetching all articles: %v\n", err)
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

			total, ok := cachedData["total"].(float64)
			if !ok {
				return nil, 0, fmt.Errorf("invalid total count in cache")
			}
			return articles, int(total), nil
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
		"data":  articleDTOs,
		"total": total,
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

func (s *articleService) UpdateArticle(id string, request dto.RequestArticleDTO, coverImage *multipart.FileHeader) (*dto.ArticleResponseDTO, error) {
	article, err := s.ArticleRepo.FindArticleByID(id)
	if err != nil {
		return nil, fmt.Errorf("article not found: %v", err)
	}

	article.Title = request.Title
	article.Heading = request.Heading
	article.Content = request.Content
	article.Author = request.Author

	var coverImagePath string
	if coverImage != nil {
		coverImageDir := "./public/uploads/articles"
		if _, err := os.Stat(coverImageDir); os.IsNotExist(err) {
			err := os.MkdirAll(coverImageDir, os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf("failed to create directory for cover image: %v", err)
			}
		}

		extension := filepath.Ext(coverImage.Filename)
		if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
			return nil, fmt.Errorf("invalid file type, only .jpg, .jpeg, and .png are allowed")
		}

		coverImageFileName := fmt.Sprintf("%s_cover%s", uuid.New().String(), extension)
		coverImagePath = filepath.Join(coverImageDir, coverImageFileName)

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

		_, err = dst.ReadFrom(src)
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
	err = utils.DeleteData(articleCacheKey)
	if err != nil {
		fmt.Printf("Error deleting old cache for article: %v\n", err)
	}

	cacheData := map[string]interface{}{
		"data": articleResponseDTO,
	}
	err = utils.SetJSONData(articleCacheKey, cacheData, time.Hour*24)
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

		cacheData = map[string]interface{}{
			"data": articleDTOs,
		}
		err = utils.SetJSONData(articlesCacheKey, cacheData, time.Hour*24)
		if err != nil {
			fmt.Printf("Error caching updated articles to Redis: %v\n", err)
		}
	}

	return articleResponseDTO, nil
}
