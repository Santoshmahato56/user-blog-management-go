package service

import "github.com/userblog/management/internal/models"

// IUserService defines the interface for user operations
type IUserService interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(page, perPage int) ([]models.User, int, error)
}
