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
	app.Put("/update-user", middleware.AuthMiddleware, controllers.UpdateUser)
	app.Post("/user/update-password", middleware.AuthMiddleware, controllers.UpdatePassword)

	// # user set pin
	app.Post("/user/set-pin", middleware.AuthMiddleware, controllers.CreatePin)
	app.Get("/user/get-pin", middleware.AuthMiddleware, controllers.GetPin)
	app.Put("/user/update-pin", middleware.AuthMiddleware, controllers.UpdatePin)

	// # address routing
	app.Get("/list-address", middleware.AuthMiddleware, controllers.GetListAddress)
	app.Get("/address/:id", middleware.AuthMiddleware, controllers.GetAddressByID)
	app.Post("/create-address", middleware.AuthMiddleware, controllers.CreateAddress)
	app.Put("/address/:id", middleware.AuthMiddleware, controllers.UpdateAddress)
	app.Delete("/address/:id", middleware.AuthMiddleware, controllers.DeleteAddress)
}
