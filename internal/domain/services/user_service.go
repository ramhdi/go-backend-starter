package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-backend-starter/internal/domain/models"
	"go-backend-starter/internal/utils"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(pool *pgxpool.Pool) *UserService {
	return &UserService{db: pool}
}

func (s *UserService) CreateUser(ctx context.Context, input *models.CreateUserInput) (*models.User, error) {
	// Check if username or email already exists
	existingUser, err := s.GetUserByUsername(ctx, input.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingUser, err = s.GetUserByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	var user models.User
	err = pgxscan.Get(ctx, s.db, &user, `
		INSERT INTO users (username, password_hash, email, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, password_hash, email, role, created_at, updated_at
	`, input.Username, passwordHash, input.Email, input.Role, time.Now(), time.Now())

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(ctx, s.db, &user, `
		SELECT id, username, password_hash, email, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(ctx, s.db, &user, `
		SELECT id, username, password_hash, email, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`, username)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(ctx, s.db, &user, `
		SELECT id, username, password_hash, email, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, input *models.UpdateUserInput) (*models.User, error) {
	user, err := s.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	setClause := ""
	args := []interface{}{}
	paramCounter := 1

	if input.Username != "" {
		// Check if new username already exists
		if input.Username != user.Username {
			existingUser, err := s.GetUserByUsername(ctx, input.Username)
			if err == nil && existingUser != nil {
				return nil, errors.New("username already exists")
			}
		}
		setClause += fmt.Sprintf("username = $%d, ", paramCounter)
		args = append(args, input.Username)
		paramCounter++
	}

	if input.Email != "" {
		// Check if new email already exists
		if input.Email != user.Email {
			existingUser, err := s.GetUserByEmail(ctx, input.Email)
			if err == nil && existingUser != nil {
				return nil, errors.New("email already exists")
			}
		}
		setClause += fmt.Sprintf("email = $%d, ", paramCounter)
		args = append(args, input.Email)
		paramCounter++
	}

	if input.Password != "" {
		passwordHash, err := utils.HashPassword(input.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		setClause += fmt.Sprintf("password_hash = $%d, ", paramCounter)
		args = append(args, passwordHash)
		paramCounter++
	}

	if input.Role != "" {
		setClause += fmt.Sprintf("role = $%d, ", paramCounter)
		args = append(args, input.Role)
		paramCounter++
	}

	if setClause == "" {
		// No fields to update
		return user, nil
	}

	// Add updated_at field
	setClause += fmt.Sprintf("updated_at = $%d ", paramCounter)
	args = append(args, time.Now())
	paramCounter++

	// Add user ID to arguments
	args = append(args, id)

	// Execute update query
	var updatedUser models.User
	err = pgxscan.Get(ctx, s.db, &updatedUser, fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d
		RETURNING id, username, password_hash, email, role, created_at, updated_at
	`, setClause[:len(setClause)-2], paramCounter), args...)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	_, err := s.db.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) ListUsers(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	err := pgxscan.Select(ctx, s.db, &users, `
		SELECT id, username, password_hash, email, role, created_at, updated_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}
