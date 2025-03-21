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

func WilayahRouter(api fiber.Router) {

	wilayahRepo := repositories.NewWilayahIndonesiaRepository(config.DB)
	wilayahService := services.NewWilayahIndonesiaService(wilayahRepo)
	wilayahHandler := handler.NewWilayahImportHandler(wilayahService)

	api.Post("/import/data-wilayah-indonesia", middleware.RoleMiddleware(utils.RoleAdministrator), wilayahHandler.ImportWilayahData)

	wilayahAPI := api.Group("/wilayah-indonesia")

	wilayahAPI.Get("/provinces", wilayahHandler.GetProvinces)
	wilayahAPI.Get("/provinces/:provinceid", wilayahHandler.GetProvinceByID)

	wilayahAPI.Get("/regencies", wilayahHandler.GetAllRegencies)
	wilayahAPI.Get("/regencies/:regencyid", wilayahHandler.GetRegencyByID)

	wilayahAPI.Get("/districts", wilayahHandler.GetAllDistricts)
	wilayahAPI.Get("/districts/:districtid", wilayahHandler.GetDistrictByID)

	wilayahAPI.Get("/villages", wilayahHandler.GetAllVillages)
	wilayahAPI.Get("/villages/:villageid", wilayahHandler.GetVillageByID)

}
