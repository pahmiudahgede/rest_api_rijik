package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserProfileRouter(api fiber.Router) {
	userProfileRepo := repositories.NewUserProfileRepository(config.DB)
	userProfileService := services.NewUserProfileService(userProfileRepo)
	userProfileHandler := handler.NewUserProfileHandler(userProfileService)

	userProfilRoute := api.Group("/user")

	userProfilRoute.Get("/info", middleware.AuthMiddleware, userProfileHandler.GetUserProfile)

	userProfilRoute.Get("/show-all", middleware.AuthMiddleware, userProfileHandler.GetAllUsers)
	userProfilRoute.Get("/:userid", middleware.AuthMiddleware, userProfileHandler.GetUserProfileById)
	userProfilRoute.Get("/:roleid", middleware.AuthMiddleware, userProfileHandler.GetUsersByRoleID)

	userProfilRoute.Put("/update-user", middleware.AuthMiddleware, userProfileHandler.UpdateUserProfile)
	// userProfilRoute.Patch("/update-user-password", middleware.AuthMiddleware, userProfileHandler.UpdateUserPassword)
	userProfilRoute.Patch("/upload-photoprofile", middleware.AuthMiddleware, userProfileHandler.UpdateUserAvatar)
}
