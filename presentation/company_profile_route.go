package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func CompanyProfileRouter(api fiber.Router) {

	companyProfileRepo := repositories.NewCompanyProfileRepository(config.DB)
	companyProfileService := services.NewCompanyProfileService(companyProfileRepo)
	companyProfileHandler := handler.NewCompanyProfileHandler(companyProfileService)

	companyProfileAPI := api.Group("/company-profile")
	companyProfileAPI.Use(middleware.AuthMiddleware())	

	companyProfileAPI.Post("/create", companyProfileHandler.CreateCompanyProfile)
	companyProfileAPI.Get("/get/:company_id", companyProfileHandler.GetCompanyProfileByID)
	companyProfileAPI.Get("/get", companyProfileHandler.GetCompanyProfilesByUserID)
	companyProfileAPI.Put("/update/:company_id", companyProfileHandler.UpdateCompanyProfile)
	companyProfileAPI.Delete("/delete/:company_id", companyProfileHandler.DeleteCompanyProfile)
}
