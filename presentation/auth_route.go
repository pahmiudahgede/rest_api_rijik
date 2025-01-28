package presentation

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/internal/handler"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
)

func AuthRouter(app *fiber.App) {
	api := app.Group("/apirijikid")

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		panic("SECRET_KEY is not set in the environment variables")
	}

	userRepo := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepo, secretKey)
	userHandler := handler.NewUserHandler(userService)

	api.Post("/login", userHandler.Login)
}
