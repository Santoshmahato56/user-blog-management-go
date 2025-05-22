package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/userblog/management/internal/models"
)

// IUserRepository defines the interface for user database operations
type IUserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(offset, limit int) ([]models.User, int, error)
}

// UserRepository handles all database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository with the given database connection
func NewUserRepository(database *gorm.DB) *UserRepository {
	return &UserRepository{
		db: database,
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role.Permissions").First(&user, id).Error
	return &user, err
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role.Permissions").Where("username = ?", username).First(&user).Error
	return &user, err
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role.Permissions").Where("email = ?", email).First(&user).Error
	return &user, err
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List returns a list of users with pagination
func (r *UserRepository) List(offset, limit int) ([]models.User, int, error) {
	var users []models.User
	var count int

	// Get the total count
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get the users with pagination
	err := r.db.Preload("Role").Offset(offset).Limit(limit).Find(&users).Error
	return users, count, err
}
