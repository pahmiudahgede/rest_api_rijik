package identitycart

import (
	"rijig/config"
	"rijig/internal/authentication"
	"rijig/internal/userprofile"
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func UserIdentityCardRoute(api fiber.Router) {
	identityRepo := NewIdentityCardRepository(config.DB)
	authRepo := authentication.NewAuthenticationRepository(config.DB)
	userRepo := userprofile.NewUserProfileRepository(config.DB)
	identityService := NewIdentityCardService(identityRepo, authRepo, userRepo)
	identityHandler := NewIdentityCardHandler(identityService)

	identity := api.Group("/identity")

	identity.Post("/create",
		middleware.AuthMiddleware(),
		middleware.RequireRoles(utils.RolePengepul),
		identityHandler.CreateIdentityCardHandler,
	)
	identity.Get("/:id",
		middleware.AuthMiddleware(),
		identityHandler.GetIdentityByID,
	)
	identity.Get("/s",
		middleware.AuthMiddleware(),
		identityHandler.GetIdentityByUserId,
	)
	identity.Get("/",
		middleware.AuthMiddleware(),
		middleware.RequireRoles(utils.RoleAdministrator),
		identityHandler.GetAllIdentityCardsByRegStatus,
	)
	identity.Patch("/:userId/status",
		middleware.AuthMiddleware(),
		middleware.RequireRoles(utils.RoleAdministrator),
		identityHandler.UpdateUserRegistrationStatusByIdentityCard,
	)

}
