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
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Username     string              `bson:"username" json:"username"`
	UniversityID string              `bson:"universityId" json:"universityId" validate:"required,university"`
	Content      string              `bson:"content" json:"content" validate:"required,max=500"`
	ReplyTo      *primitive.ObjectID `bson:"replyTo,omitempty" json:"replyTo,omitempty"`
	CreatedAt    time.Time           `bson:"createdAt" json:"createdAt"`
	Reactions    Reactions           `bson:"reactions" json:"-"` // Stored but not directly returned
}

// PostResponse is used for API responses, including reaction counts
type PostResponse struct {
	ID           primitive.ObjectID  `json:"id,omitempty"`
	Username     string              `json:"username"`
	UniversityID string              `json:"universityId"`
	Content      string              `json:"content"`
	ReplyTo      *primitive.ObjectID `json:"replyTo,omitempty"`
	CreatedAt    time.Time           `json:"createdAt"`
	Reactions    ReactionCounts      `json:"reactions"`
}

// ToResponse converts a Post to a PostResponse with reaction counts
func (p *Post) ToResponse(username string) PostResponse {
	return PostResponse{
		ID:           p.ID,
		Username:     p.Username,
		UniversityID: p.UniversityID,
		Content:      p.Content,
		ReplyTo:      p.ReplyTo,
		CreatedAt:    p.CreatedAt,
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
