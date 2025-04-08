package routes

import (
	"net/http"
	"strconv"

	"go-backend-starter/internal/api/handlers"
	"go-backend-starter/internal/api/middleware"
	"go-backend-starter/internal/domain/services"

	"github.com/gin-gonic/gin"
)

func Setup(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authService *services.AuthService,
) {
	// Apply global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CorsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// API routes that require authentication
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(authService))
	{
		// User routes - admin only
		users := api.Group("/users")
		users.Use(middleware.RequireRole("admin"))
		{
			users.POST("", userHandler.Create)
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.Get)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}

		// Current user route - all authenticated users
		api.GET("/me", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.Redirect(http.StatusTemporaryRedirect, "/api/users/"+strconv.Itoa(userID.(int)))
		})
	}
}
