package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/internal/api"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error saat memuat file .env")
	}

	config.InitConfig()
	config.InitDatabase()

}

func main() {
	app := fiber.New()

	api.AppRouter(app)

	log.Fatal(app.Listen(":" + os.Getenv("SERVER_PORT")))
}
