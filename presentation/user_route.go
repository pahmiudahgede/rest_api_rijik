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
	userProfileRepo := repositories.NewUserProfilRepository(config.DB)
	userProfileService := services.NewUserService(userProfileRepo)
	userProfileHandler := handler.NewUserHandler(userProfileService)

	userProfilRoute := api.Group("/user")

	userProfilRoute.Get("/info", middleware.AuthMiddleware(), userProfileHandler.GetUserByIDHandler)

	userProfilRoute.Get("/show-all", middleware.AuthMiddleware(), userProfileHandler.GetAllUsersHandler)
	// userProfilRoute.Get("/:userid", middleware.AuthMiddleware, userProfileHandler.GetUserProfileById)
	// userProfilRoute.Get("/:roleid", middleware.AuthMiddleware, userProfileHandler.GetUsersByRoleID)

	userProfilRoute.Put("/update-user", middleware.AuthMiddleware(), userProfileHandler.UpdateUserHandler)
	userProfilRoute.Patch("/update-user-password", middleware.AuthMiddleware(), userProfileHandler.UpdateUserPasswordHandler)
	userProfilRoute.Patch("/upload-photoprofile", middleware.AuthMiddleware(), userProfileHandler.UpdateUserAvatarHandler)
}
