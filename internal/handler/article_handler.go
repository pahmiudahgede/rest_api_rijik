package handler

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	ArticleService services.ArticleService
}

func NewArticleHandler(articleService services.ArticleService) *ArticleHandler {
	return &ArticleHandler{ArticleService: articleService}
}

func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	var request dto.RequestArticleDTO

	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	coverImage, err := c.FormFile("coverImage")
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Cover image is required")
	}

	articleResponse, err := h.ArticleService.CreateArticle(request, coverImage)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.CreateResponse(c, articleResponse, "Article created successfully")
}

func (h *ArticleHandler) GetAllArticles(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil || page < 1 {
		page = 0
	}

	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil || limit < 1 {
		limit = 0
	}

	articles, totalArticles, err := h.ArticleService.GetAllArticles(page, limit)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusInternalServerError, "Failed to fetch articles")
	}

	fmt.Printf("Total Articles: %d\n", totalArticles)

	if page == 0 && limit == 0 {
		return utils.NonPaginatedResponse(c, articles, totalArticles, "Articles fetched successfully")
	}

	return utils.PaginatedResponse(c, articles, page, limit, totalArticles, "Articles fetched successfully")
}

func (h *ArticleHandler) GetArticleByID(c *fiber.Ctx) error {
	id := c.Params("article_id")
	if id == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Article ID is required")
	}

	article, err := h.ArticleService.GetArticleByID(id)
	if err != nil {
		return utils.GenericResponse(c, fiber.StatusNotFound, "Article not found")
	}

	return utils.SuccessResponse(c, article, "Article fetched successfully")
}

func (h *ArticleHandler) UpdateArticle(c *fiber.Ctx) error {
	id := c.Params("article_id")
	if id == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Article ID is required")
	}

	var request dto.RequestArticleDTO
	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	var coverImage *multipart.FileHeader
	coverImage, err := c.FormFile("coverImage")
	if err != nil && err.Error() != "no such file" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Cover image is required")
	}

	articleResponse, err := h.ArticleService.UpdateArticle(id, request, coverImage)
	if err != nil {
		if err.Error() == fmt.Sprintf("article with ID %s not found", id) {
			return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())

	}

	return utils.SuccessResponse(c, articleResponse, "Article updated successfully")
}

func (h *ArticleHandler) DeleteArticle(c *fiber.Ctx) error {
	id := c.Params("article_id")
	if id == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Article ID is required")
	}

	err := h.ArticleService.DeleteArticle(id)
	if err != nil {

		if err.Error() == fmt.Sprintf("article with ID %s not found", id) {
			return utils.GenericResponse(c, fiber.StatusInternalServerError, err.Error())
		}

		return utils.GenericResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Article deleted successfully")
}
