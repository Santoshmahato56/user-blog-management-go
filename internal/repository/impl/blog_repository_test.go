package impl_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/repository/impl"
)

// TestBlogRepositorySuite executes all test cases for BlogRepository
func TestBlogRepositorySuite(t *testing.T) {
	for name, test := range BlogRepositoryTestBed {
		t.Run(name, test)
	}
}

var BlogRepositoryTestBed = map[string]func(*testing.T){
	"TestCreate_Success":     TestBlogRepository_Create_Success,
	"TestCreate_Error":       TestBlogRepository_Create_Error,
	"TestFindByID_Success":   TestBlogRepository_FindByID_Success,
	"TestFindByID_NotFound":  TestBlogRepository_FindByID_NotFound,
	"TestUpdate_Success":     TestBlogRepository_Update_Success,
	"TestUpdate_Error":       TestBlogRepository_Update_Error,
	"TestDelete_Success":     TestBlogRepository_Delete_Success,
	"TestDelete_Error":       TestBlogRepository_Delete_Error,
	"TestList_Success":       TestBlogRepository_List_Success,
	"TestList_Error":         TestBlogRepository_List_Error,
	"TestListByUser_Success": TestBlogRepository_ListByUser_Success,
	"TestListByUser_Error":   TestBlogRepository_ListByUser_Error,
}

// setupMockDB creates a mock database connection
func setupBlogMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
	// Create a sqlmock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// Create a gorm DB instance which uses the mock DB
	gormDB, err := gorm.Open("mysql", db)
	require.NoError(t, err)

	// Set logger
	gormDB.LogMode(true)

	return gormDB, mock, err
}

// TestBlogRepository_Create_Success tests successful blog creation
func TestBlogRepository_Create_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide valid blog data for creation
	blog := &models.Blog{
		Title:     "Test Blog",
		Content:   "This is a test blog content",
		Published: false,
		UserID:    uint(1),
	}

	// Set expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blogs`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			blog.Title,
			blog.Content,
			blog.Published,
			blog.UserID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the repository
	err = repo.Create(blog)

	// Assert expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_Create_Error tests blog creation when database returns an error
func TestBlogRepository_Create_Error(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide valid blog data for creation
	blog := &models.Blog{
		Title:     "Test Blog",
		Content:   "This is a test blog content",
		Published: false,
		UserID:    uint(1),
	}

	// Mock database error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blogs`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			blog.Title,
			blog.Content,
			blog.Published,
			blog.UserID,
		).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	// Call the repository
	err = repo.Create(blog)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_FindByID_Success tests successful blog retrieval by ID
