package service

import (
	"context"
	"errors"
	"fmt"
	"go-backend-starter/internal/models"
)

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// CreateUser creates a new user with validation
func (s *Service) CreateUser(ctx context.Context, input *models.CreateUserInput) (*models.User, error) {
	// Check if username already exists
	existingUser, err := s.repo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	return s.repo.CreateUser(ctx, input)
}

// UpdateUser updates an existing user with validation
func (s *Service) UpdateUser(ctx context.Context, id int, input *models.UpdateUserInput) (*models.User, error) {
	// Check if user exists
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Validate username uniqueness if changed
	if input.Username != "" && input.Username != user.Username {
		existingUser, err := s.repo.GetUserByUsername(ctx, input.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username: %w", err)
		}
		if existingUser != nil {
			return nil, errors.New("username already exists")
		}
	}

	// Validate email uniqueness if changed
	if input.Email != "" && input.Email != user.Email {
		existingUser, err := s.repo.GetUserByEmail(ctx, input.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if existingUser != nil {
			return nil, errors.New("email already exists")
		}
	}

	return s.repo.UpdateUser(ctx, id, input)
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

// ListUsers retrieves a list of users with pagination
func (s *Service) ListUsers(ctx context.Context, offset, limit int) ([]*models.User, error) {
	return s.repo.ListUsers(ctx, offset, limit)
}
