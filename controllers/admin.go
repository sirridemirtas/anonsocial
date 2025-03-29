package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateUserRole allows administrators (role=2) to update other users' roles
func UpdateUserRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the admin's role from the context (set by Auth middleware)
	adminRole := c.GetInt("userRole")
	if adminRole != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu işlem için yönetici yetkileri gerekiyor"}) // Admin privileges required
		return
	}

	// Get username of the user to update
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı parametresi zorunludur"}) // Username parameter is required
		return
	}

	// Parse the requested role from request body
	// Using pointer to int (*int) to distinguish between "not provided" and "value is 0"
	var input struct {
		Role *int `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if Role was provided
	if input.Role == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Yetki seviyesi zorunludur"}) // Role level is required
		return
	}

	// Validate the role value (must be 0 or 1)
	if *input.Role != 0 && *input.Role != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Yetki seviyesi 0 veya 1 olmalıdır"}) // Role must be 0 or 1
		return
	}

	// Use collation option for case-insensitive search
	opts := options.Update().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2, // Case-insensitive comparison
	})

	// Update the user's role
	update := bson.M{
		"$set": bson.M{
			"role": *input.Role,
		},
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"username": username}, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"}) // User not found
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı yetkisi güncellendi"}) // User role updated
}
