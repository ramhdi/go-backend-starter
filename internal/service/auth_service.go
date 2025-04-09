package service

import (
	"context"
	"errors"
	"fmt"
	"go-backend-starter/internal/models"
	"go-backend-starter/internal/utils"
)

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, input *models.LoginInput) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, input.Username)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return "", errors.New("invalid username or password")
	}

	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		return "", errors.New("invalid username or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *Service) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	claims, err := utils.ValidateJWT(tokenString, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return claims, nil
}
