package entity

import (
	"time"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserSignUp struct {
	Email    string `json:"email" form:"Email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" form:"Password" binding:"required,min=6" example:"password123"`
	Name     string `json:"name" form:"Name" binding:"required" example:"John Doe"`
}

type UserLogin struct {
	Email    string `json:"email" form:"Email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" form:"Password" binding:"required" example:"password123"`
}

type UserResponse struct {
	ID    int64  `json:"id" example:"1"`
	Email string `json:"email" example:"user@example.com"`
	Name  string `json:"name" example:"John Doe"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsI..."`
}