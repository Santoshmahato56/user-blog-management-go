package service

import "github.com/userblog/management/internal/models"

// IBlogService defines the interface for blog operations
type IBlogService interface {
	Create(blog *models.Blog, userID uint) error
	GetByID(id uint) (*models.Blog, error)
	Update(blog *models.Blog, userID uint) error
	Delete(id uint, userID uint) error
	List(page, perPage int, publishedOnly bool) ([]models.Blog, int, error)
	ListByUser(userID uint, page, perPage int) ([]models.Blog, int, error)
}
