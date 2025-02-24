package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/router"
)

func main() {
	config.SetupConfig()

	app := fiber.New()

	router.SetupRoutes(app)

	config.StartServer(app)
}
