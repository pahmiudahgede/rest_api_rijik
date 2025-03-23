package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(api fiber.Router) {
	userRepo := repositories.NewUserRepository(config.DB)
	roleRepo := repositories.NewRoleRepository(config.DB)
	authService := services.NewAuthService(userRepo, roleRepo)

	authHandler := handler.NewAuthHandler(authService)

	api.Post("/auth", authHandler.RegisterOrLoginHandler)
	api.Post("/logout", middleware.AuthMiddleware, authHandler.LogoutHandler)
	api.Post("/verify-otp", authHandler.VerifyOTPHandler)
}