func TestBlogRepository_FindByID_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide a valid blog ID for this test case
	blogID := uint(1)

	// Mock blog data that would be returned from DB
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "content", "published", "user_id"}).
		AddRow(blogID, "2023-01-01", "2023-01-01", nil, "Test Blog", "This is a test blog content", true, 1)

	// User data for preload
	userRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "first_name", "last_name", "role_id"}).
		AddRow(1, "2023-01-01", "2023-01-01", nil, "testuser", "test@example.com", "hashedpassword", "Test", "User", 1)

	// Set expectations
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blogs` WHERE `blogs`.`id` = ? AND `blogs`.`deleted_at` IS NULL ORDER BY `blogs`.`id` ASC LIMIT 1")).
		WithArgs(blogID).
		WillReturnRows(rows)

	// Expect preload query for User
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).
		WithArgs(1).
		WillReturnRows(userRows)

	// Call the repository
	blog, err := repo.FindByID(blogID)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, blog)
	assert.Equal(t, blogID, blog.ID)
	assert.Equal(t, "Test Blog", blog.Title)
	assert.Equal(t, "This is a test blog content", blog.Content)
	assert.True(t, blog.Published)
	assert.Equal(t, uint(1), blog.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_FindByID_NotFound tests blog retrieval when blog is not found
func TestBlogRepository_FindByID_NotFound(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide a non-existent blog ID for this test case
	blogID := uint(999)

	// Set expectations for not found
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blogs` WHERE `blogs`.`id` = ? AND `blogs`.`deleted_at` IS NULL ORDER BY `blogs`.`id` ASC LIMIT 1")).
		WithArgs(blogID).
		WillReturnError(gorm.ErrRecordNotFound)

	// Call the repository
	blog, err := repo.FindByID(blogID)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NotNil(t, blog) // GORM returns an empty struct, not nil
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_Update_Success tests successful blog update
func TestBlogRepository_Update_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide valid blog data for update
	blog := &models.Blog{
		Title:     "Updated Blog Title",
		Content:   "Updated content for test blog",
		Published: true,
		UserID:    uint(1),
	}
	blog.ID = uint(1)

	// Set expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `blogs` SET")).
		WithArgs(
			sqlmock.AnyArg(), // UpdatedAt
			blog.Title,
			blog.Content,
			blog.Published,
			blog.UserID,
			blog.ID,
			sqlmock.AnyArg(), // DeletedAt
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the repository
	err = repo.Update(blog)

	// Assert expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_Update_Error tests blog update when database returns an error
func TestBlogRepository_Update_Error(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide valid blog data for update
	blog := &models.Blog{
		Title:     "Updated Blog Title",
		Content:   "Updated content for test blog",
		Published: true,
		UserID:    uint(1),
	}
	blog.ID = uint(1)

	// Mock database error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `blogs` SET")).
		WithArgs(
			sqlmock.AnyArg(), // UpdatedAt
			blog.Title,
			blog.Content,
			blog.Published,
			blog.UserID,
			blog.ID,
			sqlmock.AnyArg(), // DeletedAt
		).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	// Call the repository
	err = repo.Update(blog)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_Delete_Success tests successful blog deletion
func TestBlogRepository_Delete_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide a valid blog ID for deletion
	blogID := uint(1)

	// Set expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `blogs` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			blogID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the repository
	err = repo.Delete(blogID)

	// Assert expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_Delete_Error tests blog deletion when database returns an error
func TestBlogRepository_Delete_Error(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide a valid blog ID for deletion
	blogID := uint(1)

	// Mock database error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `blogs` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			blogID,
		).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	// Call the repository
	err = repo.Delete(blogID)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_List_Success tests successful blog listing
func TestBlogRepository_List_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide valid data for blog listing
	offset := 0
	limit := 10
	published := true

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// Mock blogs data
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "content", "published", "user_id"}).
		AddRow(1, "2023-01-01", "2023-01-01", nil, "Blog 1", "Content of blog 1", true, 1).
		AddRow(2, "2023-01-02", "2023-01-02", nil, "Blog 2", "Content of blog 2", true, 2)

	// Set expectations for count
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `blogs` WHERE `blogs`.`deleted_at` IS NULL AND ((published = ?))")).
		WithArgs(published).
		WillReturnRows(countRows)

	// Set expectations for getting blogs
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blogs` WHERE `blogs`.`deleted_at` IS NULL AND ((published = ?)) LIMIT 10 OFFSET 0")).
		WithArgs(published).
		WillReturnRows(rows)

	// Expect preload queries for User
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` IN (?,?) AND `users`.`deleted_at` IS NULL")).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "first_name", "last_name", "role_id"}).
			AddRow(1, "2023-01-01", "2023-01-01", nil, "user1", "user1@example.com", "hashedpassword", "User", "One", 1).
			AddRow(2, "2023-01-01", "2023-01-01", nil, "user2", "user2@example.com", "hashedpassword", "User", "Two", 1))

	// Call the repository
	blogs, count, err := repo.List(offset, limit, published)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, blogs, 2)
	assert.Equal(t, "Blog 1", blogs[0].Title)
	assert.Equal(t, "Blog 2", blogs[1].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_List_Error tests blog listing when database returns an error
func TestBlogRepository_List_Error(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide valid data for blog listing
	offset := 0
	limit := 10
	published := true

	// Mock database error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `blogs` WHERE `blogs`.`deleted_at` IS NULL AND ((published = ?))")).
		WithArgs(published).
		WillReturnError(errors.New("database error"))

	// Call the repository
	blogs, count, err := repo.List(offset, limit, published)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.Empty(t, blogs)
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_ListByUser_Success tests successful blog listing by user
func TestBlogRepository_ListByUser_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide valid data for blog listing by user
	userID := uint(1)
	offset := 0
	limit := 10

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// Mock blogs data
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "content", "published", "user_id"}).
		AddRow(1, "2023-01-01", "2023-01-01", nil, "User Blog 1", "Content of user blog 1", true, userID).
		AddRow(3, "2023-01-03", "2023-01-03", nil, "User Blog 2", "Content of user blog 2", false, userID)

	// User data for preload
	userRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "first_name", "last_name", "role_id"}).
		AddRow(userID, "2023-01-01", "2023-01-01", nil, "testuser", "test@example.com", "hashedpassword", "Test", "User", 1)

	// Set expectations for count
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `blogs` WHERE `blogs`.`deleted_at` IS NULL AND ((user_id = ?))")).
		WithArgs(userID).
		WillReturnRows(countRows)

	// Set expectations for getting blogs
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blogs` WHERE `blogs`.`deleted_at` IS NULL AND ((user_id = ?)) LIMIT 10 OFFSET 0")).
		WithArgs(userID).
		WillReturnRows(rows)

	// Expect preload query for User
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL")).
		WithArgs(userID).
		WillReturnRows(userRows)

	// Call the repository
	blogs, count, err := repo.ListByUser(userID, offset, limit)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, blogs, 2)
	assert.Equal(t, "User Blog 1", blogs[0].Title)
	assert.Equal(t, "User Blog 2", blogs[1].Title)
	assert.Equal(t, userID, blogs[0].UserID)
	assert.Equal(t, userID, blogs[1].UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBlogRepository_ListByUser_Error tests blog listing by user when database returns an error
func TestBlogRepository_ListByUser_Error(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupBlogMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewBlogRepository(mockDB)

	// Test data
	// TODO: Provide valid data for blog listing by user
	userID := uint(1)
	offset := 0
	limit := 10

	// Mock database error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `blogs` WHERE `blogs`.`deleted_at` IS NULL AND ((user_id = ?))")).
		WithArgs(userID).
		WillReturnError(errors.New("database error"))

	// Call the repository
	blogs, count, err := repo.ListByUser(userID, offset, limit)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.Empty(t, blogs)
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
