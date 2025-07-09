package admin

import (
	"rijig/config"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func ApprovalRoutes(api fiber.Router) {
	baseRepo := NewApprovalRepository(config.DB)
	baseService := NewApprovalService(baseRepo)
	baseHandler := NewApprovalHandler(baseService)

	adminGroup := api.Group("/needapprove")
	adminGroup.Use(middleware.RequireAdminRole(), middleware.AuthMiddleware())

	adminGroup.Get("/pending", baseHandler.GetPendingUsers)

	adminGroup.Get("/:user_id/approval-details", baseHandler.GetUserApprovalDetails)
	adminGroup.Post("/approval-action", baseHandler.ProcessApprovalAction)
	adminGroup.Post("/bulk-approval", baseHandler.BulkProcessApproval)
	adminGroup.Post("/:user_id/approve", baseHandler.ApproveUser)
	adminGroup.Post("/:user_id/reject", baseHandler.RejectUser)
}
