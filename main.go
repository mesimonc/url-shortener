package main

import (
	"log"
	"url-shortener/config"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Connect to PostgreSQL via GORM
	db, err := repository.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// Run auto migration
	if err := repository.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Connect to Redis
	cache, err := repository.NewCache(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// Wire up dependencies
	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo, cache)
	urlHandler := handler.NewURLHandler(urlService)

	// Set up structured logger
	logger := handler.NewLogger()
	defer logger.Sync()

	// Use gin.New() to disable default logger, add zap instead
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(handler.LoggerMiddleware(logger))

	// Routes
	r.POST("/shorten", handler.NewRateLimiter("5-M"), urlHandler.Shorten)
	r.GET("/:code", urlHandler.Redirect)
	r.GET("/api/stats/:code", urlHandler.Stats)

	r.Run(cfg.Port)
}