package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/service"
	"net/http"
)

// IAuthMiddleware defines the interface for authentication middleware
type IAuthMiddleware interface {
	JWTAuth() gin.HandlerFunc
	RequirePermission(resource, action string) gin.HandlerFunc
}

// AuthMiddleware implements the IAuthMiddleware interface
type AuthMiddleware struct {
	authService service.IAuthService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService service.IAuthService) IAuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// JWTAuth middleware for JWT authentication
func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")

		// Extract token from header
		tokenString, err := m.authService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Validate token and get user
		user, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Set user in the context
		c.Set("user", *user)
		c.Next()
	}
}

// RequirePermission middleware to check if the user has the required permission
func (m *AuthMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			c.Abort()
			return
		}

		user, ok := userInterface.(models.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast user from context"})
			c.Abort()
			return
		}

		// Check if user has the required permission
		if !m.authService.ValidatePermission(&user, resource, action) {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Permission denied: %s:%s", resource, action)})
			c.Abort()
			return
		}

		c.Next()
	}
}
