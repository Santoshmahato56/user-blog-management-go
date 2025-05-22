package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/userblog/management/internal/models"
)

// IBlogRepository defines the interface for blog database operations
type IBlogRepository interface {
	Create(blog *models.Blog) error
	FindByID(id uint) (*models.Blog, error)
	Update(blog *models.Blog) error
	Delete(id uint) error
	List(offset, limit int, published bool) ([]models.Blog, int, error)
	ListByUser(userID uint, offset, limit int) ([]models.Blog, int, error)
}

// BlogRepository handles all database operations for blogs
type BlogRepository struct {
	db *gorm.DB
}

// NewBlogRepository creates a new blog repository with the given database connection
func NewBlogRepository(database *gorm.DB) *BlogRepository {
	return &BlogRepository{
		db: database,
	}
}

// Create creates a new blog
func (r *BlogRepository) Create(blog *models.Blog) error {
	return r.db.Create(blog).Error
}

// FindByID finds a blog by ID
func (r *BlogRepository) FindByID(id uint) (*models.Blog, error) {
	var blog models.Blog
	err := r.db.Preload("User").First(&blog, id).Error
	return &blog, err
}

// Update updates a blog
func (r *BlogRepository) Update(blog *models.Blog) error {
	return r.db.Save(blog).Error
}

// Delete deletes a blog
func (r *BlogRepository) Delete(id uint) error {
	return r.db.Delete(&models.Blog{}, id).Error
}

// List returns a list of blogs with pagination
func (r *BlogRepository) List(offset, limit int, published bool) ([]models.Blog, int, error) {
	var blogs []models.Blog
	var count int

	query := r.db.Model(&models.Blog{})
	if published {
		query = query.Where("published = ?", true)
	}

	// Get the total count
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get the blogs with pagination
	err := query.Preload("User").Offset(offset).Limit(limit).Find(&blogs).Error
	return blogs, count, err
}

// ListByUser returns a list of blogs by user with pagination
func (r *BlogRepository) ListByUser(userID uint, offset, limit int) ([]models.Blog, int, error) {
	var blogs []models.Blog
	var count int

	// Get the total count
	if err := r.db.Model(&models.Blog{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get the blogs with pagination
	err := r.db.Preload("User").Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&blogs).Error
	return blogs, count, err
}
