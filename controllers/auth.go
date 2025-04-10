package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/sirridemirtas/anonsocial/config"
	"github.com/sirridemirtas/anonsocial/middleware"
	"github.com/sirridemirtas/anonsocial/models"
	"github.com/sirridemirtas/anonsocial/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input struct {
		Username     string `json:"username" binding:"required"`
		Password     string `json:"password" binding:"required"`
		UniversityID string `json:"universityId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Create user object with only allowed fields
	user := models.User{
		Username:     input.Username,
		UniversityID: input.UniversityID,
		CreatedAt:    time.Now(),
		Role:         0, // Default role
		IsPrivate:    false,
	}

	// Generate salt and hash password
	user.Salt = models.GenerateSalt()
	user.Password = user.HashPassword(input.Password)

	// Validate user data
	if errors := utils.ValidateUser(&user); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) || strings.Contains(err.Error(), "duplicate key error") {
			c.JSON(http.StatusConflict, gin.H{"message": "Kullanıcı adı zaten alınmış"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	/* // Generate JWT token after successful registration
	claims := &middleware.Claims{
		UserID:       user.ID.Hex(),
		Username:     user.Username,
		Role:         user.Role,
		UniversityID: user.UniversityID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: getTokenExpiration().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token oluşturulamadı"})
		return
	}

	cookieDomain := config.AppConfig.CookieDomain
	if cookieDomain == "" {
		cookieDomain = ""
	}

	c.SetCookie("token", tokenString, 3600*24, "/", cookieDomain, false, true) */
	c.JSON(http.StatusCreated, gin.H{
		"message": "Kayıt başarılı",
		"user": gin.H{
			"id":           user.ID,
			"username":     user.Username,
			"role":         user.Role,
			"universityId": user.UniversityID,
		},
	})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Kullanıcı adı veya şifre hatalı"})
		return
	}

	if !user.ValidatePassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Kullanıcı adı veya şifre hatalı"})
		return
	}

	claims := &middleware.Claims{
		UserID:       user.ID.Hex(),
		Username:     user.Username,
		Role:         user.Role,
		UniversityID: user.UniversityID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: getTokenExpiration().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token oluşturulamadı"})
		return
	}

	// Cookie domain'i config'den al
	cookieDomain := config.AppConfig.CookieDomain
	if cookieDomain == "" {
		cookieDomain = ""
	}

	c.SetCookie("token", tokenString, 3600*24, "/", cookieDomain, false, true)
	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID.Hex(),
		"username":     user.Username,
		"role":         user.Role,
		"universityId": user.UniversityID,
	},
	)
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Çıkış başarılı"})
}

func TokenInfo(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token bilgileri alınamadı"})
		return
	}

	tokenClaims := claims.(*middleware.Claims)

	// Check if refresh parameter is set to true
	refresh := c.Query("refresh")
	if refresh == "true" {
		// Get up-to-date user information from database
		var user models.User
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := userCollection.FindOne(ctx, bson.M{"username": tokenClaims.Username}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Kullanıcı bilgileri alınamadı"})
			return
		}

		// Create new claims with latest user data
		newClaims := &middleware.Claims{
			UserID:       user.ID.Hex(),
			Username:     user.Username,
			Role:         user.Role,
			UniversityID: user.UniversityID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: getTokenExpiration().Unix(),
			},
		}

		// Generate new token string
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
		tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Yeni token oluşturulamadı"})
			return
		}

		// Get cookie domain from config
		cookieDomain := config.AppConfig.CookieDomain
		if cookieDomain == "" {
			cookieDomain = ""
		}

		// Set the new cookie with the refreshed token
		c.SetCookie("token", tokenString, 3600*24, "/", cookieDomain, false, true)

		// Update token claims for the response
		tokenClaims = newClaims
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":       tokenClaims.UserID,
		"username":     tokenClaims.Username,
		"role":         tokenClaims.Role,
		"universityId": tokenClaims.UniversityID,
		"expiresAt":    time.Unix(tokenClaims.ExpiresAt, 0),
		"refreshed":    refresh == "true",
	})
}

func getTokenExpiration() time.Time {
	// Get the expiration time from config
	expiresInStr := config.AppConfig.JWTExpiresIn

	// Parse the string to an integer (assuming it's stored in hours)
	expiresInHours, err := strconv.Atoi(expiresInStr)
	if err != nil {
		// Default to 24 hours if there's an error parsing
		expiresInHours = 24
	}

	// Calculate expiration time
	expirationTime := time.Now().Add(time.Duration(expiresInHours) * time.Hour)

	return expirationTime
}

func RefreshToken(c *gin.Context) {
	// Get claims from context that were set by Auth middleware
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token bilgileri alınamadı"})
		return
	}

	tokenClaims := claims.(*middleware.Claims)

	// Create a new token with the same user information but new expiration time
	newClaims := &middleware.Claims{
		UserID:       tokenClaims.UserID,
		Username:     tokenClaims.Username,
		Role:         tokenClaims.Role,
		UniversityID: tokenClaims.UniversityID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: getTokenExpiration().Unix(),
		},
	}

	// Generate new token string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Yeni token oluşturulamadı"})
		return
	}

	// Get cookie domain from config
	cookieDomain := config.AppConfig.CookieDomain
	if cookieDomain == "" {
		cookieDomain = ""
	}

	// Set the new cookie with the refreshed token (replacing the old one)
	// Using the same cookie settings as in the Login function
	c.SetCookie("token", tokenString, 3600*24, "/", cookieDomain, false, true)

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Token yenilendi"})
}
