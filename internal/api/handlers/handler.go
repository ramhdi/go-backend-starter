package handlers

import (
	"go-backend-starter/internal/service"
)

// Handler manages HTTP requests
type Handler struct {
	service *service.Service
}

// NewHandler creates a new handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}
