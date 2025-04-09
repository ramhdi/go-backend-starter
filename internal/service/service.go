package service

import (
	"go-backend-starter/internal/repository"
)

// Service handles all business logic
type Service struct {
	repo          repository.Repository
	jwtSecret     string
	jwtExpiration int
}

// NewService creates a new service
func NewService(repo repository.Repository, jwtSecret string, jwtExpiration int) *Service {
	return &Service{
		repo:          repo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}
