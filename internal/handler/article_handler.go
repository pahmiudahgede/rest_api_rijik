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
	var requestArticleDTO dto.RequestArticleDTO
	if err := c.BodyParser(&requestArticleDTO); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := requestArticleDTO.Validate()
	if !valid {
		return utils.ValidationErrorResponse(c, errors)
	}

	articleResponse, err := h.ArticleService.CreateArticle(requestArticleDTO)
	if err != nil {
		return utils.GenericErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.CreateResponse(c, articleResponse, "Article created successfully")
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
