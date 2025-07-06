package company

import (
	"rijig/config"
	"rijig/internal/authentication"
	"rijig/internal/userprofile"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func CompanyRouter(api fiber.Router) {
	companyProfileRepo := NewCompanyProfileRepository(config.DB)
	authRepo := authentication.NewAuthenticationRepository(config.DB)
	userRepo := userprofile.NewUserProfileRepository(config.DB)
	companyProfileService := NewCompanyProfileService(companyProfileRepo, authRepo, userRepo)
	companyProfileHandler := NewCompanyProfileHandler(companyProfileService)

	companyProfileAPI := api.Group("/companyprofile")
	companyProfileAPI.Use(middleware.AuthMiddleware())

	companyProfileAPI.Post("/create", companyProfileHandler.CreateCompanyProfile)
	companyProfileAPI.Get("/get/:id", companyProfileHandler.GetCompanyProfileByID)
	companyProfileAPI.Get("/get", companyProfileHandler.GetCompanyProfilesByUserID)
	companyProfileAPI.Put("/update", companyProfileHandler.UpdateCompanyProfile)
	companyProfileAPI.Delete("/delete", companyProfileHandler.DeleteCompanyProfile)
}
