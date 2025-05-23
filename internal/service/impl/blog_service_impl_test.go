package impl_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/userblog/management/internal/mocks"
	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/service/impl"
)

// TestBlogServiceSuite executes all test cases for BlogService
func TestBlogServiceSuite(t *testing.T) {
	for name, test := range BlogServiceTestBed {
		t.Run(name, test)
	}
}

var BlogServiceTestBed = map[string]func(*testing.T){
	"TestCreate_Success":             TestBlogService_Create_Success,
	"TestCreate_RepositoryError":     TestBlogService_Create_RepositoryError,
	"TestGetByID_Success":            TestBlogService_GetByID_Success,
	"TestGetByID_NotFound":           TestBlogService_GetByID_NotFound,
	"TestUpdate_Success":             TestBlogService_Update_Success,
	"TestUpdate_NotFound":            TestBlogService_Update_NotFound,
	"TestUpdate_Unauthorized":        TestBlogService_Update_Unauthorized,
	"TestDelete_Success":             TestBlogService_Delete_Success,
	"TestDelete_NotFound":            TestBlogService_Delete_NotFound,
	"TestDelete_Unauthorized":        TestBlogService_Delete_Unauthorized,
	"TestList_Success":               TestBlogService_List_Success,
	"TestList_RepositoryError":       TestBlogService_List_RepositoryError,
	"TestListByUser_Success":         TestBlogService_ListByUser_Success,
	"TestListByUser_RepositoryError": TestBlogService_ListByUser_RepositoryError,
}

// TestBlogService_Create_Success tests successful blog creation
func TestBlogService_Create_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide valid data for blog creation
	blog := &models.Blog{
		Title:     "Test Blog",
		Content:   "This is a test blog content",
		Published: false,
	}
	userID := uint(1)

	// Set expectations
	mockRepo.On("Create", mock.MatchedBy(func(b *models.Blog) bool {
		return b.Title == blog.Title && b.Content == blog.Content && b.UserID == userID
	})).Return(nil)

	// Call the service
	err := blogService.Create(blog, userID)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_Create_RepositoryError tests blog creation when repository returns an error
func TestBlogService_Create_RepositoryError(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide valid data for blog creation
	blog := &models.Blog{
		Title:     "Test Blog",
		Content:   "This is a test blog content",
		Published: false,
	}
	userID := uint(1)
	expectedError := errors.New("database error")

	// Set expectations
	mockRepo.On("Create", mock.MatchedBy(func(b *models.Blog) bool {
		return b.Title == blog.Title && b.Content == blog.Content && b.UserID == userID
	})).Return(expectedError)

	// Call the service
	err := blogService.Create(blog, userID)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_GetByID_Success tests successful blog retrieval by ID
func TestBlogService_GetByID_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide a valid blog ID for this test case
	blogID := uint(1)
	expectedBlog := &models.Blog{
		Title:     "Test Blog",
		Content:   "This is a test blog content",
		Published: true,
		UserID:    uint(1),
	}
	expectedBlog.ID = blogID

	// Set expectations
	mockRepo.On("FindByID", blogID).Return(expectedBlog, nil)

	// Call the service
	blog, err := blogService.GetByID(blogID)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedBlog, blog)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_GetByID_NotFound tests blog retrieval when blog is not found
func TestBlogService_GetByID_NotFound(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide a non-existent blog ID for this test case
	blogID := uint(999)
	expectedError := errors.New("record not found")

	// Set expectations
	mockRepo.On("FindByID", blogID).Return(nil, expectedError)

	// Call the service
	blog, err := blogService.GetByID(blogID)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, blog)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_Update_Success tests successful blog update
func TestBlogService_Update_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide valid data for blog update
	blogID := uint(1)
	userID := uint(1)
	updateBlog := &models.Blog{
		Title:     "Updated Blog Title",
		Content:   "Updated content for test blog",
		Published: true,
	}
	updateBlog.ID = blogID

	existingBlog := &models.Blog{
		Title:     "Original Blog Title",
		Content:   "Original content",
		Published: false,
		UserID:    userID,
	}
	existingBlog.ID = blogID

	// Set expectations
	mockRepo.On("FindByID", blogID).Return(existingBlog, nil)
	mockRepo.On("Update", mock.MatchedBy(func(b *models.Blog) bool {
		return b.ID == blogID &&
			b.Title == updateBlog.Title &&
			b.Content == updateBlog.Content &&
			b.Published == updateBlog.Published &&
			b.UserID == userID
	})).Return(nil)

	// Call the service
	err := blogService.Update(updateBlog, userID)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_Update_NotFound tests blog update when blog is not found
