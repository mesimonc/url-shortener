package main

import (
    "url-shortener/config"
    "url-shortener/internal/handler"
    "url-shortener/internal/service"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()

    r := gin.Default()

    urlService := service.NewURLService()
    urlHandler := handler.NewURLHandler(urlService)

    r.POST("/shorten", urlHandler.Shorten)
    r.GET("/:code", urlHandler.Redirect)

    r.Run(cfg.Port)
}