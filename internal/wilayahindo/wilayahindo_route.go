package wilayahindo

import (
	"rijig/config"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func WilayahRouter(api fiber.Router) {
	wilayahRepo := NewWilayahIndonesiaRepository(config.DB)
	wilayahService := NewWilayahIndonesiaService(wilayahRepo)
	wilayahHandler := NewWilayahIndonesiaHandler(wilayahService)

	api.Post("/import/data-wilayah-indonesia", middleware.RequireAdminRole(), wilayahHandler.ImportDataFromCSV)

	wilayahAPI := api.Group("/wilayah-indonesia")

	wilayahAPI.Get("/provinces", wilayahHandler.GetAllProvinces)
	wilayahAPI.Get("/provinces/:id", wilayahHandler.GetProvinceByID)

	wilayahAPI.Get("/regencies", wilayahHandler.GetAllRegencies)
	wilayahAPI.Get("/regencies/:id", wilayahHandler.GetRegencyByID)

	wilayahAPI.Get("/districts", wilayahHandler.GetAllDistricts)
	wilayahAPI.Get("/districts/:id", wilayahHandler.GetDistrictByID)

	wilayahAPI.Get("/villages", wilayahHandler.GetAllVillages)
	wilayahAPI.Get("/villages/:id", wilayahHandler.GetVillageByID)
}
