package router

import (
	"os"

	"rijig/middleware"
	"rijig/presentation"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Static(os.Getenv("BASE_URL")+"/uploads", "./public"+os.Getenv("BASE_URL")+"/uploads")

	api := app.Group(os.Getenv("BASE_URL"))
	api.Use(middleware.APIKeyMiddleware)

	presentation.AuthRouter(api)
	presentation.UserProfileRouter(api)
	presentation.UserPinRouter(api)
	presentation.RoleRouter(api)
	presentation.WilayahRouter(api)
	presentation.AddressRouter(api)
	presentation.ArticleRouter(api)
	presentation.BannerRouter(api)
	presentation.InitialCointRoute(api)
	presentation.TrashRouter(api)
	presentation.StoreRouter(api)
	presentation.ProductRouter(api)
}
