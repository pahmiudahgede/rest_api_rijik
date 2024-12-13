package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/controllers"
	"github.com/pahmiudahgede/senggoldong/internal/middleware"
)

func AppRouter(app *fiber.App) {
	// # api group domain endpoint
	api := app.Group("/apirijikid")
	
	// # API Secure
	api.Use(middleware.APIKeyMiddleware)

	// # user initial coint
	api.Get("/user/initial-coint", controllers.GetUserInitialCoint)
	api.Get("/user/initial-coint/:id", controllers.GetUserInitialCointById)
	api.Post("/user/initial-coint", controllers.CreatePoint)
	api.Put("/user/initial-coint/:id", controllers.UpdatePoint)
	api.Delete("/user/initial-coint/:id", controllers.DeletePoint)

	//# coverage area
	api.Get("/coverage-areas", controllers.GetCoverageAreas)
	api.Get("/coverage-areas-district/:id", controllers.GetCoverageAreaByIDProvince)
	api.Get("/coverage-areas-subdistrict/:id", controllers.GetCoverageAreaByIDDistrict)
	api.Post("/coverage-areas", controllers.CreateCoverageArea)
	api.Post("/coverage-areas-district", controllers.CreateCoverageDetail)
	api.Post("/coverage-areas-subdistrict", controllers.CreateLocationSpecific)

	// # role
	api.Get("/roles", controllers.GetAllUserRoles)
	api.Get("/role/:id", controllers.GetUserRoleByID)

	// # authentication
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)

	// # userinfo
	api.Get("/user", middleware.AuthMiddleware, controllers.GetUserInfo)
	api.Post("/user/update-password", middleware.AuthMiddleware, controllers.UpdatePassword)
	api.Put("/user/update-user", middleware.AuthMiddleware, controllers.UpdateUser)

	// # user set pin
	api.Get("/user/verif-pin", middleware.AuthMiddleware, controllers.GetPin)
	api.Get("/user/cek-pin-status", middleware.AuthMiddleware, controllers.GetPinStatus)
	api.Post("/user/set-pin", middleware.AuthMiddleware, controllers.CreatePin)
	api.Put("/user/update-pin", middleware.AuthMiddleware, controllers.UpdatePin)
	api.Put("/user/update-pin", middleware.AuthMiddleware, controllers.UpdatePin)

	// # address routing
	api.Get("/addresses", middleware.AuthMiddleware, controllers.GetListAddress)
	api.Get("/address/:id", middleware.AuthMiddleware, controllers.GetAddressByID)
	api.Post("/address/create-address", middleware.AuthMiddleware, controllers.CreateAddress)
	api.Put("/address/update-address/:id", middleware.AuthMiddleware, controllers.UpdateAddress)
	api.Delete("/address/delete-address/:id", middleware.AuthMiddleware, controllers.DeleteAddress)

	// # article
	api.Get("/articles", middleware.AuthMiddleware, controllers.GetArticles)
	api.Get("/article/:id", middleware.AuthMiddleware, controllers.GetArticleByID)
	api.Post("/article/create-article", middleware.AuthMiddleware, controllers.CreateArticle)
	api.Put("/article/update-article/:id", middleware.AuthMiddleware, controllers.UpdateArticle)
	api.Delete("/article/delete-article/:id", middleware.AuthMiddleware, controllers.DeleteArticle)

	// # trash type
	api.Get("/trash-categorys", controllers.GetTrashCategories)
	api.Get("/trash-category/:id", controllers.GetTrashCategoryDetail)
	api.Post("/trash-category/create-trash-category", controllers.CreateTrashCategory)
	api.Post("/trash-category/create-trash-categorydetail", controllers.CreateTrashDetail)
	api.Put("/trash-category/update-trash-category/:id", controllers.UpdateTrashCategory)
	api.Put("/trash-category/update-trash-detail/:id", controllers.UpdateTrashDetail)
	api.Delete("/trash-category/delete-trash-category/:id", controllers.DeleteTrashCategory)
	api.Delete("/trash-category/delete-trash-detail/:id", controllers.DeleteTrashDetail)

	// # banner
	api.Get("/banners", controllers.GetBanners)
	api.Get("/banner/:id", controllers.GetBannerByID)
	api.Post("/banner/create-banner", controllers.CreateBanner)
	api.Put("/banner/update-banner/:id", controllers.UpdateBanner)
	api.Delete("/banner/delete-banner/:id", controllers.DeleteBanner)
}