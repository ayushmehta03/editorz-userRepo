package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type PostUser struct {
	ID           primitive.ObjectID `bson:"id" json:"id"`
	Name         string             `bson:"name" json:"name"`
	ProfileImage string             `bson:"profileImage" json:"profileImage"`
}

type Post struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	User        PostUser           `bson:"user" json:"user"`
	Caption     string             `bson:"caption" json:"caption"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
	Suggestions []Suggestion       `bson:"suggestions" json:"suggestions"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Suggestion struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	User      PostUser          `bson:"user" json:"user"`
	Text      string            `bson:"text" json:"text"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
}
