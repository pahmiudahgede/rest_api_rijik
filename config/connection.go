package config

import (
	"fmt"
	"log"
	"os"

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
	)
	if err != nil {
		log.Fatal("Error: Failed to auto migrate domain:", err)
	}

	fmt.Println("Koneksi ke database berhasil dan migrasi dilakukan")
}
