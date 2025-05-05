package router

import (
	"os"

	"rijig/middleware"
	"rijig/presentation"
	presentationn "rijig/presentation/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group(os.Getenv("BASE_URL"))
	api.Use(middleware.APIKeyMiddleware)
	api.Static("/uploads", "./public"+os.Getenv("BASE_URL")+"/uploads")

	// || auth router || //
	// presentation.AuthRouter(api)
	presentationn.AdminAuthRouter(api)
	presentationn.AuthPengelolaRouter(api)
	presentationn.AuthPengepulRouter(api)
	presentationn.AuthMasyarakatRouter(api)
	// || auth router || //
	presentation.IdentityCardRouter(api)

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
