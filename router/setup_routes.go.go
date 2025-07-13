package router

import (
	"os"

	"rijig/internal/about"
	"rijig/internal/address"
	"rijig/internal/admin"
	"rijig/internal/article"
	"rijig/internal/authentication"
	"rijig/internal/cart"
	"rijig/internal/company"
	"rijig/internal/identitycart"
	"rijig/internal/requestpickup"
	"rijig/internal/role"
	"rijig/internal/trash"
	"rijig/internal/userpin"
	"rijig/internal/userprofile"
	"rijig/internal/whatsapp"
	"rijig/internal/wilayahindo"
	"rijig/middleware"

	// "rijig/presentation"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	apa := app.Group(os.Getenv("BASE_URL"))
	apa.Static("/uploads", "./public"+os.Getenv("BASE_URL")+"/uploads")
	// a := app.Group(os.Getenv("BASE_URL"))
	// whatsapp.WhatsAppRouter(a)

	api := app.Group(os.Getenv("BASE_URL"))
	api.Use(middleware.APIKeyMiddleware)

	authentication.AuthenticationRouter(api)
	identitycart.UserIdentityCardRoute(api)
	company.CompanyRouter(api)
	userpin.UsersPinRoute(api)
	role.UserRoleRouter(api)

	article.ArticleRouter(api)
	userprofile.UserProfileRouter(api)
	wilayahindo.WilayahRouter(api)
	trash.TrashRouter(api)
	about.AboutRouter(api)
	whatsapp.WhatsAppRouter(api)
	admin.ApprovalRoutes(api)

	// || auth router || //
	// || auth router || //
	// presentation.IdentityCardRouter(api)
	// presentation.CompanyProfileRouter(api)
	requestpickup.RequestPickupRouter(api)
	// presentation.PickupMatchingRouter(api)
	// presentation.PickupRatingRouter(api)

	// presentation.CollectorRouter(api)
	cart.TrashCartRouter(api)

	// presentation.UserProfileRouter(api)
	// presentation.UserPinRouter(api)
	// // presentation.RoleRouter(api)
	// presentation.WilayahRouter(api)
	address.AddressRouter(api)
	// // presentation.ArticleRouter(api)
	// // presentation.AboutRouter(api)
	// presentation.TrashRouter(api)
	// presentation.CoverageAreaRouter(api)
}
