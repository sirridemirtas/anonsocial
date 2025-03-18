package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Username     string              `bson:"username" json:"username"`
	UniversityID string              `bson:"universityId" json:"universityId" validate:"required,university"`
	Content      string              `bson:"content" json:"content" validate:"required,max=500"`
	ReplyTo      *primitive.ObjectID `bson:"replyTo,omitempty" json:"replyTo,omitempty"`
	CreatedAt    time.Time           `bson:"createdAt" json:"createdAt"`
}
