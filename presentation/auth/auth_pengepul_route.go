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

func AuthPengepulRouter(api fiber.Router) {
	authPengepulRepo := repository.NewAuthPengepulRepositories(config.DB)
	roleRepo := repositories.NewRoleRepository(config.DB)
	authPengepulService := services.NewAuthPengepulService(authPengepulRepo, roleRepo)

	authHandler := handler.NewAuthPengepulHandler(authPengepulService)

	authPengepul := api.Group("/authpengepul")

	authPengepul.Post("/auth", authHandler.RegisterOrLoginHandler)
	authPengepul.Post("/logout", middleware.AuthMiddleware, authHandler.LogoutHandler)
	authPengepul.Post("/verify-otp", authHandler.VerifyOTPHandler)
}
