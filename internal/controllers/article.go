package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type ArticleController struct {
	service *services.ArticleService
}

func NewArticleController(service *services.ArticleService) *ArticleController {
	return &ArticleController{service: service}
}

func (ac *ArticleController) GetAllArticles(c *fiber.Ctx) error {
	articles, err := ac.service.GetAllArticles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch articles",
		))
	}
	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Articles fetched successfully",
		articles,
	))
}

func (ac *ArticleController) GetArticleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	article, err := ac.service.GetArticleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
			fiber.StatusNotFound,
			"Article not found",
		))
	}
	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Article fetched successfully",
		article,
	))
}

func (ac *ArticleController) CreateArticle(c *fiber.Ctx) error {
	var request dto.ArticleCreateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
		))
	}

	article, err := ac.service.CreateArticle(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Article created successfully",
		article,
	))
}

func (ac *ArticleController) UpdateArticle(c *fiber.Ctx) error {
	id := c.Params("id")
	var request dto.ArticleUpdateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
		))
	}

	article, err := ac.service.UpdateArticle(id, &request)
	if err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
				fiber.StatusNotFound,
				"Article not found",
			))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Article updated successfully",
		article,
	))
}

func (ac *ArticleController) DeleteArticle(c *fiber.Ctx) error {
	id := c.Params("id")

	err := ac.service.DeleteArticle(id)
	if err != nil {
		if err.Error() == "article not found" {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
				fiber.StatusNotFound,
				"Article not found",
			))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			fiber.StatusInternalServerError,
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Article deleted successfully",
		nil,
	))
}