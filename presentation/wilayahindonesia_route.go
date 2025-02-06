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

func WilayahRouter(api fiber.Router) {

	wilayahRepo := repositories.NewWilayahIndonesiaRepository(config.DB)
	wilayahService := services.NewWilayahIndonesiaService(wilayahRepo)
	wilayahHandler := handler.NewWilayahImportHandler(wilayahService)

	api.Post("/import/data-wilayah-indonesia", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator), wilayahHandler.ImportWilayahData)

	wilayahAPI := api.Group("/wilayah-indonesia")

	wilayahAPI.Get("/provinces", middleware.AuthMiddleware, wilayahHandler.GetProvinces)
	wilayahAPI.Get("/provinces/:provinceid", middleware.AuthMiddleware, wilayahHandler.GetProvinceByID)

	wilayahAPI.Get("/regencies", middleware.AuthMiddleware, wilayahHandler.GetAllRegencies)
	wilayahAPI.Get("/regencies/:regencyid", middleware.AuthMiddleware, wilayahHandler.GetRegencyByID)

	wilayahAPI.Get("/districts", middleware.AuthMiddleware, wilayahHandler.GetAllDistricts)
	wilayahAPI.Get("/districts/:districtid", middleware.AuthMiddleware, wilayahHandler.GetDistrictByID)

	wilayahAPI.Get("/villages", middleware.AuthMiddleware, wilayahHandler.GetAllVillages)
	wilayahAPI.Get("/villages/:villageid", middleware.AuthMiddleware, wilayahHandler.GetVillageByID)

}
