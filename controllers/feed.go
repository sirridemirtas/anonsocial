package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Get current user's username
	username := c.GetString("username")

	// Transform posts to include reaction counts
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse(username))
	}

	c.JSON(http.StatusOK, postResponses)
}

// GetFeedPostReplies returns replies with reaction counts
func GetFeedPostReplies(c *gin.Context) {
	pageNum, pageSize, err := getPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz sayfa parametresi"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"}) // Invalid ID
		return
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: 1}}).
		SetSkip(int64((pageNum - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := postCollection.Find(ctx, bson.M{"replyTo": postId}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	// Initialize an empty array (not null)
	replies := []models.Post{}
	if err = cursor.All(ctx, &replies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get current user's username
	username := c.GetString("username")

	// Transform replies to include reaction counts
	var replyResponses []models.PostResponse
	for _, reply := range replies {
		replyResponses = append(replyResponses, reply.ToResponse(username))
	}

	c.JSON(http.StatusOK, replyResponses)
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

	username := c.Param("username")

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64((pageNum - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := postCollection.Find(ctx, bson.M{
		"username": username,
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

	// Get current user's username
	currentUsername := c.GetString("username")

	// Transform posts to include reaction counts
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse(currentUsername))
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

	// Get current user's username
	username := c.GetString("username")

	// Transform posts to include reaction counts
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse(username))
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
