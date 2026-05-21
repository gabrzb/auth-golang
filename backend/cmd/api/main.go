package main

import (
	"log"
	"net/http"

	"github.com/gabrzb/auth-go-gin/internal/config"
	"github.com/gabrzb/auth-go-gin/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
