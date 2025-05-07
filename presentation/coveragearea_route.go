package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"

	"github.com/gofiber/fiber/v2"
)

func CoverageAreaRouter(api fiber.Router) {
	coverageAreaRepo := repositories.NewCoverageAreaRepository(config.DB)
	wilayahRepo := repositories.NewWilayahIndonesiaRepository(config.DB)
	coverageAreaService := services.NewCoverageAreaService(coverageAreaRepo, wilayahRepo)
	coverageAreaHandler := handler.NewCoverageAreaHandler(coverageAreaService)

	coverage := api.Group("/coveragearea")

	coverage.Post("/", coverageAreaHandler.CreateCoverageArea)
	coverage.Get("/", coverageAreaHandler.GetAllCoverageAreas)
	coverage.Get("/:id", coverageAreaHandler.GetCoverageAreaByID)
	coverage.Put("/:id", coverageAreaHandler.UpdateCoverageArea)
	coverage.Delete("/:id", coverageAreaHandler.DeleteCoverageArea)
}
