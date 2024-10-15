package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	mu                    sync.Mutex
	current               int
	maxConcurrentRequests = 2 // Set your limit here
)

func LimitConcurrentRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		mu.Lock()
		if current >= maxConcurrentRequests {
			mu.Unlock()
			c.AbortWithStatus(http.StatusTooManyRequests) // Return 429 Too Many Requests
			return
		}
		current++
		mu.Unlock()

		defer func() {
			mu.Lock()
			current--
			mu.Unlock()
		}()

		c.Next() // Process the request
	}
}
