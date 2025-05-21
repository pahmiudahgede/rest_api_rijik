package main

import (
	"log"
	"os"
	"rijig/config"
	"rijig/internal/worker"
	"rijig/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/robfig/cron"
)

func main() {
	config.SetupConfig()
	logFile, _ := os.OpenFile("logs/cart_commit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log.SetOutput(logFile)

	go func() {
		c := cron.New()
		c.AddFunc("@every 1m", func() {
			_ = worker.CommitExpiredCartsToDB()
		})
		c.Start()
	}()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Content-Type,x-api-key",
	}))

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     "http://localhost:3000",
	// 	AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	// 	AllowHeaders:     "Origin, Content-Type, Accept, Authorization, x-api-key",
	// 	AllowCredentials: true,
	// }))

	// app.Use(func(c *fiber.Ctx) error {
	// 	c.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// 	c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	// 	c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, x-api-key")
	// 	c.Set("Access-Control-Allow-Credentials", "true")
	// 	return c.Next()
	// })

	// app.Options("*", func(c *fiber.Ctx) error {
	// 	c.Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// 	c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	// 	c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, x-api-key")
	// 	c.Set("Access-Control-Allow-Credentials", "true")
	// 	return c.SendStatus(fiber.StatusNoContent)
	// })

	router.SetupRoutes(app)
	config.StartServer(app)
}
