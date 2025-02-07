package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
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
		return utils.GenericErrorResponse(c, fiber.StatusBadRequest, "Cover image is required")
	}

	articleResponse, err := h.ArticleService.CreateArticle(request, coverImage)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, articleResponse, "Article created successfully")
}

func (h *ArticleHandler) GetAllArticles(c *fiber.Ctx) error {

	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		limit = 0
	}

	article, totalArticle, err := h.ArticleService.GetAllArticles(page, limit)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch article")
	}

	if page > 0 && limit > 0 {
		return utils.PaginatedResponse(c, article, page, limit, totalArticle, "Article fetched successfully")
	}

	return utils.NonPaginatedResponse(c, article, totalArticle, "Article fetched successfully")

}

func (h *ArticleHandler) GetArticleByID(c *fiber.Ctx) error {

	id := c.Params("article_id")
	if id == "" {
		return utils.GenericErrorResponse(c, fiber.StatusBadRequest, "Article ID is required")
	}

	article, err := h.ArticleService.GetArticleByID(id)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, article, "Article fetched successfully")
}
