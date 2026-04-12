package models

import "time"

type User struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	Username     string    `bson:"username" json:"username"`
	PasswordHash string    `bson:"passwordHash" json:"-"`
	Role         string    `bson:"role" json:"role"` // ADMIN or USER
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
}

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"` // Defaults to USER if empty or invalid
}

type SigninRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
