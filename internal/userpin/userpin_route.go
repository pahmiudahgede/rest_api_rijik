package userpin

import (
	"rijig/config"
	"rijig/internal/authentication"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func UsersPinRoute(api fiber.Router) {
	userPinRepo := NewUserPinRepository(config.DB)
	authRepo := authentication.NewAuthenticationRepository(config.DB)

	userPinService := NewUserPinService(userPinRepo, authRepo)

	userPinHandler := NewUserPinHandler(userPinService)

	pin := api.Group("/pin", middleware.AuthMiddleware())

	pin.Post("/create", userPinHandler.CreateUserPinHandler)
	pin.Post("/verif", userPinHandler.VerifyPinHandler)
}
