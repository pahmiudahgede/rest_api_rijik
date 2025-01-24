package config

import (
	"fmt"
	"log"

	"github.com/pahmiudahgede/senggoldong/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {

	InitConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBName, DBPassword,
	)

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error: Gagal terhubung ke database: %v", err)
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
		log.Fatalf("Error: Gagal melakukan migrasi schema: %v", err)
	}

	log.Println("Koneksi ke database berhasil dan migrasi schema juga berhasil")
}
