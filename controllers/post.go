package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirridemirtas/anonsocial/config"
	"github.com/sirridemirtas/anonsocial/middleware"
	"github.com/sirridemirtas/anonsocial/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var postCollection *mongo.Collection

func SetPostCollection(client *mongo.Client) {
	postCollection = client.Database("anonsocial").Collection("posts")
}

func CreatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input struct {
		Content string `json:"content" binding:"required,max=500"`
		ReplyTo string `json:"replyTo,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")
	universityID := c.GetString("universityId")

	post := models.Post{
		Username:  username,
		Content:   input.Content,
		CreatedAt: time.Now(),
		Reactions: models.Reactions{
			Likes:    []string{},
			Dislikes: []string{},
		},
	}

	if input.ReplyTo != "" {
		// Check if this is a reply to a post
		replyToID, err := primitive.ObjectIDFromHex(input.ReplyTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz cevap gönderi kimliği"}) // Invalid reply post ID
			return
		}

		// Check if parent post exists and is not a reply itself
		var parentPost models.Post
		err = postCollection.FindOne(ctx, bson.M{"_id": replyToID}).Decode(&parentPost)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cevaplamak istediğiniz gönderi bulunamadı"}) // Parent post not found
			return
		}

		if parentPost.ReplyTo != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bir cevaba cevap veremezsiniz"}) // Cannot reply to a reply
			return
		}

		post.ReplyTo = &replyToID

		// After saving the reply, create notifications
		// 1. Notify the parent post owner about the reply
		CreateOrUpdateReplyNotification(replyToID, parentPost.Username, parentPost.Content, username, false)

		// 2. If this is not a direct reply to the parent post owner's post,
		// also find other people who replied to notify them
		if parentPost.Username != username {
			// Find all unique users who replied to this post (excluding current user and post owner)
			pipeline := []bson.M{
				{"$match": bson.M{"replyTo": replyToID}},
				{"$group": bson.M{"_id": "$username"}},
			}

			cursor, err := postCollection.Aggregate(ctx, pipeline)
			if err == nil {
				var results []struct {
					Username string `bson:"_id"`
				}
				if err := cursor.All(ctx, &results); err == nil {
					for _, result := range results {
						if result.Username != username && result.Username != parentPost.Username {
							// Notify other users who replied to this post
							CreateOrUpdateReplyNotification(replyToID, result.Username, parentPost.Content, username, true)
						}
					}
				}
				cursor.Close(ctx)
			}
		}
	} else {
		post.UniversityID = universityID
	}

	result, err := postCollection.InsertOne(ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, post)
}

// GetPosts needs to be updated to include reaction info
func GetPosts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username from context or token
	username := getUsernameFromRequest(c)

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := postCollection.Find(ctx, bson.M{"replyTo": nil}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert posts to response format with reaction info
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse(username))
	}

	c.JSON(http.StatusOK, postResponses)
}

func GetPost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"}) // Invalid ID
		return
	}

	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gönderi bulunamadı"}) // Post not found
		return
	}

	// Get username from context (if auth middleware has been applied)
	username := c.GetString("username")

	// If no username in context, try to extract it from the token if available
	if username == "" {
		cookie, err := c.Cookie("token")
		if err == nil {
			// Token exists, try to parse it
			token, err := jwt.ParseWithClaims(cookie, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.AppConfig.JWTSecret), nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(*middleware.Claims); ok {
					username = claims.Username
				}
			}
		}
	}

	// Convert to response format with reaction counts
	response := post.ToResponse(username)

	c.JSON(http.StatusOK, response)
}

// GetPostsByUniversity needs to be updated to include reaction info
func GetPostsByUniversity(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username from context or token
	username := getUsernameFromRequest(c)

	universityId := c.Param("universityId")
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := postCollection.Find(ctx, bson.M{
		"universityId": universityId,
		"replyTo":      nil,
	}, opts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert posts to response format with reaction info
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse(username))
	}

	c.JSON(http.StatusOK, postResponses)
}

// GetPostReplies needs to be updated to include reaction info
func GetPostReplies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username from context or token
	username := getUsernameFromRequest(c)

	postId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"}) // Invalid ID
		return
	}

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}})
	cursor, err := postCollection.Find(ctx, bson.M{"replyTo": postId}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var replies []models.Post
	if err = cursor.All(ctx, &replies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert replies to response format with reaction info
	var replyResponses []models.PostResponse
	for _, reply := range replies {
		replyResponses = append(replyResponses, reply.ToResponse(username))
	}

	c.JSON(http.StatusOK, replyResponses)
}

func DeletePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"}) // Invalid ID
		return
	}

	username := c.GetString("username")

	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": postId}).Decode(&post)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.Username != username && c.GetInt("userRole") != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu gönderiyi silme izniniz yok"}) // Not authorized to delete this post
		return
	}

	_, err = postCollection.DeleteOne(ctx, bson.M{"_id": postId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete all replies to this post
	_, err = postCollection.DeleteMany(ctx, bson.M{"replyTo": postId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gönderi ve cevapları silindi"}) // Post and its replies deleted successfully
}

// GetPostsByUser needs to be updated to include reaction info
func GetPostsByUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username from context or token
	username := getUsernameFromRequest(c)

	targetUsername := c.Param("username")
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := postCollection.Find(ctx, bson.M{
		"username": targetUsername,
		"replyTo":  nil,
	}, opts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert posts to response format with reaction info
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse(username))
	}

	c.JSON(http.StatusOK, postResponses)
}

// LikePost handles adding a like to a post
func LikePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz gönderi kimliği"}) // Invalid post ID
		return
	}

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"}) // Username not found
		return
	}

	// Find the post first to check if it exists
	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gönderi bulunamadı"}) // Post not found
		return
	}

	// Add username to likes and remove from dislikes if present
	update := bson.M{
		"$addToSet": bson.M{"reactions.likes": username},
		"$pull":     bson.M{"reactions.dislikes": username},
	}

	_, err = postCollection.UpdateOne(ctx, bson.M{"_id": postID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a notification for the post owner if it's not the same user
	if post.Username != username {
		CreateOrUpdateReactionNotification(postID, post.Username, post.Content, true)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gönderi beğenildi"}) // Post liked
}

// DislikePost handles adding a dislike to a post
func DislikePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz gönderi kimliği"}) // Invalid post ID
		return
	}

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"}) // Username not found
		return
	}

	// Find the post first to check if it exists
	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gönderi bulunamadı"}) // Post not found
		return
	}

	// Add username to dislikes and remove from likes if present
	update := bson.M{
		"$addToSet": bson.M{"reactions.dislikes": username},
		"$pull":     bson.M{"reactions.likes": username},
	}

	_, err = postCollection.UpdateOne(ctx, bson.M{"_id": postID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a notification for the post owner if it's not the same user
	if post.Username != username {
		CreateOrUpdateReactionNotification(postID, post.Username, post.Content, false)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gönderi beğenilmedi"}) // Post disliked
}

// RemoveLikePost handles removing a like from a post
func RemoveLikePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz gönderi kimliği"}) // Invalid post ID
		return
	}

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"}) // Username not found
		return
	}

	// Find the post first to check if it exists
	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gönderi bulunamadı"}) // Post not found
		return
	}

	// Remove username from likes if present
	update := bson.M{
		"$pull": bson.M{"reactions.likes": username},
	}

	_, err = postCollection.UpdateOne(ctx, bson.M{"_id": postID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beğeni kaldırıldı"}) // Like removed
}

// RemoveDislikePost handles removing a dislike from a post
func RemoveDislikePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz gönderi kimliği"}) // Invalid post ID
		return
	}

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"}) // Username not found
		return
	}

	// Find the post first to check if it exists
	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gönderi bulunamadı"}) // Post not found
		return
	}

	// Remove username from dislikes if present
	update := bson.M{
		"$pull": bson.M{"reactions.dislikes": username},
	}

	_, err = postCollection.UpdateOne(ctx, bson.M{"_id": postID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Beğenmeme kaldırıldı"}) // Dislike removed
}

// Add a helper function to extract username from request context or token
func getUsernameFromRequest(c *gin.Context) string {
	// Get username from context (if auth middleware has been applied)
	username := c.GetString("username")

	// If no username in context, try to extract it from the token if available
	if username == "" {
		cookie, err := c.Cookie("token")
		if err == nil {
			// Token exists, try to parse it
			token, err := jwt.ParseWithClaims(cookie, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.AppConfig.JWTSecret), nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(*middleware.Claims); ok {
					username = claims.Username
				}
			}
		}
	}

	return username
}
