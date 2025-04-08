package handlers

import (
	"net/http"

	"go-backend-starter/internal/domain/models"
	"go-backend-starter/internal/domain/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), &input)
	if err != nil {
		log.Error().Err(err).Str("username", input.Username).Msg("Login failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
