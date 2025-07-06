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
	wilayahAPI.Get("/provinces/:provinceid", wilayahHandler.GetProvinceByID)

	wilayahAPI.Get("/regencies", wilayahHandler.GetAllRegencies)
	wilayahAPI.Get("/regencies/:regencyid", wilayahHandler.GetRegencyByID)

	wilayahAPI.Get("/districts", wilayahHandler.GetAllDistricts)
	wilayahAPI.Get("/districts/:districtid", wilayahHandler.GetDistrictByID)

	wilayahAPI.Get("/villages", wilayahHandler.GetAllVillages)
	wilayahAPI.Get("/villages/:villageid", wilayahHandler.GetVillageByID)

}
