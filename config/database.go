package config

import (
	"fmt"
	"log"
	"os"

	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Database connected successfully!")

	err = DB.AutoMigrate(
		// ==wilayah indonesia==
		&model.Province{},
		&model.Regency{},
		&model.District{},
		&model.Village{},
		// ==wilayah indonesia==

		// ==main feature==
		&model.User{},
		&model.Role{},
		&model.UserPin{},
		&model.Address{},
		&model.Article{},
		&model.Banner{},
		// ==main feature==
	)
	if err != nil {
		log.Fatalf("Error performing auto-migration: %v", err)
	}
	log.Println("Database migrated successfully!")
}
