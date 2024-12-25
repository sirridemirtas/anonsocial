package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func FormData() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "multipart/form-data") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be multipart/form-data"})
			c.Abort()
			return
		}

		err := c.Request.ParseMultipartForm(32 << 20) // 32MB max memory
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
			c.Abort()
			return
		}

		c.Next()
	}
}
