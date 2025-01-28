package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func StartServer(app *fiber.App) {
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Server is running on http://%s", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
