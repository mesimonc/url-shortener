package main

// @title           URL Shortener API
// @version         1.0
// @description     A URL shortening service built with Go, Gin, PostgreSQL and Redis.
// @host            localhost:8080
// @BasePath        /

import (
	"log"
	"url-shortener/config"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
	_ "url-shortener/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(cfg.Port)
}