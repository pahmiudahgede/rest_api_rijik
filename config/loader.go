package config

import (
	"log"
	"os"
)

var (
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	ServerHost string
	ServerPort string
	APIKey     string
)

func InitConfig() {

	ServerHost = os.Getenv("SERVER_HOST")
	ServerPort = os.Getenv("SERVER_PORT")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	APIKey = os.Getenv("API_KEY")

	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	RedisDB = 0

	if ServerHost == "" || ServerPort == "" || DBHost == "" || DBPort == "" || DBName == "" || DBUser == "" || DBPassword == "" || APIKey == "" {
		log.Fatal("Error: Beberapa environment variables tidak ditemukan.")
	}
}
