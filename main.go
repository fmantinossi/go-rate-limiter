package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-rate-limiter/config"
	"github.com/go-rate-limiter/limiter"
	"github.com/go-rate-limiter/storage"
)

func main() {
	cfg := config.LoadConfig()

	redisStorage, err := storage.NewRedisStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis storage: %v", err)
	}

	rateLimiter := limiter.NewRateLimiter(redisStorage, cfg)

	router := gin.Default()

	router.Use(rateLimiter.Middleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
