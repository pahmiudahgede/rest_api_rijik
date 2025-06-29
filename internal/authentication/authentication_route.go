package authentication

import (
	"rijig/config"
	"rijig/internal/role"
	"rijig/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthenticationRouter(api fiber.Router) {
	repoAuth := NewAuthenticationRepository(config.DB)
	repoRole := role.NewRoleRepository(config.DB)

	authService := NewAuthenticationService(repoAuth, repoRole)
	authHandler := NewAuthenticationHandler(authService)

	authRoute := api.Group("/auth")

	authRoute.Post("/refresh-token",
		middleware.AuthMiddleware(),
		middleware.DeviceValidation(),
		authHandler.RefreshToken,
	)

	// authRoute.Get("/me",
	// 	middleware.AuthMiddleware(),
	// 	middleware.CheckRefreshTokenTTL(30*time.Second),
	// 	middleware.RequireApprovedRegistration(),
	// 	authHandler.GetMe,
	// )

	authRoute.Get("/cekapproval", middleware.AuthMiddleware(), authHandler.GetRegistrationStatus)
	authRoute.Post("/login/admin", authHandler.Login)
	authRoute.Post("/register/admin", authHandler.RegisterAdmin)

	authRoute.Post("/verify-email", authHandler.VerifyEmail)
	authRoute.Post("/resend-verification", authHandler.ResendEmailVerification)

	authRoute.Post("/verify-otp-admin", authHandler.VerifyAdminOTP)
	authRoute.Post("/resend-otp-admin", authHandler.ResendAdminOTP)
	
	authRoute.Post("/forgot-password", authHandler.ForgotPassword)
	authRoute.Post("/reset-password", authHandler.ResetPassword)

	authRoute.Post("/request-otp", authHandler.RequestOtpHandler)
	authRoute.Post("/verif-otp", authHandler.VerifyOtpHandler)
	authRoute.Post("/request-otp/register", authHandler.RequestOtpRegistHandler)
	authRoute.Post("/verif-otp/register", authHandler.VerifyOtpRegistHandler)
	authRoute.Post("/logout", middleware.AuthMiddleware(), authHandler.LogoutAuthentication)
}
