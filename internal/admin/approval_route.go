package admin

import (
	"rijig/config"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func ApprovalRoutes(api fiber.Router) {
	repo := NewAdminRepository(config.DB)
	service := NewAdminService(repo)
	handler := NewAdminHandler(service)

	admin := api.Group("/admusers")
	admin.Use(middleware.RequireAdminRole(), middleware.AuthMiddleware())

	admin.Get("/getalluser", handler.GetAllUsers)
	admin.Patch("/reguser/:userid", handler.UpdateRegistrationStatus)

	admin.Get("/statistics", handler.GetUserStatistics)
	admin.Get("/export", handler.GetAllUsersExport)
	admin.Get("/role/:role", handler.GetUsersByRole)
	admin.Get("/health", handler.HealthCheck)
}
