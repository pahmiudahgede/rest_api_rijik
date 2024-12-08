package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/controllers"
	"github.com/pahmiudahgede/senggoldong/internal/middleware"
)

func AppRouter(app *fiber.App) {
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	app.Get("/user", middleware.AuthMiddleware, controllers.GetUserInfo)
	app.Put("/update-user", middleware.AuthMiddleware, controllers.UpdateUser)
	app.Post("/user/update-password", middleware.AuthMiddleware, controllers.UpdatePassword)

	app.Get("/list-address", middleware.AuthMiddleware, controllers.GetListAddress)
	app.Get("/address/:id", middleware.AuthMiddleware, controllers.GetAddressByID)
	app.Post("/create-address", middleware.AuthMiddleware, controllers.CreateAddress)
	app.Put("/address/:id", middleware.AuthMiddleware, controllers.UpdateAddress)
	app.Delete("/address/:id", middleware.AuthMiddleware, controllers.DeleteAddress)
}
