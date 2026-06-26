package handler

import (
	"github.com/gin-gonic/gin"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	limiter "github.com/ulule/limiter/v3"
)

// NewRateLimiter creates a rate limiter middleware.
// Limits requests to the given rate (e.g. "5-M" = 5 requests per minute).
func NewRateLimiter(rate string) gin.HandlerFunc {
	r, err := limiter.NewRateFromFormatted(rate)
	if err != nil {
		panic(err)
	}

	store := memory.NewStore()
	instance := limiter.New(store, r)

	return mgin.NewMiddleware(instance)
}