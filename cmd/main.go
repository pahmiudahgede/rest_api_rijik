package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/router"
)

func main() {
	config.SetupConfig()
	
	app := fiber.New()
	// app.Static(utils.BaseUrl+"/uploads", "./public"+utils.BaseUrl+"/uploads")
	
	router.SetupRoutes(app)
	config.StartServer(app)

}
