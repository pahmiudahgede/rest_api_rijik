package controllers

import (
	"github.com/go-playground/validator/v10"
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
			[]dto.ArticleResponse{},
		))
	}

	for i := range articles {
		articles[i].PublishedAtFormatted = utils.FormatDateToIndonesianFormat(articles[i].PublishedAt)
		articles[i].UpdatedAtFormatted = utils.FormatDateToIndonesianFormat(articles[i].UpdatedAt)
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Articles fetched successfully",
		articles,
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

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Article fetched successfully",
		article,
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

	validate := validator.New()
	err := validate.Struct(articleUpdateRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Validation error",
			err.Error(),
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

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Article updated successfully",
		updatedArticle,
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
