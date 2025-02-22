package presentation

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/internal/handler"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/middleware"
)

func AuthRouter(api fiber.Router) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is not set in the environment variables")
		os.Exit(1)
	}

	userRepo := repositories.NewUserRepository(config.DB)
	roleRepo := repositories.NewRoleRepository(config.DB)
	userService := services.NewUserService(userRepo, roleRepo, secretKey)
	userHandler := handler.NewUserHandler(userService)

	api.Post("/login", userHandler.Login)
	api.Post("/register", userHandler.Register)
	api.Post("/logout", middleware.AuthMiddleware, userHandler.Logout)

}
