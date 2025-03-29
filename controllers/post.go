package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirridemirtas/anonsocial/config"
	"github.com/sirridemirtas/anonsocial/data"
	"github.com/sirridemirtas/anonsocial/middleware"
	"github.com/sirridemirtas/anonsocial/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection

func SetPostCollection(client *mongo.Client) {
	postCollection = client.Database(config.AppConfig.MongoDB_DB).Collection("posts")
}

func CreatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input struct {
		Content      string `json:"content" binding:"required,max=500"`
		ReplyTo      string `json:"replyTo,omitempty"`
		UniversityID string `json:"universityId,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")
	userUniversityID := c.GetString("universityId")

	// Determine which universityId to use for post placement
	postUniversityID := userUniversityID
	if input.UniversityID != "" {
		// If a universityId was provided, validate it using IsValidUniversityID
		if data.IsValidUniversityID(input.UniversityID) {
			postUniversityID = input.UniversityID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz üniversite ID'si"})
			return
		}
	}

	post := models.Post{
		Username:         username,
		UniversityID:     postUniversityID, // University where the post appears
		UserUniversityID: userUniversityID, // Always user's own university ID
		Content:          input.Content,
		CreatedAt:        time.Now(),
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
		// We need to also get the privacy status of the post owner
		pipeline := []bson.M{
			{
				"$match": bson.M{"_id": replyToID},
			},
			{
				"$lookup": bson.M{
					"from":         "users",
					"localField":   "username",
					"foreignField": "username",
					"as":           "user",
				},
			},
			{
				"$addFields": bson.M{
					"userIsPrivate": bson.M{
						"$cond": bson.M{
							"if": bson.M{"$gt": []interface{}{bson.M{"$size": "$user"}, 0}},
							"then": bson.M{
								"$toBool": bson.M{
									"$arrayElemAt": []interface{}{"$user.isPrivate", 0},
								},
							},
							"else": false,
						},
					},
				},
			},
			{
				"$project": bson.M{
					"user": 0, // Remove the user array from output
				},
			},
		}

		cursor, err := postCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var parentPosts []models.Post
		if err = cursor.All(ctx, &parentPosts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(parentPosts) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cevaplamak istediğiniz gönderi bulunamadı"}) // Parent post not found
			return
		}

		parentPost := parentPosts[0]

		if parentPost.ReplyTo != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bir cevaba cevap veremezsiniz"}) // Cannot reply to a reply
			return
		}

		// Allow replies to private users' posts (we'll just hide their username in the UI)
		// No need to block replies based on privacy

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
	}

	result, err := postCollection.InsertOne(ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, post)
}

// GetPosts needs to be updated to handle privacy
func GetPosts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username from context or token
	username := getUsernameFromRequest(c)

	// Create pipeline to join with users collection
	pipeline := []bson.M{
		{
			"$match": bson.M{"replyTo": nil}, // Only top-level posts
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "username",
				"foreignField": "username",
				"as":           "user",
			},
		},
		{
			"$addFields": bson.M{
				"userIsPrivate": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$gt": []interface{}{bson.M{"$size": "$user"}, 0}},
						// Use $toBool to convert any type to boolean
						"then": bson.M{
							"$toBool": bson.M{
								"$arrayElemAt": []interface{}{"$user.isPrivate", 0},
							},
						},
						"else": false,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"user": 0, // Remove the user array from output
			},
		},
		{
			"$sort": bson.M{"createdAt": -1}, // Sort by createdAt descending
		},
	}

	cursor, err := postCollection.Aggregate(ctx, pipeline)
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

	// Convert posts to response format with reaction info and privacy handling
	var postResponses []models.PostResponse
	for _, post := range posts {
		// Include all posts, but sanitize usernames for private users
		postResponses = append(postResponses, post.ToResponse(username))
	}

	c.JSON(http.StatusOK, postResponses)
}

func GetPost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"}) // Invalid ID
		return
	}

	// Get username from context or token
	username := getUsernameFromRequest(c)

	// Create a pipeline to join with users collection to get privacy status
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": postID},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "username",
				"foreignField": "username",
				"as":           "user",
			},
		},
		{
			"$addFields": bson.M{
				"userIsPrivate": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$gt": []interface{}{bson.M{"$size": "$user"}, 0}},
						// Use $toBool to convert any type to boolean
						"then": bson.M{
							"$toBool": bson.M{
								"$arrayElemAt": []interface{}{"$user.isPrivate", 0},
							},
						},
						"else": false,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"user": 0, // Remove the user array from output
			},
		},
	}

	cursor, err := postCollection.Aggregate(ctx, pipeline)
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

	if len(posts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gönderi bulunamadı"}) // Post not found
		return
	}

	post := posts[0]

	// Convert to response format with reaction counts
	// This will automatically hide the username if the user is private
	// and the requesting user is not the post owner
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
	if universityId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Üniversite ID parametresi gerekli"})
		return
	}

	// Get all posts from the university and join with users to apply privacy
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"universityId": universityId,
				"replyTo":      nil, // Only include top-level posts, not replies
			},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "username",
				"foreignField": "username",
				"as":           "user",
			},
		},
		{
			"$addFields": bson.M{
				"userIsPrivate": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$gt": []interface{}{bson.M{"$size": "$user"}, 0}},
						// Use $toBool to convert any type to boolean
						"then": bson.M{
							"$toBool": bson.M{
								"$arrayElemAt": []interface{}{"$user.isPrivate", 0},
							},
						},
						"else": false,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"user": 0, // Remove the user array from output
			},
		},
		{
			"$sort": bson.M{"createdAt": -1}, // Sort by createdAt descending
		},
	}

	cursor, err := postCollection.Aggregate(ctx, pipeline)
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

	// Convert posts to response format with reaction info and privacy handling
	var postResponses []models.PostResponse
	for _, post := range posts {
		// Include all posts, but sanitize usernames for private users
		postResponses = append(postResponses, post.ToResponse(username))
	}

	c.JSON(http.StatusOK, postResponses)
}

// GetPostReplies needs to be updated to include reaction info
func GetPostReplies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"}) // Invalid ID
		return
	}

	// Get username from context or token
	username := getUsernameFromRequest(c)

	// First check if parent post exists and if the owner is private
	// Create a pipeline to join with users collection to get privacy status
	pipelineParent := []bson.M{
		{
			"$match": bson.M{"_id": postID},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "username",
				"foreignField": "username",
				"as":           "user",
			},
		},
		{
			"$addFields": bson.M{
				"userIsPrivate": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$gt": []interface{}{bson.M{"$size": "$user"}, 0}},
						// Use $toBool to convert any type to boolean
						"then": bson.M{
							"$toBool": bson.M{
								"$arrayElemAt": []interface{}{"$user.isPrivate", 0},
							},
						},
						"else": false,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"user": 0, // Remove the user array from output
			},
		},
	}

	cursorParent, err := postCollection.Aggregate(ctx, pipelineParent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursorParent.Close(ctx)

	var parentPosts []models.Post
	if err = cursorParent.All(ctx, &parentPosts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(parentPosts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gönderi bulunamadı"}) // Post not found
		return
	}

	// Now get all replies and join with users to get their privacy status
	pipelineReplies := []bson.M{
		{
			"$match": bson.M{"replyTo": postID},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "username",
				"foreignField": "username",
				"as":           "user",
			},
		},
		{
			"$addFields": bson.M{
				"userIsPrivate": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$gt": []interface{}{bson.M{"$size": "$user"}, 0}},
						// Use $toBool to convert any type to boolean
						"then": bson.M{
							"$toBool": bson.M{
								"$arrayElemAt": []interface{}{"$user.isPrivate", 0},
							},
						},
						"else": false,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"user": 0, // Remove the user array from output
			},
		},
		{
			"$sort": bson.M{"createdAt": 1}, // Sort by createdAt ascending
		},
	}

	cursorReplies, err := postCollection.Aggregate(ctx, pipelineReplies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursorReplies.Close(ctx)

	var replies []models.Post
	if err = cursorReplies.All(ctx, &replies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert replies to response format with reaction info and privacy handling
	var replyResponses []models.PostResponse
	for _, reply := range replies {
		// This will handle username privacy - if user is private and requester is not the owner, username will be empty
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
	userRole := c.GetInt("userRole")

	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": postId}).Decode(&post)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Allow users with role 1 or 2 to delete any post
	if post.Username != username && userRole != 1 && userRole != 2 {
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

// GetPostsByUser retrieves all posts by a specific user
func GetPostsByUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username from context or token
	requestingUsername := getUsernameFromRequest(c)

	targetUsername := c.Param("username")
	if targetUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı parametresi gerekli"})
		return
	}

	// First check if user exists and is private
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": targetUsername}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	// If user is private and requester is not the same user, return empty array
	if user.IsPrivate && targetUsername != requestingUsername {
		c.JSON(http.StatusOK, []models.PostResponse{}) // Empty response
		return
	}

	// Now get all posts from the user and apply privacy settings
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"username": targetUsername,
				"replyTo":  nil, // Only include top-level posts, not replies
			},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "username",
				"foreignField": "username",
				"as":           "user",
			},
		},
		{
			"$addFields": bson.M{
				"userIsPrivate": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$gt": []interface{}{bson.M{"$size": "$user"}, 0}},
						// Use $toBool to convert any type to boolean
						"then": bson.M{
							"$toBool": bson.M{
								"$arrayElemAt": []interface{}{"$user.isPrivate", 0},
							},
						},
						"else": false,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"user": 0, // Remove the user array from output
			},
		},
		{
			"$sort": bson.M{"createdAt": -1}, // Sort by createdAt descending
		},
	}

	cursor, err := postCollection.Aggregate(ctx, pipeline)
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

	// Convert posts to response format with reaction info and privacy handling
	var postResponses []models.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, post.ToResponse(requestingUsername))
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
