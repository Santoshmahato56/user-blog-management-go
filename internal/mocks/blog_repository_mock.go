package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/userblog/management/internal/models"
)

// MockBlogRepository is a mock implementation of IBlogRepository
type MockBlogRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockBlogRepository) Create(blog *models.Blog) error {
	args := m.Called(blog)
	return args.Error(0)
}

// FindByID mocks the FindByID method
func (m *MockBlogRepository) FindByID(id uint) (*models.Blog, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Blog), args.Error(1)
}

// Update mocks the Update method
func (m *MockBlogRepository) Update(blog *models.Blog) error {
	args := m.Called(blog)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockBlogRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockBlogRepository) List(offset, limit int, published bool) ([]models.Blog, int, error) {
	args := m.Called(offset, limit, published)
	return args.Get(0).([]models.Blog), args.Int(1), args.Error(2)
}

// ListByUser mocks the ListByUser method
func (m *MockBlogRepository) ListByUser(userID uint, offset, limit int) ([]models.Blog, int, error) {
	args := m.Called(userID, offset, limit)
	return args.Get(0).([]models.Blog), args.Int(1), args.Error(2)
}
