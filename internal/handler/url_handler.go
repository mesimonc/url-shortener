package handler

import (
	"net/http"
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	svc *service.URLService
}

// NewURLHandler creates a new URLHandler with the given service.
func NewURLHandler(svc *service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

// Shorten handles POST /shorten, creates a short code for the given URL.
func (h *URLHandler) Shorten(c *gin.Context) {
	var req struct {
		URL        string `json:"url"`
		CustomCode string `json:"custom_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.svc.Shorten(req.URL, req.CustomCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to shorten"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": code})
}

// Redirect handles GET /:code, redirects to the original URL.
func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	originalURL, err := h.svc.Resolve(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if originalURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// Stats handles GET /api/stats/:code, returns click statistics.
func (h *URLHandler) Stats(c *gin.Context) {
	code := c.Param("code")

	url, err := h.svc.GetStats(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if url == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":         url.Code,
		"original_url": url.OriginalURL,
		"clicks":       url.Clicks,
		"created_at":   url.CreatedAt,
	})
}