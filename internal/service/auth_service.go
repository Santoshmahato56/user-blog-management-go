package service

import "github.com/userblog/management/internal/models"

// IAuthService defines the interface for authentication operations
type IAuthService interface {
	Register(user *models.User) error
	Login(username, password string) (string, error)
	GetUserByID(id uint) (*models.User, error)
	ValidateToken(tokenString string) (*models.User, error)
	ValidatePermission(user *models.User, resource, action string) bool
	ExtractTokenFromHeader(authHeader string) (string, error)
}
