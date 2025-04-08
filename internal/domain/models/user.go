package models

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Don't expose password hash
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

type UpdateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password" binding:"omitempty,min=8"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role" binding:"omitempty,oneof=admin user"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
