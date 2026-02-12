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
	Title       string               `bson:"title" json:"title"`
	User        PostUser           `bson:"user" json:"user"`
	AuthorID primitive.ObjectID `bson:"author_id" json:"author_id"`
	Slug         string                  `bson:"slug" json:"slug"`
	Caption     string             `bson:"caption" json:"caption"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
	Published    bool                   `bson:"published" json:"published"`
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

