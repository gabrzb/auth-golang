package database

import (
	"fmt"
	"log"

	"github.com/gabrzb/auth-go-gin/internal/config"
	"github.com/gabrzb/auth-go-gin/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected")
	return db
}

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrated")
}
