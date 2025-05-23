package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/userblog/management/internal/models"
)

// MockUserRepository is a mock implementation of IUserRepository
type MockUserRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// FindByUsername mocks the FindByUsername method
func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// FindByEmail mocks the FindByEmail method
func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Update mocks the Update method
func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockUserRepository) List(offset, limit int) ([]models.User, int, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]models.User), args.Int(1), args.Error(2)
}
