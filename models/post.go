package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	UserID       primitive.ObjectID  `bson:"userId" json:"userId"`
	UniversityID string              `bson:"universityId" json:"universityId" validate:"required,university"`
	Content      string              `bson:"content" json:"content" validate:"required,max=300"`
	ReplyTo      *primitive.ObjectID `bson:"replyTo,omitempty" json:"replyTo,omitempty"`
	CreatedAt    time.Time           `bson:"createdAt" json:"createdAt"`
}
