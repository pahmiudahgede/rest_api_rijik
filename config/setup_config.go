package config

import (
	"log"

	"github.com/joho/godotenv"
)

func SetupConfig() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ConnectDatabase()
	ConnectRedis()
	InitWhatsApp() 
}
