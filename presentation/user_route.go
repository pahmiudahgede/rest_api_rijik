package presentation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/internal/handler"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/middleware"
)

func UserProfileRouter(api fiber.Router) {
	userProfileRepo := repositories.NewUserProfileRepository(config.DB)
	userProfileService := services.NewUserProfileService(userProfileRepo)
	userProfileHandler := handler.NewUserProfileHandler(userProfileService)

	api.Get("/user", middleware.AuthMiddleware, userProfileHandler.GetUserProfile)
	api.Put("/user/update-user", middleware.AuthMiddleware, userProfileHandler.UpdateUserProfile)
	api.Post("/user/update-user-password", middleware.AuthMiddleware, userProfileHandler.UpdateUserPassword)
}
