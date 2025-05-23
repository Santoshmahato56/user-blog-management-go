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

// TestUserRepositorySuite executes all test cases for UserRepository
func TestUserRepositorySuite(t *testing.T) {
	for name, test := range UserRepositoryTestBed {
		t.Run(name, test)
	}
}

var UserRepositoryTestBed = map[string]func(*testing.T){
	"TestCreate_Success":          TestUserRepository_Create_Success,
	"TestCreate_Error":            TestUserRepository_Create_Error,
	"TestFindByID_Success":        TestUserRepository_FindByID_Success,
	"TestFindByID_NotFound":       TestUserRepository_FindByID_NotFound,
	"TestFindByUsername_Success":  TestUserRepository_FindByUsername_Success,
	"TestFindByUsername_NotFound": TestUserRepository_FindByUsername_NotFound,
	"TestFindByEmail_Success":     TestUserRepository_FindByEmail_Success,
	"TestFindByEmail_NotFound":    TestUserRepository_FindByEmail_NotFound,
	"TestUpdate_Success":          TestUserRepository_Update_Success,
	"TestUpdate_Error":            TestUserRepository_Update_Error,
	"TestDelete_Success":          TestUserRepository_Delete_Success,
	"TestDelete_Error":            TestUserRepository_Delete_Error,
	"TestList_Success":            TestUserRepository_List_Success,
	"TestList_Error":              TestUserRepository_List_Error,
}

// setupMockDB creates a mock database connection
func setupUserMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
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

// TestUserRepository_Create_Success tests successful user creation
func TestUserRepository_Create_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide valid user data for creation
	user := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(2),
	}

	// Set expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			user.Username,
			user.Email,
			sqlmock.AnyArg(), // Password will be hashed
			user.FirstName,
			user.LastName,
			user.RoleID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the repository
	err = repo.Create(user)

	// Assert expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_Create_Error tests user creation when database returns an error
func TestUserRepository_Create_Error(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide valid user data for creation
	user := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(2),
	}

	// Mock database error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			user.Username,
			user.Email,
			sqlmock.AnyArg(), // Password will be hashed
			user.FirstName,
			user.LastName,
			user.RoleID,
		).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	// Call the repository
	err = repo.Create(user)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_FindByID_Success tests successful user retrieval by ID
