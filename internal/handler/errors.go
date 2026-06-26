package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// errorResponse represents a standard error response body.
type errorResponse struct {
	Error string `json:"error"`
}

// badRequest returns a 400 response with the given message.
func badRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, errorResponse{Error: msg})
}

// notFound returns a 404 response with the given message.
func notFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, errorResponse{Error: msg})
}

// internalError returns a 500 response with a generic message.
func internalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, errorResponse{Error: "internal server error"})
}
