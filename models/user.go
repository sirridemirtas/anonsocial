package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username     string             `bson:"username" json:"username" validate:"required,alphanum,min=3,max=16"`
	Password     string             `bson:"password" json:"-"` // Never send in JSON
	IsPrivate    bool               `bson:"isPrivate" json:"isPrivate"`
	Role         int                `bson:"role" json:"role"`
	UniversityID string             `bson:"universityId" json:"universityId" validate:"required,university"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	Salt         string             `bson:"salt" json:"-"` // Never send in JSON
}

func GenerateSalt() string {
	salt := make([]byte, 16)
	rand.Read(salt)
	return hex.EncodeToString(salt)
}

func (u *User) HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password + u.Salt))
	return hex.EncodeToString(hash.Sum(nil))
}

func (u *User) ValidatePassword(password string) bool {
	return u.Password == u.HashPassword(password)
}
