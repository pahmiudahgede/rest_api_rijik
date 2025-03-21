package presentation

import (
	"rijig/config"
	"rijig/internal/handler"
	"rijig/internal/repositories"
	"rijig/internal/services"

	"github.com/gofiber/fiber/v2"
	// "gorm.io/gorm"
	// "rijig/middleware"
)

func AuthRouter(api fiber.Router) {
	// userRepo := repositories.NewUserRepository(config.DB)
	// roleRepo := repositories.NewRoleRepository(config.DB)
	// userService := services.NewUserService(userRepo, roleRepo, secretKey)
	// userHandler := handler.NewUserHandler(userService)

	// api.Post("/login", userHandler.Login)
	// api.Post("/register", userHandler.Register)
	// api.Post("/logout", middleware.AuthMiddleware, userHandler.Logout)
	// userRepo := repositories.NewUserRepository(config.DB)
	// authService := services.NewAuthService(userRepo, secretKey)

	// // Inisialisasi handler
	// authHandler := handler.NewAuthHandler(authService)

	// // Endpoint OTP
	// authRoutes := api.Group("/auth")
	// authRoutes.Post("/send-otp", authHandler.SendOTP)
	// authRoutes.Post("/verify-otp", authHandler.VerifyOTP)
	// userRepo := repositories.NewUserRepository(config.DB)
	// authService := services.NewAuthService(userRepo)

	// authHandler := handler.NewAuthHandler(authService)

	// // Routes
	// api.Post("/register", authHandler.Register)
	// api.Post("/verify-otp", authHandler.VerifyOTP)
	userRepo := repositories.NewUserRepository(config.DB)
	roleRepo := repositories.NewRoleRepository(config.DB)
	redisRepo := repositories.NewRedisRepository(config.RedisClient)

	// Setup Service
	authService := services.NewAuthService(userRepo, roleRepo, redisRepo)

	// Setup Handler
	authHandler := handler.NewAuthHandler(authService)

	// Define Routes
	api.Post("/register", authHandler.Register) // Route untuk registrasi
	api.Post("/verify-otp", authHandler.VerifyOTP)
}
