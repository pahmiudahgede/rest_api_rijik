package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func IdentityCardRouter(api fiber.Router) {
	identityCardRepo := repositories.NewIdentityCardRepository(config.DB)
	userRepo := repositories.NewUserProfilRepository(config.DB)
	identityCardService := services.NewIdentityCardService(identityCardRepo, userRepo)
	identityCardHandler := handler.NewIdentityCardHandler(identityCardService)

	identityCardApi := api.Group("/identitycard")
	identityCardApi.Use(middleware.AuthMiddleware)

	identityCardApi.Post("/create", middleware.RoleMiddleware(utils.RolePengelola, utils.RolePengepul), identityCardHandler.CreateIdentityCard)
	identityCardApi.Get("/get/:identity_id", identityCardHandler.GetIdentityCardById)
	identityCardApi.Get("/get", identityCardHandler.GetIdentityCard)
	identityCardApi.Put("/update/:identity_id", middleware.RoleMiddleware(utils.RolePengelola, utils.RolePengepul), identityCardHandler.UpdateIdentityCard)
	identityCardApi.Delete("/delete/:identity_id", middleware.RoleMiddleware(utils.RolePengelola, utils.RolePengepul), identityCardHandler.DeleteIdentityCard)
}
