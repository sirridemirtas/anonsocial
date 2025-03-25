package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Reactions stores user reactions to a post
type Reactions struct {
	Likes    []string `bson:"likes" json:"-"`    // List of usernames who liked
	Dislikes []string `bson:"dislikes" json:"-"` // List of usernames who disliked
}

// ReactionCounts represents reaction counts for API responses
type ReactionCounts struct {
	LikeCount    int  `json:"likeCount"`
	DislikeCount int  `json:"dislikeCount"`
	Liked        bool `json:"liked,omitempty"`
	Disliked     bool `json:"disliked,omitempty"`
}

type Post struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Username         string              `bson:"username" json:"username"`
	UniversityID     string              `bson:"universityId" json:"universityId" validate:"required,university"`
	UserUniversityID string              `bson:"userUniversityId" json:"userUniversityId"` // User's own university ID
	Content          string              `bson:"content" json:"content" validate:"required,max=500"`
	ReplyTo          *primitive.ObjectID `bson:"replyTo,omitempty" json:"replyTo,omitempty"`
	CreatedAt        time.Time           `bson:"createdAt" json:"createdAt"`
	Reactions        Reactions           `bson:"reactions" json:"-"`     // Stored but not directly returned
	UserIsPrivate    bool                `bson:"userIsPrivate" json:"-"` // Internal field not to be exposed in JSON
}

// PostResponse is used for API responses, including reaction counts
type PostResponse struct {
	ID               primitive.ObjectID  `json:"id,omitempty"`
	Username         string              `json:"username"`
	UniversityID     string              `json:"universityId"`
	UserUniversityID string              `json:"userUniversityId"` // User's own university ID
	Content          string              `json:"content"`
	ReplyTo          *primitive.ObjectID `json:"replyTo,omitempty"`
	CreatedAt        time.Time           `json:"createdAt"`
	Reactions        ReactionCounts      `json:"reactions"`
}

// ToResponse converts a Post to a PostResponse with reaction counts
func (p *Post) ToResponse(username string) PostResponse {
	// Create a copy of the post to sanitize
	postCopy := *p

	// Apply privacy settings to the username
	if postCopy.UserIsPrivate && postCopy.Username != username {
		postCopy.Username = "" // Hide username if user is private and requester is not the owner
	}

	return PostResponse{
		ID:               postCopy.ID,
		Username:         postCopy.Username, // This will be empty if user is private and requester is not the owner
		UniversityID:     postCopy.UniversityID,
		UserUniversityID: postCopy.UserUniversityID, // Include user's university ID in all responses
		Content:          postCopy.Content,
		ReplyTo:          postCopy.ReplyTo,
		CreatedAt:        postCopy.CreatedAt,
		Reactions: ReactionCounts{
			LikeCount:    len(p.Reactions.Likes),
			DislikeCount: len(p.Reactions.Dislikes),
			Liked:        contains(p.Reactions.Likes, username),
			Disliked:     contains(p.Reactions.Dislikes, username),
		},
	}
}

// contains checks if a slice contains a string
func contains(slice []string, str string) bool {
	if str == "" {
		return false
	}

	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
