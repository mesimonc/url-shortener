package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type URLHandler struct{}

func NewURLHandler() *URLHandler {
    return &URLHandler{}
}

func (h *URLHandler) Shorten(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "shorten endpoint",
    })
}

func (h *URLHandler) Redirect(c *gin.Context) {
    code := c.Param("code")
    c.JSON(http.StatusOK, gin.H{
        "code": code,
    })
}
