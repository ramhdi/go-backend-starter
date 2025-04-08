package handlers

import (
	"net/http"
	"strconv"

	"go-backend-starter/internal/domain/models"
	"go-backend-starter/internal/domain/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &input)
	if err != nil {
		log.Error().Err(err).Interface("input", input).Msg("Create user failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Get user failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, &input)
	if err != nil {
		log.Error().Err(err).Int("id", id).Interface("input", input).Msg("Update user failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		log.Error().Err(err).Int("id", id).Msg("Delete user failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) List(c *gin.Context) {
	// Default pagination values
	offset := 0
	limit := 10

	// Parse offset and limit from query parameters
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	if o, err := strconv.Atoi(offsetStr); err == nil {
		offset = o
	}

	if l, err := strconv.Atoi(limitStr); err == nil {
		limit = l
	}

	users, err := h.userService.ListUsers(c.Request.Context(), offset, limit)
	if err != nil {
		log.Error().Err(err).Int("offset", offset).Int("limit", limit).Msg("List users failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
