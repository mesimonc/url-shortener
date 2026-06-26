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
    //defer db.Close()

    // Run database migrations
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

    r := gin.Default()
	// Apply rate limiter to /shorten only
    r.POST("/shorten", handler.NewRateLimiter("5-M"), urlHandler.Shorten)
    r.GET("/:code", urlHandler.Redirect)
	r.GET("/api/stats/:code", urlHandler.Stats)

    r.Run(cfg.Port)
}