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

	validate := validator.New()
	err := validate.Struct(articleRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Validation error",
			err.Error(),
		))
	}

	createdArticle, err := services.CreateArticle(&articleRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create article",
			nil,
		))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Article created successfully",
		createdArticle,
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
