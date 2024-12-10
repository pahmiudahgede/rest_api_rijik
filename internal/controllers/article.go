package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func CreateArticle(c *fiber.Ctx) error {
	var articleRequest dto.ArticleRequest

	if err := c.BodyParser(&articleRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input",
			nil,
		))
	}

	if err := articleRequest.ValidatePostArticle(); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			err.Error(),
			nil,
		))
	}

	articleResponse, err := services.CreateArticle(&articleRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create article",
			nil,
		))
	}

	articleResponse.PublishedAtFormatted = utils.FormatDateToIndonesianFormat(articleResponse.PublishedAt)
	articleResponse.UpdatedAtFormatted = utils.FormatDateToIndonesianFormat(articleResponse.UpdatedAt)

	response := dto.FormattedResponse{
		ID:                   articleResponse.ID,
		Title:                articleResponse.Title,
		CoverImage:           articleResponse.CoverImage,
		Author:               articleResponse.Author,
		Heading:              articleResponse.Heading,
		Content:              articleResponse.Content,
		PublishedAtFormatted: articleResponse.PublishedAtFormatted,
		UpdatedAtFormatted:   articleResponse.UpdatedAtFormatted,
	}

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Article created successfully",
		response,
	))
}

func GetArticles(c *fiber.Ctx) error {

	articles, err := services.GetArticles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch articles",
			nil,
		))
	}

	if len(articles) == 0 {
		return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
			fiber.StatusOK,
			"Articles fetched successfully but data is empty",
			[]dto.FormattedResponse{},
		))
	}

	var formattedArticles []dto.FormattedResponse
	for _, article := range articles {

		article.PublishedAtFormatted = utils.FormatDateToIndonesianFormat(article.PublishedAt)
		article.UpdatedAtFormatted = utils.FormatDateToIndonesianFormat(article.UpdatedAt)

		formattedArticles = append(formattedArticles, dto.FormattedResponse{
			ID:                   article.ID,
			Title:                article.Title,
			CoverImage:           article.CoverImage,
			Author:               article.Author,
			Heading:              article.Heading,
			Content:              article.Content,
			PublishedAtFormatted: article.PublishedAtFormatted,
			UpdatedAtFormatted:   article.UpdatedAtFormatted,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Articles fetched successfully",
		formattedArticles,
	))
}

func GetArticleByID(c *fiber.Ctx) error {
	id := c.Params("id")

	article, err := services.GetArticleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.FormatResponse(
			fiber.StatusNotFound,
			"Article not found",
			nil,
		))
	}

	article.PublishedAtFormatted = utils.FormatDateToIndonesianFormat(article.PublishedAt)
	article.UpdatedAtFormatted = utils.FormatDateToIndonesianFormat(article.UpdatedAt)

	response := dto.FormattedResponse{
		ID:                   article.ID,
		Title:                article.Title,
		CoverImage:           article.CoverImage,
		Author:               article.Author,
		Heading:              article.Heading,
		Content:              article.Content,
		PublishedAtFormatted: article.PublishedAtFormatted,
		UpdatedAtFormatted:   article.UpdatedAtFormatted,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Article fetched successfully",
		response,
	))
}

func UpdateArticle(c *fiber.Ctx) error {
	id := c.Params("id")

	var articleUpdateRequest dto.ArticleUpdateRequest

	if err := c.BodyParser(&articleUpdateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid input",
			nil,
		))
	}

	if err := articleUpdateRequest.ValidateUpdateArticle(); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			err.Error(),
			nil,
		))
	}

	updatedArticle, err := services.UpdateArticle(id, &articleUpdateRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to update article",
			nil,
		))
	}

	updatedArticle.PublishedAtFormatted = utils.FormatDateToIndonesianFormat(updatedArticle.PublishedAt)
	updatedArticle.UpdatedAtFormatted = utils.FormatDateToIndonesianFormat(updatedArticle.UpdatedAt)

	response := dto.FormattedResponse{
		ID:                   updatedArticle.ID,
		Title:                updatedArticle.Title,
		CoverImage:           updatedArticle.CoverImage,
		Author:               updatedArticle.Author,
		Heading:              updatedArticle.Heading,
		Content:              updatedArticle.Content,
		PublishedAtFormatted: updatedArticle.PublishedAtFormatted,
		UpdatedAtFormatted:   updatedArticle.UpdatedAtFormatted,
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Article updated successfully",
		response,
	))
}

func DeleteArticle(c *fiber.Ctx) error {
	id := c.Params("id")

	err := services.DeleteArticle(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to delete article",
			nil,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Article deleted successfully",
		nil,
	))
}
