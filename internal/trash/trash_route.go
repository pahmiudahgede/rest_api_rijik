// ===internal/trash/trash_route.go===
package trash

import (
	"rijig/config"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func TrashRouter(api fiber.Router) {
	trashRepo := NewTrashRepository(config.DB)
	trashService := NewTrashService(trashRepo)
	trashHandler := NewTrashHandler(trashService)

	trashAPI := api.Group("/trash")
	trashAPI.Use(middleware.AuthMiddleware())

	// ============= TRASH CATEGORY ROUTES =============

	// Create trash category (JSON)
	trashAPI.Post("/category", trashHandler.CreateTrashCategory)

	// Create trash category with icon (form-data)
	trashAPI.Post("/category/with-icon", trashHandler.CreateTrashCategoryWithIcon)

	// Create trash category with details (JSON)
	trashAPI.Post("/category/with-details", trashHandler.CreateTrashCategoryWithDetails)

	// Get all trash categories (with optional query param: ?with_details=true)
	trashAPI.Get("/category", trashHandler.GetAllTrashCategories)

	// Get trash category by ID (with optional query param: ?with_details=true)
	trashAPI.Get("/category/:id", trashHandler.GetTrashCategoryByID)

	// Update trash category (JSON)
	trashAPI.Put("/category/:id", trashHandler.UpdateTrashCategory)

	// Update trash category with icon (form-data)
	trashAPI.Put("/category/:id/with-icon", trashHandler.UpdateTrashCategoryWithIcon)

	// Delete trash category
	trashAPI.Delete("/category/:id", trashHandler.DeleteTrashCategory)

	// ============= TRASH DETAIL ROUTES =============

	// Create trash detail (JSON)
	trashAPI.Post("/detail", trashHandler.CreateTrashDetail)

	// Create trash detail with icon (form-data)
	trashAPI.Post("/detail/with-icon", trashHandler.CreateTrashDetailWithIcon)

	// Add trash detail to specific category (JSON)
	trashAPI.Post("/category/:categoryId/detail", trashHandler.AddTrashDetailToCategory)

	// Add trash detail to specific category with icon (form-data)
	trashAPI.Post("/category/:categoryId/detail/with-icon", trashHandler.AddTrashDetailToCategoryWithIcon)

	// Get trash details by category ID
	trashAPI.Get("/category/:categoryId/details", trashHandler.GetTrashDetailsByCategory)

	// Get trash detail by ID
	trashAPI.Get("/detail/:id", trashHandler.GetTrashDetailByID)

	// Update trash detail (JSON)
	trashAPI.Put("/detail/:id", trashHandler.UpdateTrashDetail)

	// Update trash detail with icon (form-data)
	trashAPI.Put("/detail/:id/with-icon", trashHandler.UpdateTrashDetailWithIcon)

	// Delete trash detail
	trashAPI.Delete("/detail/:id", trashHandler.DeleteTrashDetail)

	// ============= BULK OPERATIONS ROUTES =============

	// Bulk create trash details for specific category
	trashAPI.Post("/category/:categoryId/details/bulk", trashHandler.BulkCreateTrashDetails)

	// Bulk delete trash details
	trashAPI.Delete("/details/bulk", trashHandler.BulkDeleteTrashDetails)

	// Reorder trash details within a category
	trashAPI.Put("/category/:categoryId/details/reorder", trashHandler.ReorderTrashDetails)
}
