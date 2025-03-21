package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserPinRouter(api fiber.Router) {

	userPinRepo := repositories.NewUserPinRepository(config.DB)

	userPinService := services.NewUserPinService(userPinRepo)

	userPinHandler := handler.NewUserPinHandler(userPinService)

	api.Post("/user/set-pin", middleware.AuthMiddleware, userPinHandler.CreateUserPin)
	api.Post("/user/verif-pin", middleware.AuthMiddleware, userPinHandler.VerifyUserPin)
	api.Get("/user/cek-pin-status", middleware.AuthMiddleware, userPinHandler.CheckPinStatus)
	api.Patch("/user/update-pin", middleware.AuthMiddleware, userPinHandler.UpdateUserPin)
}
