package presentation

import (
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func PickupRatingRouter(api fiber.Router) {
	ratingRepo := repositories.NewPickupRatingRepository()
	ratingService := services.NewPickupRatingService(ratingRepo)
	ratingHandler := handler.NewPickupRatingHandler(ratingService)

	rating := api.Group("/pickup")
	rating.Use(middleware.AuthMiddleware())
	rating.Post("/:id/rating", ratingHandler.CreateRating)

	collector := api.Group("/collector")
	collector.Get("/:id/ratings", ratingHandler.GetRatingsByCollector)
	collector.Get("/:id/ratings/average", ratingHandler.GetAverageRating)
}