func TestUserRepository_FindByID_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide a valid user ID for this test case
	userID := uint(1)

	// Mock user data that would be returned from DB
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "first_name", "last_name", "role_id"}).
		AddRow(userID, "2023-01-01", "2023-01-01", nil, "testuser", "test@example.com", "hashedpassword", "Test", "User", 2)

	// Role data for preload
	roleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "description"}).
		AddRow(2, "2023-01-01", "2023-01-01", nil, "user", "Regular user role")

	// Permission data for preload
	permissionRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "resource", "action", "description"}).
		AddRow(1, "2023-01-01", "2023-01-01", nil, "blogs", "read", "Can read blogs")

	// Role permissions data for preload
	rolePermissionsRows := sqlmock.NewRows([]string{"role_id", "permission_id"}).
		AddRow(2, 1)

	// Set expectations
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
		WithArgs(userID).
		WillReturnRows(rows)

	// Expect preload query for Role
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `roles` WHERE `roles`.`id` = ? AND `roles`.`deleted_at` IS NULL")).
		WithArgs(2).
		WillReturnRows(roleRows)

	// Expect preload query for Role.Permissions
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `role_permissions` WHERE `role_permissions`.`role_id` = ?")).
		WithArgs(2).
		WillReturnRows(rolePermissionsRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `permissions` WHERE `permissions`.`id` = ? AND `permissions`.`deleted_at` IS NULL")).
		WithArgs(1).
		WillReturnRows(permissionRows)

	// Call the repository
	user, err := repo.FindByID(userID)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, uint(2), user.RoleID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_FindByID_NotFound tests user retrieval when user is not found
func TestUserRepository_FindByID_NotFound(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide a non-existent user ID for this test case
	userID := uint(999)

	// Set expectations for not found
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	// Call the repository
	user, err := repo.FindByID(userID)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NotNil(t, user) // GORM returns an empty struct, not nil
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_FindByUsername_Success tests successful user retrieval by username
func TestUserRepository_FindByUsername_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide a valid username for this test case
	username := "testuser"

	// Mock user data that would be returned from DB
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "first_name", "last_name", "role_id"}).
		AddRow(1, "2023-01-01", "2023-01-01", nil, username, "test@example.com", "hashedpassword", "Test", "User", 2)

	// Role data for preload
	roleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "description"}).
		AddRow(2, "2023-01-01", "2023-01-01", nil, "user", "Regular user role")

	// Permission data for preload
	permissionRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "resource", "action", "description"}).
		AddRow(1, "2023-01-01", "2023-01-01", nil, "blogs", "read", "Can read blogs")

	// Role permissions data for preload
	rolePermissionsRows := sqlmock.NewRows([]string{"role_id", "permission_id"}).
		AddRow(2, 1)

	// Set expectations
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
		WithArgs(username).
		WillReturnRows(rows)

	// Expect preload query for Role
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `roles` WHERE `roles`.`id` = ? AND `roles`.`deleted_at` IS NULL")).
		WithArgs(2).
		WillReturnRows(roleRows)

	// Expect preload query for Role.Permissions
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `role_permissions` WHERE `role_permissions`.`role_id` = ?")).
		WithArgs(2).
		WillReturnRows(rolePermissionsRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `permissions` WHERE `permissions`.`id` = ? AND `permissions`.`deleted_at` IS NULL")).
		WithArgs(1).
		WillReturnRows(permissionRows)

	// Call the repository
	user, err := repo.FindByUsername(username)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, uint(2), user.RoleID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_FindByUsername_NotFound tests user retrieval when username is not found
func TestUserRepository_FindByUsername_NotFound(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide a non-existent username for this test case
	username := "nonexistentuser"

	// Set expectations for not found
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (username = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
		WithArgs(username).
		WillReturnError(gorm.ErrRecordNotFound)

	// Call the repository
	user, err := repo.FindByUsername(username)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NotNil(t, user) // GORM returns an empty struct, not nil
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_FindByEmail_Success tests successful user retrieval by email
func TestUserRepository_FindByEmail_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide a valid email for this test case
	email := "test@example.com"

	// Mock user data that would be returned from DB
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "first_name", "last_name", "role_id"}).
		AddRow(1, "2023-01-01", "2023-01-01", nil, "testuser", email, "hashedpassword", "Test", "User", 2)

	// Role data for preload
	roleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "description"}).
		AddRow(2, "2023-01-01", "2023-01-01", nil, "user", "Regular user role")

	// Permission data for preload
	permissionRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "resource", "action", "description"}).
		AddRow(1, "2023-01-01", "2023-01-01", nil, "blogs", "read", "Can read blogs")

	// Role permissions data for preload
	rolePermissionsRows := sqlmock.NewRows([]string{"role_id", "permission_id"}).
		AddRow(2, 1)

	// Set expectations
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (email = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
		WithArgs(email).
		WillReturnRows(rows)

	// Expect preload query for Role
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `roles` WHERE `roles`.`id` = ? AND `roles`.`deleted_at` IS NULL")).
		WithArgs(2).
		WillReturnRows(roleRows)

	// Expect preload query for Role.Permissions
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `role_permissions` WHERE `role_permissions`.`role_id` = ?")).
		WithArgs(2).
		WillReturnRows(rolePermissionsRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `permissions` WHERE `permissions`.`id` = ? AND `permissions`.`deleted_at` IS NULL")).
		WithArgs(1).
		WillReturnRows(permissionRows)

	// Call the repository
	user, err := repo.FindByEmail(email)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, uint(2), user.RoleID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_FindByEmail_NotFound tests user retrieval when email is not found
func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide a non-existent email for this test case
	email := "nonexistent@example.com"

	// Set expectations for not found
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (email = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` ASC LIMIT 1")).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound)

	// Call the repository
	user, err := repo.FindByEmail(email)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NotNil(t, user) // GORM returns an empty struct, not nil
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_Update_Success tests successful user update
func TestUserRepository_Update_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide valid user data for update
	user := &models.User{
		Username:  "updateduser",
		Email:     "updated@example.com",
		Password:  "updatedpassword",
		FirstName: "Updated",
		LastName:  "User",
		RoleID:    uint(2),
	}
	user.ID = uint(1)

	// Set expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET")).
		WithArgs(
			sqlmock.AnyArg(), // UpdatedAt
			user.Username,
			user.Email,
			sqlmock.AnyArg(), // Password will be hashed
			user.FirstName,
			user.LastName,
			user.RoleID,
			user.ID,
			sqlmock.AnyArg(), // DeletedAt
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the repository
	err = repo.Update(user)

	// Assert expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_Update_Error tests user update when database returns an error
func TestUserRepository_Update_Error(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide valid user data for update
	user := &models.User{
		Username:  "updateduser",
		Email:     "updated@example.com",
		Password:  "updatedpassword",
		FirstName: "Updated",
		LastName:  "User",
		RoleID:    uint(2),
	}
	user.ID = uint(1)

	// Mock database error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET")).
		WithArgs(
			sqlmock.AnyArg(), // UpdatedAt
			user.Username,
			user.Email,
			sqlmock.AnyArg(), // Password will be hashed
			user.FirstName,
			user.LastName,
			user.RoleID,
			user.ID,
			sqlmock.AnyArg(), // DeletedAt
		).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	// Call the repository
	err = repo.Update(user)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_Delete_Success tests successful user deletion
func TestUserRepository_Delete_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide a valid user ID for deletion
	userID := uint(1)

	// Set expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			userID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the repository
	err = repo.Delete(userID)

	// Assert expectations
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_Delete_Error tests user deletion when database returns an error
func TestUserRepository_Delete_Error(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide a valid user ID for deletion
	userID := uint(1)

	// Mock database error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			userID,
		).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	// Call the repository
	err = repo.Delete(userID)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_List_Success tests successful user listing
func TestUserRepository_List_Success(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide valid data for user listing
	offset := 0
	limit := 10

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// Mock users data
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "email", "password", "first_name", "last_name", "role_id"}).
		AddRow(1, "2023-01-01", "2023-01-01", nil, "user1", "user1@example.com", "hashedpassword", "User", "One", 2).
		AddRow(2, "2023-01-02", "2023-01-02", nil, "user2", "user2@example.com", "hashedpassword", "User", "Two", 2)

	// Role data for preload
	roleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "description"}).
		AddRow(2, "2023-01-01", "2023-01-01", nil, "user", "Regular user role")

	// Set expectations for count
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`deleted_at` IS NULL")).
		WillReturnRows(countRows)

	// Set expectations for getting users
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT 10 OFFSET 0")).
		WillReturnRows(rows)

	// Expect preload query for Role
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `roles` WHERE `roles`.`id` = ? AND `roles`.`deleted_at` IS NULL")).
		WithArgs(2).
		WillReturnRows(roleRows)

	// Call the repository
	users, count, err := repo.List(offset, limit)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1", users[0].Username)
	assert.Equal(t, "user2", users[1].Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestUserRepository_List_Error tests user listing when database returns an error
func TestUserRepository_List_Error(t *testing.T) {
	// Setup mock DB
	mockDB, mock, err := setupUserMockDB(t)
	require.NoError(t, err)
	defer mockDB.Close()

	// Create repository with mock DB
	repo := impl.NewUserRepository(mockDB)

	// Test data
	// TODO: Provide valid data for user listing
	offset := 0
	limit := 10

	// Mock database error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`deleted_at` IS NULL")).
		WillReturnError(errors.New("database error"))

	// Call the repository
	users, count, err := repo.List(offset, limit)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.Empty(t, users)
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
