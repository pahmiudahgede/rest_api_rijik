package router

import (
	"os"

	"rijig/middleware"
	"rijig/presentation"
	presentationn "rijig/presentation/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	apa := app.Group(os.Getenv("BASE_URL"))
	apa.Static("/uploads", "./public"+os.Getenv("BASE_URL")+"/uploads")

	api := app.Group(os.Getenv("BASE_URL"))
	api.Use(middleware.APIKeyMiddleware)

	// || auth router || //
	// presentation.AuthRouter(api)
	presentationn.AuthAdminRouter(api)
	presentationn.AuthPengelolaRouter(api)
	presentationn.AuthPengepulRouter(api)
	presentationn.AuthMasyarakatRouter(api)
	// || auth router || //
	presentation.IdentityCardRouter(api)
	presentation.CompanyProfileRouter(api)
	presentation.RequestPickupRouter(api)
	presentation.PickupMatchingRouter(api)
	presentation.PickupRatingRouter(api)

	presentation.CollectorRouter(api)
	presentation.TrashCartRouter(api)

	presentation.UserProfileRouter(api)
	presentation.UserPinRouter(api)
	presentation.RoleRouter(api)
	presentation.WilayahRouter(api)
	presentation.AddressRouter(api)
	presentation.ArticleRouter(api)
	presentation.BannerRouter(api)
	presentation.InitialCointRoute(api)
	presentation.AboutRouter(api)
	presentation.TrashRouter(api)
	presentation.CoverageAreaRouter(api)
	presentation.StoreRouter(api)
	presentation.ProductRouter(api)
	presentation.WhatsAppRouter(api)

}
