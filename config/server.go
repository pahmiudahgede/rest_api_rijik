package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetSecretKey() string {
	return os.Getenv("SECRET_KEY")
}

func StartServer(app *fiber.App) {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "0.0.0.0" // Default untuk Railway
	}

	// Railway menggunakan PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		// Fallback ke SERVER_PORT untuk development lokal
		port = os.Getenv("SERVER_PORT")
		if port == "" {
			port = "7000" // Default fallback
		}
	}

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Server starting on %s", address)
	log.Printf("Environment: %s", getEnvironment())
	
	if err := app.Listen(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnvironment() string {
	if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
		return "Railway Production"
	}
	if os.Getenv("DOCKER_ENV") != "" {
		return "Docker Development"
	}
	return "Local Development"
}