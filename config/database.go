package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var dsn string

	// Check if running on Railway (DATABASE_URL provided)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Println("Using Railway DATABASE_URL")
		dsn = processDatabaseURL(databaseURL)
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

	log.Printf("Connecting to database with DSN: %s", sanitizeDSN(dsn))

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

// processDatabaseURL handles Railway's DATABASE_URL format
func processDatabaseURL(databaseURL string) string {
	// Railway DATABASE_URL format: 
	// postgresql://username:password@hostname:port/database?sslmode=require
	
	// Jika sudah format yang benar, langsung gunakan
	if strings.HasPrefix(databaseURL, "postgresql://") {
		// Convert postgresql:// to postgres:// if needed
		dsn := strings.Replace(databaseURL, "postgresql://", "postgres://", 1)
		
		// Ensure sslmode is set correctly for Railway
		if !strings.Contains(dsn, "sslmode=") {
			separator := "?"
			if strings.Contains(dsn, "?") {
				separator = "&"
			}
			dsn += separator + "sslmode=require"
		}
		
		return dsn
	}
	
	// Jika format postgres://, langsung gunakan
	if strings.HasPrefix(databaseURL, "postgres://") {
		// Ensure sslmode is set correctly for Railway
		if !strings.Contains(databaseURL, "sslmode=") {
			separator := "?"
			if strings.Contains(databaseURL, "?") {
				separator = "&"
			}
			databaseURL += separator + "sslmode=require"
		}
		
		return databaseURL
	}
	
	// Fallback: if it's not URL format, assume it's already a DSN
	return databaseURL
}

// sanitizeDSN removes password from DSN for logging
func sanitizeDSN(dsn string) string {
	// Hide password in logs for security
	if strings.Contains(dsn, "password=") {
		parts := strings.Split(dsn, " ")
		for i, part := range parts {
			if strings.HasPrefix(part, "password=") {
				parts[i] = "password=****"
			}
		}
		return strings.Join(parts, " ")
	}
	
	// For URL format, hide password
	if strings.Contains(dsn, "://") && strings.Contains(dsn, "@") {
		// Find the password part in URL
		beforeAt := strings.Split(dsn, "@")[0]
		afterAt := strings.Split(dsn, "@")[1]
		
		if strings.Contains(beforeAt, ":") {
			protocol := strings.Split(beforeAt, "://")[0]
			userpass := strings.Split(beforeAt, "://")[1]
			
			if strings.Contains(userpass, ":") {
				user := strings.Split(userpass, ":")[0]
				return protocol + "://" + user + ":****@" + afterAt
			}
		}
	}
	
	return dsn
}