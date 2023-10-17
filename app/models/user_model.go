package models

import (
	"time"

	"github.com/google/uuid"
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
	Id           uuid.UUID   `db:"id" json:"id"`
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updated_at"`
	Email        string      `db:"email" json:"email"`
	Username     string      `db:"username" json:"username"`
	PasswordHash string      `db:"password_hash" json:"password_hash"`
	ProfileImage string      `db:"profile_image" json:"profile_image"`
	UserStatus   Status      `db:"status" json:"status"`
	Followers    []uuid.UUID `db:"followers" json:"followers"`
	Following    []uuid.UUID `db:"following" json:"following"`
}
