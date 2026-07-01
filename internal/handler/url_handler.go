package handler

import (
	"url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
    URL           string `json:"url" binding:"required,url"`
    CustomCode    string `json:"custom_code"`
    ExpiresIn int    `json:"expires_in_days"`
}
type URLHandler struct {
	svc *service.URLService
}

// NewURLHandler creates a new URLHandler with the given service.
func NewURLHandler(svc *service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

// @Summary      Shorten a URL
// @Description  Generate a short code for the given URL
// @Accept       json
// @Produce      json
// @Param        request  body      ShortenRequest  true  "Request body"
// @Success      201      {object}  map[string]string
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /shorten [post]
// Shorten handles POST /shorten, creates a short code for the given URL.
func (h *URLHandler) Shorten(c *gin.Context) {
	var req ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	code, err := h.svc.Shorten(req.URL, req.CustomCode, req.ExpiresIn)
    if err != nil {
        if err.Error() == "code already taken" {
            badRequest(c, "custom code already taken")
            return
        }
        internalError(c)
        return
    }

	c.JSON(201, gin.H{"code": code})
}

// @Summary      Redirect to original URL
// @Description  Redirects to the original URL for the given short code
// @Param        code  path  string  true  "Short code"
// @Success      301
// @Failure      404  {object}  errorResponse
// @Router       /{code} [get]
// Redirect handles GET /:code, redirects to the original URL.
func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	originalURL, err := h.svc.Resolve(code)
	if err != nil {
		internalError(c)
		return
	}
	if originalURL == "" {
		notFound(c, "URL not found")
		return
	}

	c.Redirect(301, originalURL)
}

// @Summary      Get URL stats
// @Description  Returns click statistics for the given short code
// @Param        code  path      string  true  "Short code"
// @Success      200   {object}  map[string]interface{}
// @Failure      404   {object}  errorResponse
// @Router       /api/stats/{code} [get]
// Stats handles GET /api/stats/:code, returns click statistics.
func (h *URLHandler) Stats(c *gin.Context) {
	code := c.Param("code")

	url, err := h.svc.GetStats(code)
	if err != nil {
		internalError(c)
		return
	}
	if url == nil {
		notFound(c, "URL not found")
		return
	}

	c.JSON(200, gin.H{
		"code":         url.Code,
		"original_url": url.OriginalURL,
		"clicks":       url.Clicks,
		"created_at":   url.CreatedAt,
	})
}