package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirridemirtas/anonsocial/config"
)

type Claims struct {
	UserID       string `json:"userId"`
	Username     string `json:"username"`
	Role         int    `json:"role"`
	UniversityID string `json:"universityId"`
	jwt.StandardClaims
}

func Auth(requiredRole int) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Kimlik doğrulama gerekiyor"}) // Authentication required
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"}) // Invalid token
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || claims.Role < requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Yetkiniz yok"}) // Insufficient permissions
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("userRole", claims.Role)
		c.Set("universityId", claims.UniversityID)
		c.Next()
	}
}

// GetUsernameFromToken extracts username from token if valid
func GetUsernameFromToken(cookie string) (string, bool) {
	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return "", false
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", false
	}

	return claims.Username, true
}

// OptionalAuth middleware for endpoints that work with or without authentication
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			// No token, continue with empty username
			c.Set("username", "")
			c.Next()
			return
		}

		username, ok := GetUsernameFromToken(cookie)
		if !ok {
			// Invalid token, continue with empty username
			c.Set("username", "")
			c.Next()
			return
		}

		// Set username in context
		c.Set("username", username)
		c.Next()
	}
}
