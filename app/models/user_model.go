package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status int8

const (
	OFFLINE Status = iota
	PRIVATE
	ONLINE
	DEACTIVATED
	DELETED
)

type User struct {
	ID           primitive.ObjectID   `bson:"_id" json:"_id"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at" json:"updated_at"`
	Name         string               `bson:"name" json:"name"`
	Email        string               `bson:"email" json:"email"`
	Username     string               `bson:"username" json:"username"`
	PasswordHash string               `bson:"password_hash" json:"password_hash"`
	ProfileImage string               `bson:"profile_image" json:"profile_image"`
	UserStatus   Status               `bson:"status" json:"status"`
	Followers    []primitive.ObjectID `bson:"followers" json:"followers"`
	Following    []primitive.ObjectID `bson:"following" json:"following"`
}
