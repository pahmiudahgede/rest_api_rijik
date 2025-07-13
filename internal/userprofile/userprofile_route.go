package userprofile

import (
	"rijig/config"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserProfileRouter(api fiber.Router) {
	userProfileRepo := NewUserProfileRepository(config.DB)
	userProfileService := NewUserProfileService(userProfileRepo)
	userProfileHandler := NewUserProfileHandler(userProfileService)

	userRoute := api.Group("/userprofile")
	userRoute.Use(middleware.AuthMiddleware())

	userRoute.Get("/", userProfileHandler.GetUserProfile)
	userRoute.Put("/update", userProfileHandler.UpdateUserProfile)
}
