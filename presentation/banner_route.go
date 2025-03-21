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

func BannerRouter(api fiber.Router) {
	bannerRepo := repositories.NewBannerRepository(config.DB)
	bannerService := services.NewBannerService(bannerRepo)
	BannerHandler := handler.NewBannerHandler(bannerService)

	bannerAPI := api.Group("/banner-rijik")

	bannerAPI.Post("/create-banner", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), BannerHandler.CreateBanner)
	bannerAPI.Get("/getall-banner", BannerHandler.GetAllBanners)
	bannerAPI.Get("/get-banner/:banner_id", BannerHandler.GetBannerByID)
	bannerAPI.Put("/update-banner/:banner_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), BannerHandler.UpdateBanner)
	bannerAPI.Delete("/delete-banner/:banner_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), BannerHandler.DeleteBanner)
}
