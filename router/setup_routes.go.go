package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/presentation"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/apirijikid")

	presentation.AuthRouter(api)
}
