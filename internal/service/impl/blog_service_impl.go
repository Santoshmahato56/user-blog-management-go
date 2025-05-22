package impl

import (
	"errors"

	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/repository"
	"github.com/userblog/management/internal/service"
)

// BlogService implements the IBlogService interface
type BlogService struct {
	blogRepo repository.IBlogRepository
}

// NewBlogService creates a new blog service
func NewBlogService(blogRepo repository.IBlogRepository) service.IBlogService {
	return &BlogService{
		blogRepo: blogRepo,
	}
}

// Create creates a new blog
func (s *BlogService) Create(blog *models.Blog, userID uint) error {
	blog.UserID = userID
	return s.blogRepo.Create(blog)
}

// GetByID returns a blog by ID
func (s *BlogService) GetByID(id uint) (*models.Blog, error) {
	return s.blogRepo.FindByID(id)
}

// Update updates a blog
func (s *BlogService) Update(blog *models.Blog, userID uint) error {
	// Get the existing blog
	existingBlog, err := s.blogRepo.FindByID(blog.ID)
	if err != nil {
		return err
	}

	// Check if the user owns the blog
	if existingBlog.UserID != userID {
		return errors.New("unauthorized: you can only update your own blogs")
	}

	// Update only allowed fields
	existingBlog.Title = blog.Title
	existingBlog.Content = blog.Content
	existingBlog.Published = blog.Published

	return s.blogRepo.Update(existingBlog)
}

// Delete deletes a blog
func (s *BlogService) Delete(id uint, userID uint) error {
	// Get the existing blog
	existingBlog, err := s.blogRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check if the user owns the blog
	if existingBlog.UserID != userID {
		return errors.New("unauthorized: you can only delete your own blogs")
	}

	return s.blogRepo.Delete(id)
}

// List returns a list of blogs with pagination
func (s *BlogService) List(page, perPage int, publishedOnly bool) ([]models.Blog, int, error) {
	offset := (page - 1) * perPage
	return s.blogRepo.List(offset, perPage, publishedOnly)
}

// ListByUser returns a list of blogs by user with pagination
func (s *BlogService) ListByUser(userID uint, page, perPage int) ([]models.Blog, int, error) {
	offset := (page - 1) * perPage
	return s.blogRepo.ListByUser(userID, offset, perPage)
}
