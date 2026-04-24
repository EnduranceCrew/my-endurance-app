package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout aborts requests that exceed the given duration.
func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		done := make(chan struct{}, 1)
		go func() {
			c.Next()
			done <- struct{}{}
		}()
		select {
		case <-done:
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"success": false,
				"error":   "tempo limite da requisição excedido",
			})
		}
	}
}
