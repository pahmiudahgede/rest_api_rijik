package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/middleware"
	"github.com/pahmiudahgede/senggoldong/presentation"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/apirijikid/v2")
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
}
