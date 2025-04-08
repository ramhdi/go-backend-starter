package services

import (
	"context"
	"errors"
	"fmt"

	"go-backend-starter/internal/domain/models"
	"go-backend-starter/internal/utils"
)

type AuthService struct {
	userService   *UserService
	jwtSecret     string
	jwtExpiration int
}

func NewAuthService(userService *UserService, jwtSecret string, jwtExpiration int) *AuthService {
	return &AuthService{
		userService:   userService,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

func (s *AuthService) Login(ctx context.Context, input *models.LoginInput) (string, error) {
	user, err := s.userService.GetUserByUsername(ctx, input.Username)
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

func (s *AuthService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	claims, err := utils.ValidateJWT(tokenString, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return claims, nil
}
