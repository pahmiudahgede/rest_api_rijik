package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func SetupConfig() {

	if _, exists := os.LookupEnv("DOCKER_ENV"); exists {

		log.Println("Running in Docker container, using environment variables")
	} else {

		err := godotenv.Load(".env.dev")
		if err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
			log.Println("Trying to use system environment variables...")
		} else {
			log.Println("Loaded environment from .env.dev file")
		}
	}
	ConnectDatabase()
	ConnectRedis()
	go func() {
		InitWhatsApp() // Ini tidak akan blocking startup server
	}()
}
