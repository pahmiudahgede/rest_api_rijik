package presentation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/internal/handler"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/middleware"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func ArticleRouter(api fiber.Router) {
	articleRepo := repositories.NewArticleRepository(config.DB)
	articleService := services.NewArticleService(articleRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	articleAPI := api.Group("/article-rijik")

	articleAPI.Post("/create-article", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), articleHandler.CreateArticle)
	articleAPI.Get("/view-article", articleHandler.GetAllArticles)
	articleAPI.Get("/view-article/:article_id", articleHandler.GetArticleByID)
	articleAPI.Put("/update-article/:article_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), articleHandler.UpdateArticle)
}
