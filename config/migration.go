package config

import (
	"log"

	"rijig/model"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("Starting database migration...")

	err := db.AutoMigrate(
		// Location models
		&model.Province{},
		&model.Regency{},
		&model.District{},
		&model.Village{},

		// User related models
		&model.User{},
		&model.Collector{},
		&model.AvaibleTrashByCollector{},
		&model.Role{},
		&model.UserPin{},
		&model.Address{},
		&model.IdentityCard{},
		&model.CompanyProfile{},

		// Pickup related models
		&model.RequestPickup{},
		&model.RequestPickupItem{},
		&model.PickupStatusHistory{},
		&model.PickupRating{},

		// Cart related models
		&model.Cart{},
		&model.CartItem{},

		// Store related models
		&model.Store{},
		&model.Product{},
		&model.ProductImage{},

		// Content models
		&model.Article{},
		&model.Banner{},
		&model.InitialCoint{},
		&model.About{},
		&model.AboutDetail{},
		&model.CoverageArea{},

		// Trash related models
		&model.TrashCategory{},
		&model.TrashDetail{},
	)

	if err != nil {
		log.Printf("Error performing auto-migration: %v", err)
		return err
	}

	log.Println("Database migrated successfully!")
	return nil
}