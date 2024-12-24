package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username" validate:"required,alphanum,min=3,max=16"`
	Password     string             `bson:"password" json:"password" validate:"required"`
	IsPrivate    bool               `bson:"isPrivate" json:"isPrivate"`
	Role         int                `bson:"role" json:"role"`
	UniversityID string             `bson:"universityId" json:"universityId" validate:"required,university"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}
