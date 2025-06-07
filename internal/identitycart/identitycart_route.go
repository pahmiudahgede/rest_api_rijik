package identitycart

import (
	"rijig/config"
	"rijig/internal/authentication"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserIdentityCardRoute(api fiber.Router) {
	identityRepo := NewIdentityCardRepository(config.DB)
	authRepo := authentication.NewAuthenticationRepository(config.DB)
	identityService := NewIdentityCardService(identityRepo, authRepo)
	identityHandler := NewIdentityCardHandler(identityService)

	identity := api.Group("/identity")

	identity.Post("/create",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("pengelola", "pengepul"),
		identityHandler.CreateIdentityCardHandler,
	)
	identity.Get("/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("pengelola", "pengepul"),
		identityHandler.GetIdentityByID,
	)
	identity.Get("/",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("pengelola", "pengepul"),
		identityHandler.GetIdentityByUserId,
	)

}
