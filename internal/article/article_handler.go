package article

import (
	"mime/multipart"
	"rijig/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	articleService ArticleService
}

func NewArticleHandler(articleService ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService}
}

func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	var request RequestArticleDTO

	if err := c.BodyParser(&request); err != nil {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Invalid request body", map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.ValidateRequestArticleDTO()
	if !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	coverImage, err := c.FormFile("coverImage")
	if err != nil {
		return utils.BadRequest(c, "Cover image is required")
	}

	articleResponse, err := h.articleService.CreateArticle(c.Context(), request, coverImage)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.SuccessWithData(c, "Article created successfully", articleResponse)
}

func (h *ArticleHandler) GetAllArticles(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "0"))
	if err != nil || page < 0 {
		page = 0
	}

	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil || limit < 0 {
		limit = 0
	}

	articles, totalArticles, err := h.articleService.GetAllArticles(c.Context(), page, limit)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch articles")
	}

	responseData := map[string]interface{}{
		"articles": articles,
		"total":    int(totalArticles),
	}

	if page == 0 && limit == 0 {
		return utils.SuccessWithData(c, "Articles fetched successfully", responseData)
	}

	return utils.SuccessWithPagination(c, "Articles fetched successfully", responseData, page, limit)
}

func (h *ArticleHandler) GetArticleByID(c *fiber.Ctx) error {
	id := c.Params("article_id")
	if id == "" {
		return utils.BadRequest(c, "Article ID is required")
	}

	article, err := h.articleService.GetArticleByID(c.Context(), id)
	if err != nil {
		return utils.NotFound(c, "Article not found")
	}

	return utils.SuccessWithData(c, "Article fetched successfully", article)
}

func (h *ArticleHandler) UpdateArticle(c *fiber.Ctx) error {
	id := c.Params("article_id")
	if id == "" {
		return utils.BadRequest(c, "Article ID is required")
	}

	var request RequestArticleDTO
	if err := c.BodyParser(&request); err != nil {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Invalid request body", map[string][]string{"body": {"Invalid body"}})
	}

	errors, valid := request.ValidateRequestArticleDTO()
	if !valid {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	var coverImage *multipart.FileHeader
	coverImage, err := c.FormFile("coverImage")

	if err != nil && err.Error() != "no such file" && err.Error() != "there is no uploaded file associated with the given key" {
		return utils.BadRequest(c, "Invalid cover image")
	}

	articleResponse, err := h.articleService.UpdateArticle(c.Context(), id, request, coverImage)
	if err != nil {
		if isNotFoundError(err) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalServerError(c, err.Error())
	}

	return utils.SuccessWithData(c, "Article updated successfully", articleResponse)
}

func (h *ArticleHandler) DeleteArticle(c *fiber.Ctx) error {
	id := c.Params("article_id")
	if id == "" {
		return utils.BadRequest(c, "Article ID is required")
	}

	err := h.articleService.DeleteArticle(c.Context(), id)
	if err != nil {
		if isNotFoundError(err) {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalServerError(c, err.Error())
	}

	return utils.Success(c, "Article deleted successfully")
}

func isNotFoundError(err error) bool {
	return err != nil && (err.Error() == "article not found" ||
		err.Error() == "failed to find article: record not found" ||
		false)
}
