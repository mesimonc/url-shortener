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

    db, err := repository.NewPostgres(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("failed to connect to postgres: %v", err)
    }
    defer db.Close()

    if err := repository.Migrate(db); err != nil {
        log.Fatalf("failed to migrate: %v", err)
    }

    urlRepo := repository.NewURLRepository(db)
    urlService := service.NewURLService(urlRepo)
    urlHandler := handler.NewURLHandler(urlService)

    r := gin.Default()
    r.POST("/shorten", urlHandler.Shorten)
    r.GET("/:code", urlHandler.Redirect)

    r.Run(cfg.Port)
}