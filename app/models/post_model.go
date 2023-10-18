package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Visibility int8

const (
	PUBLIC Visibility = iota
	SELF
	FOLLOWERS
)

type Post struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id"`
	Creator        primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
	PostVisibility Visibility         `bson:"visibility" json:"visibility"`
	Content        string             `bson:"content" json:"content"`
	Image          string             `bson:"image_link" json:"image_link"`
}
