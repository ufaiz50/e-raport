package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://127.0.0.1",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8001",
			"http://localhost",
			"http://localhost:3000",
			"http://localhost:8001"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key", "X-Requested-With"},
		ExposeHeaders: []string{"Set-Cookie"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	})
}
