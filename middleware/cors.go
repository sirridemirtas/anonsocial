package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func Cors() gin.HandlerFunc {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://agora-313.pages.dev"},
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
		// Handle preflight OPTIONS requests specially
		if c.Request.Method == http.MethodOptions {
			corsMiddleware.HandlerFunc(c.Writer, c.Request)
			c.AbortWithStatus(http.StatusOK)
			return
		}

		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	}
}
