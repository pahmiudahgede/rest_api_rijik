package services

import (
	"time"

	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

func CreateArticle(articleRequest *dto.ArticleRequest) (*dto.ArticleResponse, error) {
	article := domain.Article{
		Title:       articleRequest.Title,
		CoverImage:  articleRequest.CoverImage,
		Author:      articleRequest.Author,
		Heading:     articleRequest.Heading,
		Content:     articleRequest.Content,
		PublishedAt: articleRequest.PublishedAt,
		UpdatedAt:   articleRequest.PublishedAt,
	}

	err := repositories.CreateArticle(&article)
	if err != nil {
		return nil, err
	}

	return &dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: article.PublishedAt,
		UpdatedAt:   article.UpdatedAt,
	}, nil
}

func GetArticles() ([]dto.ArticleResponse, error) {
	articles, err := repositories.GetArticles()
	if err != nil {
		return nil, err
	}
	var response []dto.ArticleResponse
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
	return response, nil
}

func GetArticleByID(id string) (dto.ArticleResponse, error) {
	article, err := repositories.GetArticleByID(id)
	if err != nil {
		return dto.ArticleResponse{}, err
	}
	return dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: article.PublishedAt,
		UpdatedAt:   article.UpdatedAt,
	}, nil
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

	return &dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		CoverImage:  article.CoverImage,
		Author:      article.Author,
		Heading:     article.Heading,
		Content:     article.Content,
		PublishedAt: article.PublishedAt,
		UpdatedAt:   article.UpdatedAt,
	}, nil
}

func DeleteArticle(id string) error {

	err := repositories.DeleteArticle(id)
	if err != nil {
		return err
	}
	return nil
}
