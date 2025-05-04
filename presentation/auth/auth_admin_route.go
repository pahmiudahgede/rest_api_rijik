package presentation

import (
	"log"
	"os"
	"rijig/config"
	handler "rijig/internal/handler/auth"
	"rijig/internal/repositories"
	repository "rijig/internal/repositories/auth"
	services "rijig/internal/services/auth"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func AdminAuthRouter(api fiber.Router) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is not set in the environment variables")
		os.Exit(1)
	}

	adminAuthRepo := repository.NewAuthAdminRepository(config.DB)
	roleRepo := repositories.NewRoleRepository(config.DB)

	adminAuthService := services.NewAuthAdminService(adminAuthRepo, roleRepo, secretKey)

	adminAuthHandler := handler.NewAuthAdminHandler(adminAuthService)

	adminAuthAPI := api.Group("/admin-auth")

	adminAuthAPI.Post("/register", adminAuthHandler.RegisterAdmin)
	adminAuthAPI.Post("/login", adminAuthHandler.LoginAdmin)
	adminAuthAPI.Post("/logout", middleware.AuthMiddleware, adminAuthHandler.LogoutAdmin)
}
