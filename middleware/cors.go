package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func Cors() gin.HandlerFunc {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"accept",
			"origin",
			"Cache-Control",
			"X-Requested-With",
		},
		AllowedMethods: []string{"POST", "OPTIONS", "GET", "PUT", "DELETE", "PATCH"},
		ExposedHeaders: []string{"Content-Length"},
	})

	return func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	}
}
