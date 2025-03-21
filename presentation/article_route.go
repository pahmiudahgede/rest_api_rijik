package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
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
	articleAPI.Delete("/delete-article/:article_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), articleHandler.DeleteArticle)
}
