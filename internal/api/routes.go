package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/controllers"
	"github.com/pahmiudahgede/senggoldong/internal/middleware"
)

func AppRouter(app *fiber.App) {
	// # role
	app.Get("/listrole", controllers.GetAllUserRoles)
	app.Get("/listrole/:id", controllers.GetUserRoleByID)

	// # authentication
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	// # userinfo
	app.Get("/user", middleware.AuthMiddleware, controllers.GetUserInfo)
	app.Post("/user/update-password", middleware.AuthMiddleware, controllers.UpdatePassword)
	app.Put("/update-user", middleware.AuthMiddleware, controllers.UpdateUser)

	// # user set pin
	app.Get("/user/get-pin", middleware.AuthMiddleware, controllers.GetPin)
	app.Post("/user/set-pin", middleware.AuthMiddleware, controllers.CreatePin)
	app.Put("/user/update-pin", middleware.AuthMiddleware, controllers.UpdatePin)

	// # address routing
	app.Get("/list-address", middleware.AuthMiddleware, controllers.GetListAddress)
	app.Get("/address/:id", middleware.AuthMiddleware, controllers.GetAddressByID)
	app.Post("/create-address", middleware.AuthMiddleware, controllers.CreateAddress)
	app.Put("/address/:id", middleware.AuthMiddleware, controllers.UpdateAddress)
	app.Delete("/address/:id", middleware.AuthMiddleware, controllers.DeleteAddress)

	// # article
	app.Get("/articles", middleware.AuthMiddleware, controllers.GetArticles)
	app.Get("/articles/:id", middleware.AuthMiddleware, controllers.GetArticleByID)
	app.Post("/articles", middleware.AuthMiddleware, controllers.CreateArticle)
	app.Put("/articles/:id", middleware.AuthMiddleware, controllers.UpdateArticle)
	app.Delete("/articles/:id", middleware.AuthMiddleware, controllers.DeleteArticle)

	// # trash type
	app.Get("/trash-category", controllers.GetTrashCategories)
	app.Get("/trash-categorydetail/:id", controllers.GetTrashCategoryDetail)
	app.Post("/addtrash-category", controllers.CreateTrashCategory)
	app.Post("/addtrash-categorydetail", controllers.CreateTrashDetail)
	app.Put("/updatetrash-category/:id", controllers.UpdateTrashCategory)
	app.Put("/updatetrash-detail/:id", controllers.UpdateTrashDetail)
	app.Delete("/deletetrash-category/:id", controllers.DeleteTrashCategory)
	app.Delete("/deletetrash-detail/:id", controllers.DeleteTrashDetail)
}
