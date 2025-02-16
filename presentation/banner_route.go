package presentation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/internal/handler"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/middleware"
	"github.com/pahmiudahgede/senggoldong/utils"
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
