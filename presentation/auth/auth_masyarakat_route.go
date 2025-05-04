package presentation

import (
	"rijig/config"
	handler "rijig/internal/handler/auth"
	"rijig/internal/repositories"
	repository "rijig/internal/repositories/auth"
	services "rijig/internal/services/auth"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthMasyarakatRouter(api fiber.Router) {
	authMasyarakatRepo := repository.NewAuthMasyarakatRepositories(config.DB)
	roleRepo := repositories.NewRoleRepository(config.DB)
	authMasyarakatService := services.NewAuthMasyarakatService(authMasyarakatRepo, roleRepo)

	authHandler := handler.NewAuthMasyarakatHandler(authMasyarakatService)

	authMasyarakat := api.Group("/authmasyarakat")

	authMasyarakat.Post("/auth", authHandler.RegisterOrLoginHandler)
	authMasyarakat.Post("/logout", middleware.AuthMiddleware, authHandler.LogoutHandler)
	authMasyarakat.Post("/verify-otp", authHandler.VerifyOTPHandler)
}
