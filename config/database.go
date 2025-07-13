package config

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var dsn string

	// Check if running on Railway (DATABASE_URL provided)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Println("Using Railway DATABASE_URL")
		dsn = databaseURL
	} else {
		// Fallback to individual environment variables (for local development)
		log.Println("Using individual database environment variables")
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	}

	// Parse and modify DSN if needed for Railway
	if parsedURL, err := url.Parse(dsn); err == nil && parsedURL.Scheme == "postgresql" {
		// Railway sometimes uses postgresql:// scheme, convert to postgres://
		dsn = "postgres://" + dsn[13:]
		
		// Ensure SSL mode is set correctly for Railway
		if !contains(dsn, "sslmode=") {
			if contains(dsn, "?") {
				dsn += "&sslmode=require"
			} else {
				dsn += "?sslmode=require"
			}
		}
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Database connected successfully!")

	if err := RunMigrations(DB); err != nil {
		log.Fatalf("Error performing auto-migration: %v", err)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
			 (s[:len(substr)] == substr || 
			  s[len(s)-len(substr):] == substr || 
			  containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}