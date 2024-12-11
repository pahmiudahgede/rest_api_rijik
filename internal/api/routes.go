package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/controllers"
	"github.com/pahmiudahgede/senggoldong/internal/middleware"
)

func AppRouter(app *fiber.App) {
	// # API Secure
	app.Use(middleware.APIKeyMiddleware)

	// # role
	app.Get("/apirijikid/roles", controllers.GetAllUserRoles)
	app.Get("/apirijikid/role/:id", controllers.GetUserRoleByID)

	// # authentication
	app.Post("/apirijikid/register", controllers.Register)
	app.Post("/apirijikid/login", controllers.Login)

	// # userinfo
	app.Get("/apirijikid/user", middleware.AuthMiddleware, controllers.GetUserInfo)
	app.Post("/apirijikid/user/update-password", middleware.AuthMiddleware, controllers.UpdatePassword)
	app.Put("/apirijikid/user/update-user", middleware.AuthMiddleware, controllers.UpdateUser)

	// # user set pin
	app.Get("/apirijikid/user/verif-pin", middleware.AuthMiddleware, controllers.GetPin)
	app.Get("/apirijikid/user/cek-pin-status", middleware.AuthMiddleware, controllers.GetPinStatus)
	app.Post("/apirijikid/user/set-pin", middleware.AuthMiddleware, controllers.CreatePin)
	app.Put("/apirijikid/user/update-pin", middleware.AuthMiddleware, controllers.UpdatePin)
	app.Put("/apirijikid/user/update-pin", middleware.AuthMiddleware, controllers.UpdatePin)

	// # address routing
	app.Get("/apirijikid/addresses", middleware.AuthMiddleware, controllers.GetListAddress)
	app.Get("/apirijikid/address/:id", middleware.AuthMiddleware, controllers.GetAddressByID)
	app.Post("/apirijikid/address/create-address", middleware.AuthMiddleware, controllers.CreateAddress)
	app.Put("/apirijikid/address/update-address/:id", middleware.AuthMiddleware, controllers.UpdateAddress)
	app.Delete("/apirijikid/address/delete-address/:id", middleware.AuthMiddleware, controllers.DeleteAddress)

	// # article
	app.Get("/apirijikid/articles", middleware.AuthMiddleware, controllers.GetArticles)
	app.Get("/apirijikid/article/:id", middleware.AuthMiddleware, controllers.GetArticleByID)
	app.Post("/apirijikid/article/create-article", middleware.AuthMiddleware, controllers.CreateArticle)
	app.Put("/apirijikid/article/update-article/:id", middleware.AuthMiddleware, controllers.UpdateArticle)
	app.Delete("/apirijikid/article/delete-article/:id", middleware.AuthMiddleware, controllers.DeleteArticle)

	// # trash type
	app.Get("/apirijikid/trash-categorys", controllers.GetTrashCategories)
	app.Get("/apirijikid/trash-category/:id", controllers.GetTrashCategoryDetail)
	app.Post("/apirijikid/trash-category/create-trash-category", controllers.CreateTrashCategory)
	app.Post("/apirijikid/trash-category/create-trash-categorydetail", controllers.CreateTrashDetail)
	app.Put("/apirijikid/trash-category/update-trash-category/:id", controllers.UpdateTrashCategory)
	app.Put("/apirijikid/trash-category/update-trash-detail/:id", controllers.UpdateTrashDetail)
	app.Delete("/apirijikid/trash-category/delete-trash-category/:id", controllers.DeleteTrashCategory)
	app.Delete("/apirijikid/trash-category/delete-trash-detail/:id", controllers.DeleteTrashDetail)

	// # banner
	app.Get("/apirijikid/banners", controllers.GetBanners)
	app.Get("/apirijikid/banner/:id", controllers.GetBannerByID)
	app.Post("/apirijikid/banner/create-banner", controllers.CreateBanner)
	app.Put("/apirijikid/banner/update-banner/:id", controllers.UpdateBanner)
	app.Delete("/apirijikid/banner/delete-banner/:id", controllers.DeleteBanner)
}
