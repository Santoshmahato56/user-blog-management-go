package impl

import (
	"errors"

	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/repository"
	"github.com/userblog/management/internal/service"
)

// UserService implements the IUserService interface
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository) service.IUserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Create creates a new user
func (s *UserService) Create(user *models.User) error {
	// Check if username already exists
	existingUser, err := s.userRepo.FindByUsername(user.Username)
	if err == nil && existingUser.ID != 0 {
		return errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser.ID != 0 {
		return errors.New("email already exists")
	}

	return s.userRepo.Create(user)
}

// GetByID returns a user by ID
func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// Update updates a user
func (s *UserService) Update(user *models.User) error {
	// Get the existing user
	existingUser, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return err
	}

	// Check if username is being changed and if it already exists
	if user.Username != existingUser.Username {
		newUser, err := s.userRepo.FindByUsername(user.Username)
		if err == nil && newUser.ID != 0 && newUser.ID != user.ID {
			return errors.New("username already exists")
		}
	}

	// Check if email is being changed and if it already exists
	if user.Email != existingUser.Email {
		newUser, err := s.userRepo.FindByEmail(user.Email)
		if err == nil && newUser.ID != 0 && newUser.ID != user.ID {
			return errors.New("email already exists")
		}
	}

	// Update only allowed fields
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName

	// Only update password if provided
	if user.Password != "" {
		existingUser.Password = user.Password
	}

	// Only admin can change roles
	if user.RoleID != 0 {
		existingUser.RoleID = user.RoleID
	}

	return s.userRepo.Update(existingUser)
}

// Delete deletes a user
func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

// List returns a list of users with pagination
func (s *UserService) List(page, perPage int) ([]models.User, int, error) {
	offset := (page - 1) * perPage
	return s.userRepo.List(offset, perPage)
}
