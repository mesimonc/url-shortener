package handler

import (
    "net/http"
    "url-shortener/internal/service"

    "github.com/gin-gonic/gin"
)

type URLHandler struct {
    svc *service.URLService
}

func NewURLHandler(svc *service.URLService) *URLHandler {
    return &URLHandler{svc: svc}
}

func (h *URLHandler) Shorten(c *gin.Context) {
    var req struct {
        URL string `json:"url"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    code, err := h.svc.Shorten(req.URL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to shorten"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"code": code})
}

func (h *URLHandler) Redirect(c *gin.Context) {
    code := c.Param("code")
    c.JSON(http.StatusOK, gin.H{"code": code})
}