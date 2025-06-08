package address

import (
	"rijig/config"
	"rijig/internal/wilayahindo"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func AddressRouter(api fiber.Router) {
	addressRepo := NewAddressRepository(config.DB)
	wilayahRepo := wilayahindo.NewWilayahIndonesiaRepository(config.DB)
	addressService := NewAddressService(addressRepo, wilayahRepo)
	addressHandler := NewAddressHandler(addressService)

	adddressAPI := api.Group("/user/address")

	adddressAPI.Post("/create-address", middleware.AuthMiddleware(), addressHandler.CreateAddress)
	adddressAPI.Get("/get-address", middleware.AuthMiddleware(), addressHandler.GetAddressByUserID)
	adddressAPI.Get("/get-address/:address_id", middleware.AuthMiddleware(), addressHandler.GetAddressByID)
	adddressAPI.Put("/update-address/:address_id", middleware.AuthMiddleware(), addressHandler.UpdateAddress)
	adddressAPI.Delete("/delete-address/:address_id", middleware.AuthMiddleware(), addressHandler.DeleteAddress)
}
