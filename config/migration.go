package config

import (
	"log"

	"rijig/model"

	"gorm.io/gorm"
)


func enableUUIDExtension(db *gorm.DB) error {
	// Try to enable uuid-ossp extension
	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Printf("Warning: Could not enable uuid-ossp extension: %v", err)
		// Try alternative: pgcrypto extension (usually available)
		err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error
		if err != nil {
			log.Printf("Warning: Could not enable pgcrypto extension: %v", err)
			return err
		}
	}
	log.Println("UUID extension enabled successfully!")
	return nil
}

func RunMigrations(db *gorm.DB) error {
	log.Println("Starting database migration...")

	if err := enableUUIDExtension(db); err != nil {
		log.Printf("Error enabling UUID extension: %v", err)
		return err
	}

	err := db.AutoMigrate(
		// Location models
		&model.Province{},
		&model.Regency{},
		&model.District{},
		&model.Village{},

		// User related models
		&model.Role{},
		&model.User{},
		&model.Collector{},
		&model.AvaibleTrashByCollector{},
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

	if err := RunSeeders(db); err != nil {
		log.Printf("Error running seeders: %v", err)
		return err
	}

	log.Println("Database migrated successfully!")
	return nil
}
