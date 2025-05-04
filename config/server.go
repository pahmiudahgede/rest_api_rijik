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
	port := os.Getenv("SERVER_PORT")

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("server berjalan di http://%s", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("gagal saat menjalankan server: %v", err)
	}
}
