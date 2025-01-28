package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/middleware"
	"github.com/pahmiudahgede/senggoldong/presentation"
)

func main() {
	config.SetupConfig()

	app := fiber.New()
	app.Use(middleware.APIKeyMiddleware)
	presentation.AuthRouter(app)
	config.StartServer(app)
}
