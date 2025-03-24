package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Avatar represents a user's avatar configuration
type Avatar struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Username     string             `bson:"username" json:"username"`
	FaceColor    string             `bson:"faceColor" json:"faceColor" validate:"required,hexcolor"`
	EarSize      string             `bson:"earSize" json:"earSize" validate:"required,oneof=small big"`
	HairStyle    string             `bson:"hairStyle" json:"hairStyle" validate:"required,oneof=normal thick mohawk womanLong womanShort"`
	HairColor    string             `bson:"hairColor" json:"hairColor" validate:"required,hexcolor"`
	HatStyle     string             `bson:"hatStyle" json:"hatStyle" validate:"required,oneof=none beanie turban"`
	HatColor     string             `bson:"hatColor" json:"hatColor" validate:"required,hexcolor"`
	EyeStyle     string             `bson:"eyeStyle" json:"eyeStyle" validate:"required,oneof=circle oval smile"`
	GlassesStyle string             `bson:"glassesStyle" json:"glassesStyle" validate:"required,oneof=none round square"`
	NoseStyle    string             `bson:"noseStyle" json:"noseStyle" validate:"required,oneof=short long round"`
	MouthStyle   string             `bson:"mouthStyle" json:"mouthStyle" validate:"required,oneof=laugh smile peace"`
	ShirtStyle   string             `bson:"shirtStyle" json:"shirtStyle" validate:"required,oneof=hoody short polo"`
	ShirtColor   string             `bson:"shirtColor" json:"shirtColor" validate:"required,hexcolor"`
	BgColor      string             `bson:"bgColor" json:"bgColor" validate:"required,hexcolor"`
}
