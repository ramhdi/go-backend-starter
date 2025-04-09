package routes

import (
	"go-backend-starter/internal/api/handlers"
	"go-backend-starter/internal/api/middleware"
	"go-backend-starter/internal/service"

	"github.com/gin-gonic/gin"
)

// Setup configures all API routes
func Setup(router *gin.Engine, handler *handlers.Handler, service *service.Service) {
	// Apply global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CorsMiddleware())

	// Health check
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	api := router.Group("/api")
	{
		// Auth routes
		api.POST("/auth/login", handler.Login)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(service))
	{
		// User routes - admin only
		users := protected.Group("/users")
		users.Use(middleware.RequireRole("admin"))
		{
			users.POST("", handler.CreateUser)
			users.GET("", handler.ListUsers)
			users.GET("/:id", handler.GetUser)
			users.PUT("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
		}

		// Current user route - for any authenticated user
		protected.GET("/me", handler.GetCurrentUser)
	}
}