func TestBlogService_Update_NotFound(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide a non-existent blog ID for this test case
	blogID := uint(999)
	userID := uint(1)
	updateBlog := &models.Blog{
		Title:     "Updated Blog Title",
		Content:   "Updated content for test blog",
		Published: true,
	}
	updateBlog.ID = blogID

	expectedError := errors.New("record not found")

	// Set expectations
	mockRepo.On("FindByID", blogID).Return(nil, expectedError)

	// Call the service
	err := blogService.Update(updateBlog, userID)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_Update_Unauthorized tests blog update when user is not authorized
func TestBlogService_Update_Unauthorized(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide valid data for blog update with different user IDs
	blogID := uint(1)
	ownerUserID := uint(1)
	differentUserID := uint(2)

	updateBlog := &models.Blog{
		Title:     "Updated Blog Title",
		Content:   "Updated content for test blog",
		Published: true,
	}
	updateBlog.ID = blogID

	existingBlog := &models.Blog{
		Title:     "Original Blog Title",
		Content:   "Original content",
		Published: false,
		UserID:    ownerUserID,
	}
	existingBlog.ID = blogID

	// Set expectations
	mockRepo.On("FindByID", blogID).Return(existingBlog, nil)

	// Call the service with a different user ID
	err := blogService.Update(updateBlog, differentUserID)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	mockRepo.AssertExpectations(t)
}

// TestBlogService_Delete_Success tests successful blog deletion
func TestBlogService_Delete_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide valid blog ID for deletion
	blogID := uint(1)
	userID := uint(1)

	existingBlog := &models.Blog{
		Title:     "Blog to Delete",
		Content:   "Content of blog to delete",
		Published: true,
		UserID:    userID,
	}
	existingBlog.ID = blogID

	// Set expectations
	mockRepo.On("FindByID", blogID).Return(existingBlog, nil)
	mockRepo.On("Delete", blogID).Return(nil)

	// Call the service
	err := blogService.Delete(blogID, userID)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_Delete_NotFound tests blog deletion when blog is not found
func TestBlogService_Delete_NotFound(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide a non-existent blog ID for this test case
	blogID := uint(999)
	userID := uint(1)

	expectedError := errors.New("record not found")

	// Set expectations
	mockRepo.On("FindByID", blogID).Return(nil, expectedError)

	// Call the service
	err := blogService.Delete(blogID, userID)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_Delete_Unauthorized tests blog deletion when user is not authorized
func TestBlogService_Delete_Unauthorized(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide valid data for blog deletion with different user IDs
	blogID := uint(1)
	ownerUserID := uint(1)
	differentUserID := uint(2)

	existingBlog := &models.Blog{
		Title:     "Blog to Delete",
		Content:   "Content of blog to delete",
		Published: true,
		UserID:    ownerUserID,
	}
	existingBlog.ID = blogID

	// Set expectations
	mockRepo.On("FindByID", blogID).Return(existingBlog, nil)

	// Call the service with a different user ID
	err := blogService.Delete(blogID, differentUserID)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	mockRepo.AssertExpectations(t)
}

// TestBlogService_List_Success tests successful blog listing
func TestBlogService_List_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide valid data for blog listing
	page := 1
	perPage := 10
	publishedOnly := true
	offset := (page - 1) * perPage

	expectedBlogs := []models.Blog{
		{
			Title:     "Blog 1",
			Content:   "Content of blog 1",
			Published: true,
			UserID:    uint(1),
		},
		{
			Title:     "Blog 2",
			Content:   "Content of blog 2",
			Published: true,
			UserID:    uint(2),
		},
	}
	expectedCount := 2

	// Set expectations
	mockRepo.On("List", offset, perPage, publishedOnly).Return(expectedBlogs, expectedCount, nil)

	// Call the service
	blogs, count, err := blogService.List(page, perPage, publishedOnly)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedBlogs, blogs)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_List_RepositoryError tests blog listing when repository returns an error
func TestBlogService_List_RepositoryError(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide valid data for blog listing
	page := 1
	perPage := 10
	publishedOnly := true
	offset := (page - 1) * perPage

	expectedError := errors.New("database error")

	// Set expectations
	mockRepo.On("List", offset, perPage, publishedOnly).Return([]models.Blog{}, 0, expectedError)

	// Call the service
	blogs, count, err := blogService.List(page, perPage, publishedOnly)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, blogs)
	assert.Equal(t, 0, count)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_ListByUser_Success tests successful blog listing by user
func TestBlogService_ListByUser_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide valid data for blog listing by user
	userID := uint(1)
	page := 1
	perPage := 10
	offset := (page - 1) * perPage

	expectedBlogs := []models.Blog{
		{
			Title:     "User Blog 1",
			Content:   "Content of user blog 1",
			Published: true,
			UserID:    userID,
		},
		{
			Title:     "User Blog 2",
			Content:   "Content of user blog 2",
			Published: false,
			UserID:    userID,
		},
	}
	expectedCount := 2

	// Set expectations
	mockRepo.On("ListByUser", userID, offset, perPage).Return(expectedBlogs, expectedCount, nil)

	// Call the service
	blogs, count, err := blogService.ListByUser(userID, page, perPage)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedBlogs, blogs)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}

// TestBlogService_ListByUser_RepositoryError tests blog listing by user when repository returns an error
func TestBlogService_ListByUser_RepositoryError(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockBlogRepository)

	// Create the service with mock repository
	blogService := impl.NewBlogService(mockRepo)

	// Test data
	// TODO: Provide valid data for blog listing by user
	userID := uint(1)
	page := 1
	perPage := 10
	offset := (page - 1) * perPage

	expectedError := errors.New("database error")

	// Set expectations
	mockRepo.On("ListByUser", userID, offset, perPage).Return([]models.Blog{}, 0, expectedError)

	// Call the service
	blogs, count, err := blogService.ListByUser(userID, page, perPage)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, blogs)
	assert.Equal(t, 0, count)
	mockRepo.AssertExpectations(t)
}
