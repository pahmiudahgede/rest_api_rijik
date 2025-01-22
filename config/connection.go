package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/pahmiudahgede/senggoldong/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB         *gorm.DB
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	APIKey     string
	ServerHost string
	ServerPort string

	RedisClient   *redis.Client
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
)

func Context() context.Context {
	return context.Background()
}

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
		log.Fatal("Error: environment variables yang dibutuhkan tidak ada")
	}
}

func InitDatabase() {
	InitConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBName, DBPassword)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("gagal terhubung ke database: ", err)
	}

	err = DB.AutoMigrate(
		&domain.Point{},
		&domain.User{},
		&domain.UserRole{},
		&domain.UserPin{},
		&domain.MenuAccess{},
		&domain.PlatformHandle{},
		&domain.Address{},
		&domain.Article{},
		&domain.TrashCategory{},
		&domain.TrashDetail{},
		&domain.Banner{},
		&domain.CoverageArea{},
		&domain.CoverageDistric{},
		&domain.CoverageSubdistrict{},
		&domain.RequestPickup{},
		&domain.RequestItem{},
		&domain.Product{},
		&domain.ProductImage{},
		&domain.Store{},
	)
	if err != nil {
		log.Fatal("Error: Failed to auto migrate domain:", err)
	}

	fmt.Println("Koneksi ke database berhasil dan migrasi dilakukan")
}

func InitRedis() {
	InitConfig()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", RedisHost, RedisPort),
		Password: RedisPassword,
		DB:       RedisDB,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Gagal terhubung ke Redis:", err)
	}

	fmt.Println("Koneksi ke Redis berhasil")
}
