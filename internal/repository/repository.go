package repository

import (
	"context"

	"go-backend-starter/internal/models"
)

// Repository defines all data access operations
type Repository interface {
	// User operations
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, input *models.CreateUserInput) (*models.User, error)
	UpdateUser(ctx context.Context, id int, input *models.UpdateUserInput) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, offset, limit int) ([]*models.User, error)
}
