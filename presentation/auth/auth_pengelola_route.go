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

func AuthPengelolaRouter(api fiber.Router) {
	authPengelolaRepo := repository.NewAuthPengelolaRepositories(config.DB)
	roleRepo := repositories.NewRoleRepository(config.DB)
	authPengelolaService := services.NewAuthPengelolaService(authPengelolaRepo, roleRepo)

	authHandler := handler.NewAuthPengelolaHandler(authPengelolaService)

	authPengelola := api.Group("/authpengelola")

	authPengelola.Post("/auth", authHandler.RegisterOrLoginHandler)
	authPengelola.Post("/logout", middleware.AuthMiddleware, authHandler.LogoutHandler)
	authPengelola.Post("/verify-otp", authHandler.VerifyOTPHandler)
}
