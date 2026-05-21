package main

import (
	"log"

	"github.com/gabrzb/auth-go-gin/internal/config"
	"github.com/gabrzb/auth-go-gin/internal/database"
	"github.com/gabrzb/auth-go-gin/internal/handlers"
	"github.com/gabrzb/auth-go-gin/internal/routes"
	"github.com/gabrzb/auth-go-gin/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)
	database.Migrate(db)

	jwtService, err := services.NewJWTService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize JWT service: %v", err)
	}

	authService := services.NewAuthService(db, jwtService)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()
	routes.Setup(r, authHandler)

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
