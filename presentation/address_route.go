package presentation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/internal/handler"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/middleware"
)

func AddressRouter(api fiber.Router) {
	addressRepo := repositories.NewAddressRepository(config.DB)
	wilayahRepo := repositories.NewWilayahIndonesiaRepository(config.DB)
	addressService := services.NewAddressService(addressRepo, wilayahRepo)
	addressHandler := handler.NewAddressHandler(addressService)

	adddressAPI := api.Group("/user/address")

	adddressAPI.Post("/create-address", middleware.AuthMiddleware, addressHandler.CreateAddress)
	adddressAPI.Get("/get-address", middleware.AuthMiddleware, addressHandler.GetAddressByUserID)
	adddressAPI.Get("/get-address/:address_id", middleware.AuthMiddleware, addressHandler.GetAddressByID)
}
