package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-backend-starter/internal/models"
	"go-backend-starter/internal/utils"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository implements Repository interface for PostgreSQL
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetUserByID retrieves a user by ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(ctx, r.db, &user, `
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

// GetUserByUsername retrieves a user by username
func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(ctx, r.db, &user, `
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

// GetUserByEmail retrieves a user by email
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(ctx, r.db, &user, `
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

// CreateUser creates a new user
func (r *PostgresRepository) CreateUser(ctx context.Context, input *models.CreateUserInput) (*models.User, error) {
	// Hash password
	passwordHash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	var user models.User
	err = pgxscan.Get(ctx, r.db, &user, `
		INSERT INTO users (username, password_hash, email, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, password_hash, email, role, created_at, updated_at
	`, input.Username, passwordHash, input.Email, input.Role, time.Now(), time.Now())

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// UpdateUser updates an existing user
func (r *PostgresRepository) UpdateUser(ctx context.Context, id int, input *models.UpdateUserInput) (*models.User, error) {
	// Build the set clause and arguments for the SQL query
	setClauses := []string{}
	args := []interface{}{}
	paramCounter := 1

	if input.Username != "" {
		setClauses = append(setClauses, fmt.Sprintf("username = $%d", paramCounter))
		args = append(args, input.Username)
		paramCounter++
	}

	if input.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", paramCounter))
		args = append(args, input.Email)
		paramCounter++
	}

	if input.Password != "" {
		passwordHash, err := utils.HashPassword(input.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("password_hash = $%d", paramCounter))
		args = append(args, passwordHash)
		paramCounter++
	}

	if input.Role != "" {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", paramCounter))
		args = append(args, input.Role)
		paramCounter++
	}

	// If no fields to update, just return the current user
	if len(setClauses) == 0 {
		return r.GetUserByID(ctx, id)
	}

	// Add updated_at field
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", paramCounter))
	args = append(args, time.Now())
	paramCounter++

	// Add user ID to arguments
	args = append(args, id)

	// Join set clauses with commas
	setClause := strings.Join(setClauses, ", ")

	// Execute update query
	var updatedUser models.User
	err := pgxscan.Get(ctx, r.db, &updatedUser, fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d
		RETURNING id, username, password_hash, email, role, created_at, updated_at
	`, setClause, paramCounter), args...)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &updatedUser, nil
}

// DeleteUser deletes a user
func (r *PostgresRepository) DeleteUser(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers retrieves a list of users with pagination
func (r *PostgresRepository) ListUsers(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	err := pgxscan.Select(ctx, r.db, &users, `
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
