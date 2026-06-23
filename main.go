package main

import (
    "url-shortener/config"
    "url-shortener/internal/handler"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()

    r := gin.Default()

    urlHandler := handler.NewURLHandler()

    r.POST("/shorten", urlHandler.Shorten)
    r.GET("/:code", urlHandler.Redirect)

    r.Run(cfg.Port)
}