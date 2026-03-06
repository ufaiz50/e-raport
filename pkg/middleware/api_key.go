package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == os.Getenv("API_SECRET_KEY") {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
		}
	}
}
