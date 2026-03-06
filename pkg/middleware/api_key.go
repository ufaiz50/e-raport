package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		fmt.Println("Received API Key:", apiKey)                      // Debugging line
		fmt.Println("Expected API Key:", os.Getenv("API_SECRET_KEY")) // Debugging line
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
