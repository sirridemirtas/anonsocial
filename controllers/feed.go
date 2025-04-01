package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DefaultPageSize = 50
	MaxPageSize     = 50
)

// GetFeedPosts returns posts with reaction counts
func GetFeedPosts(c *gin.Context) {
	pageNum, pageSize, err := getPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz sayfa parametresi"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username for reaction status
	username := getUsernameFromRequest(c)

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64((pageNum - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := postCollection.Find(ctx, bson.M{"replyTo": nil}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	// Initialize an empty array (not null)
	posts := []models.Post{}
	if err = cursor.All(ctx, &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Transform posts to include reaction counts and respect privacy settings
	var postResponses []models.PostResponse
	for _, post := range posts {
		// Check if post owner is private and requester is not the owner
		var isPrivate bool
		if post.Username != username {
			var user models.User
			err := userCollection.FindOne(ctx, bson.M{"username": post.Username}).Decode(&user)
			if err == nil && user.IsPrivate {
				isPrivate = true
			}
		}

		response := post.ToResponse(username)
		if isPrivate {
			response.Username = "" // Clear username for private users
		}
		postResponses = append(postResponses, response)
	}

	c.JSON(http.StatusOK, postResponses)
}

// GetFeedUserPosts returns a user's posts with reaction counts
func GetFeedUserPosts(c *gin.Context) {
	pageNum, pageSize, err := getPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz sayfa parametresi"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username for reaction status
	username := getUsernameFromRequest(c)

	targetUsername := c.Param("username")

	// First check if user exists and is private
	var targetUser models.User
	err = userCollection.FindOne(ctx, bson.M{"username": targetUsername}).Decode(&targetUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	// If user is private and requester is not the same user, return empty array
	if targetUser.IsPrivate && targetUsername != username {
		c.JSON(http.StatusOK, []models.PostResponse{}) // Empty response
		return
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64((pageNum - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := postCollection.Find(ctx, bson.M{
		"username": targetUsername,
		"replyTo":  nil,
	}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	// Initialize an empty array (not null)
	posts := []models.Post{}
	if err = cursor.All(ctx, &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Transform posts to include reaction counts
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse(username))
	}

	c.JSON(http.StatusOK, postResponses)
}

// GetFeedUniversityPosts returns university posts with reaction counts
func GetFeedUniversityPosts(c *gin.Context) {
	pageNum, pageSize, err := getPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz sayfa parametresi"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username for reaction status
	username := getUsernameFromRequest(c)

	universityId := c.Param("universityId")

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64((pageNum - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := postCollection.Find(ctx, bson.M{
		"universityId": universityId,
		"replyTo":      nil,
	}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	// Initialize an empty array (not null)
	posts := []models.Post{}
	if err = cursor.All(ctx, &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Transform posts to include reaction counts and respect privacy settings
	var postResponses []models.PostResponse
	for _, post := range posts {
		// Check if post owner is private and requester is not the owner
		var isPrivate bool
		if post.Username != username {
			var user models.User
			err := userCollection.FindOne(ctx, bson.M{"username": post.Username}).Decode(&user)
			if err == nil && user.IsPrivate {
				isPrivate = true
			}
		}

		response := post.ToResponse(username)
		if isPrivate {
			response.Username = "" // Clear username for private users
		}
		postResponses = append(postResponses, response)
	}

	c.JSON(http.StatusOK, postResponses)
}

// Helper function to get pagination parameters from the request
// Returns pageNum, pageSize, and error
func getPaginationParams(c *gin.Context) (int, int, error) {
	// Default page is 1
	pageNum := 1
	if page := c.Query("page"); page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt <= 0 {
			// Return error if page is not a valid positive number
			return 0, 0, err
		}
		pageNum = pageInt
	}

	// Default and maximum page size
	pageSize := DefaultPageSize
	if size := c.Query("size"); size != "" {
		sizeInt, err := strconv.Atoi(size)
		if err == nil && sizeInt > 0 {
			pageSize = sizeInt
			if pageSize > MaxPageSize {
				pageSize = MaxPageSize
			}
		}
	}

	return pageNum, pageSize, nil
}
