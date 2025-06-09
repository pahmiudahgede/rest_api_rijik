package router

import (
	"os"

	"rijig/internal/article"
	"rijig/internal/authentication"
	"rijig/internal/company"
	"rijig/internal/identitycart"
	"rijig/internal/role"
	"rijig/internal/userpin"
	"rijig/internal/whatsapp"
	"rijig/middleware"

	// "rijig/presentation"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	apa := app.Group(os.Getenv("BASE_URL"))
	whatsapp.WhatsAppRouter(apa)
	apa.Static("/uploads", "./public"+os.Getenv("BASE_URL")+"/uploads")

	api := app.Group(os.Getenv("BASE_URL"))
	api.Use(middleware.APIKeyMiddleware)

	authentication.AuthenticationRouter(api)
	identitycart.UserIdentityCardRoute(api)
	company.CompanyRouter(api)
	userpin.UsersPinRoute(api)
	role.UserRoleRouter(api)

	article.ArticleRouter(api)

	// || auth router || //
	// presentation.AuthRouter(api)
	// presentationn.AuthAdminRouter(api)
	// presentationn.AuthPengelolaRouter(api)
	// presentationn.AuthPengepulRouter(api)
	// presentationn.AuthMasyarakatRouter(api)
	// || auth router || //
	// presentation.IdentityCardRouter(api)
	// presentation.CompanyProfileRouter(api)
	// presentation.RequestPickupRouter(api)
	// presentation.PickupMatchingRouter(api)
	// presentation.PickupRatingRouter(api)

	// presentation.CollectorRouter(api)
	// presentation.TrashCartRouter(api)

	// presentation.UserProfileRouter(api)
	// presentation.UserPinRouter(api)
	// // presentation.RoleRouter(api)
	// presentation.WilayahRouter(api)
	// presentation.AddressRouter(api)
	// // presentation.ArticleRouter(api)
	// // presentation.AboutRouter(api)
	// presentation.TrashRouter(api)
	// presentation.CoverageAreaRouter(api)
	whatsapp.WhatsAppRouter(api)
}
