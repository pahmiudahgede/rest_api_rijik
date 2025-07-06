package article

import (
	"rijig/config"
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func ArticleRouter(api fiber.Router) {
	articleRepo := NewArticleRepository(config.DB)
	articleService := NewArticleService(articleRepo)
	articleHandler := NewArticleHandler(articleService)

	articleAPI := api.Group("/article")

	articleAPI.Post("/create", middleware.AuthMiddleware(), middleware.RequireRoles(utils.RoleAdministrator), articleHandler.CreateArticle)
	articleAPI.Get("/view", articleHandler.GetAllArticles)
	articleAPI.Get("/view/:article_id", articleHandler.GetArticleByID)
	articleAPI.Put("/update/:article_id", middleware.AuthMiddleware(), middleware.RequireRoles(utils.RoleAdministrator), articleHandler.UpdateArticle)
	articleAPI.Delete("/delete/:article_id", middleware.AuthMiddleware(), middleware.RequireRoles(utils.RoleAdministrator), articleHandler.DeleteArticle)
}
