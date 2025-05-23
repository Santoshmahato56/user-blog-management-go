package impl_test

import (
	"errors"
	"github.com/userblog/management/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/service/impl"
)

// TestUserServiceSuite executes all test cases for UserService
func TestUserServiceSuite(t *testing.T) {
	for name, test := range UserServiceTestBed {
		t.Run(name, test)
	}
}

var UserServiceTestBed = map[string]func(*testing.T){
	"TestCreate_Success":         TestUserService_Create_Success,
	"TestCreate_UsernameExists":  TestUserService_Create_UsernameExists,
	"TestCreate_EmailExists":     TestUserService_Create_EmailExists,
	"TestCreate_RepositoryError": TestUserService_Create_RepositoryError,
	"TestGetByID_Success":        TestUserService_GetByID_Success,
	"TestGetByID_NotFound":       TestUserService_GetByID_NotFound,
	"TestUpdate_Success":         TestUserService_Update_Success,
	"TestUpdate_NotFound":        TestUserService_Update_NotFound,
	"TestUpdate_UsernameExists":  TestUserService_Update_UsernameExists,
	"TestUpdate_EmailExists":     TestUserService_Update_EmailExists,
	"TestUpdate_RepositoryError": TestUserService_Update_RepositoryError,
	"TestDelete_Success":         TestUserService_Delete_Success,
	"TestDelete_Error":           TestUserService_Delete_Error,
	"TestList_Success":           TestUserService_List_Success,
	"TestList_RepositoryError":   TestUserService_List_RepositoryError,
}

// TestUserService_Create_Success tests successful user creation
func TestUserService_Create_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide valid data for user creation
	user := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(2),
	}

	// Set expectations for FindByUsername - should return no user
	mockRepo.On("FindByUsername", user.Username).Return(nil, errors.New("not found"))

	// Set expectations for FindByEmail - should return no user
	mockRepo.On("FindByEmail", user.Email).Return(nil, errors.New("not found"))

	// Set expectations for Create
	mockRepo.On("Create", user).Return(nil)

	// Call the service
	err := userService.Create(user)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Create_UsernameExists tests user creation when username already exists
func TestUserService_Create_UsernameExists(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide valid data for user creation
	user := &models.User{
		Username:  "existinguser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(2),
	}

	existingUser := &models.User{
		Username: "existinguser",
		Email:    "existing@example.com",
	}
	existingUser.ID = uint(1)

	// Set expectations for FindByUsername - should return existing user
	mockRepo.On("FindByUsername", user.Username).Return(existingUser, nil)

	// Call the service
	err := userService.Create(user)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username already exists")
	mockRepo.AssertExpectations(t)
}

// TestUserService_Create_EmailExists tests user creation when email already exists
func TestUserService_Create_EmailExists(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide valid data for user creation
	user := &models.User{
		Username:  "testuser",
		Email:     "existing@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(2),
	}

	existingUser := &models.User{
		Username: "existinguser",
		Email:    "existing@example.com",
	}
	existingUser.ID = uint(1)

	// Set expectations for FindByUsername - should return no user
	mockRepo.On("FindByUsername", user.Username).Return(nil, errors.New("not found"))

	// Set expectations for FindByEmail - should return existing user
	mockRepo.On("FindByEmail", user.Email).Return(existingUser, nil)

	// Call the service
	err := userService.Create(user)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
	mockRepo.AssertExpectations(t)
}

// TestUserService_Create_RepositoryError tests user creation when repository returns an error
func TestUserService_Create_RepositoryError(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide valid data for user creation
	user := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(2),
	}

	expectedError := errors.New("database error")

	// Set expectations for FindByUsername - should return no user
	mockRepo.On("FindByUsername", user.Username).Return(nil, errors.New("not found"))

	// Set expectations for FindByEmail - should return no user
	mockRepo.On("FindByEmail", user.Email).Return(nil, errors.New("not found"))

	// Set expectations for Create - should return error
	mockRepo.On("Create", user).Return(expectedError)

	// Call the service
	err := userService.Create(user)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetByID_Success tests successful user retrieval by ID
func TestUserService_GetByID_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide a valid user ID for this test case
	userID := uint(1)
	expectedUser := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(2),
	}
	expectedUser.ID = userID

	// Set expectations
	mockRepo.On("FindByID", userID).Return(expectedUser, nil)

	// Call the service
	user, err := userService.GetByID(userID)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

// TestUserService_GetByID_NotFound tests user retrieval when user is not found
func TestUserService_GetByID_NotFound(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide a non-existent user ID for this test case
	userID := uint(999)
	expectedError := errors.New("record not found")

	// Set expectations
	mockRepo.On("FindByID", userID).Return(nil, expectedError)

	// Call the service
	user, err := userService.GetByID(userID)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Update_Success tests successful user update
func TestUserService_Update_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide valid data for user update
	userID := uint(1)
	updateUser := &models.User{
		Username:  "updateduser",
		Email:     "updated@example.com",
		Password:  "updatedpassword",
		FirstName: "Updated",
		LastName:  "User",
		RoleID:    uint(2),
	}
	updateUser.ID = userID

	existingUser := &models.User{
		Username:  "originaluser",
		Email:     "original@example.com",
		Password:  "hashedpassword",
		FirstName: "Original",
		LastName:  "User",
		RoleID:    uint(2),
	}
	existingUser.ID = userID

	// Set expectations for FindByID - should return existing user
	mockRepo.On("FindByID", userID).Return(existingUser, nil)

	// Set expectations for FindByUsername - should return no user for new username
	mockRepo.On("FindByUsername", updateUser.Username).Return(nil, errors.New("not found"))

	// Set expectations for FindByEmail - should return no user for new email
	mockRepo.On("FindByEmail", updateUser.Email).Return(nil, errors.New("not found"))

	// Set expectations for Update with updated user fields
	mockRepo.On("Update", mock.MatchedBy(func(u *models.User) bool {
		return u.ID == existingUser.ID &&
			u.Username == updateUser.Username &&
			u.Email == updateUser.Email &&
			u.FirstName == updateUser.FirstName &&
			u.LastName == updateUser.LastName &&
			u.Password == updateUser.Password &&
			u.RoleID == updateUser.RoleID
	})).Return(nil)

	// Call the service
	err := userService.Update(updateUser)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Update_NotFound tests user update when user is not found
func TestUserService_Update_NotFound(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide a non-existent user ID for this test case
	userID := uint(999)
	updateUser := &models.User{
		Username:  "updateduser",
		Email:     "updated@example.com",
		Password:  "updatedpassword",
		FirstName: "Updated",
		LastName:  "User",
		RoleID:    uint(2),
	}
	updateUser.ID = userID

	expectedError := errors.New("record not found")

	// Set expectations for FindByID - should return error
	mockRepo.On("FindByID", userID).Return(nil, expectedError)

	// Call the service
	err := userService.Update(updateUser)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Update_UsernameExists tests user update when username already exists
func TestUserService_Update_UsernameExists(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide valid data for user update with username conflict
	userID := uint(1)
	updateUser := &models.User{
		Username:  "existinguser",
		Email:     "updated@example.com",
		Password:  "updatedpassword",
		FirstName: "Updated",
		LastName:  "User",
		RoleID:    uint(2),
	}
	updateUser.ID = userID

	existingUser := &models.User{
		Username:  "originaluser",
		Email:     "original@example.com",
		Password:  "hashedpassword",
		FirstName: "Original",
		LastName:  "User",
		RoleID:    uint(2),
	}
	existingUser.ID = userID

	conflictingUser := &models.User{
		Username: "existinguser",
		Email:    "another@example.com",
	}
	conflictingUser.ID = uint(2)

	// Set expectations for FindByID - should return existing user
	mockRepo.On("FindByID", userID).Return(existingUser, nil)

	// Set expectations for FindByUsername - should return conflicting user
	mockRepo.On("FindByUsername", updateUser.Username).Return(conflictingUser, nil)

	// Call the service
	err := userService.Update(updateUser)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username already exists")
	mockRepo.AssertExpectations(t)
}

// TestUserService_Update_EmailExists tests user update when email already exists
func TestUserService_Update_EmailExists(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide valid data for user update with email conflict
	userID := uint(1)
	updateUser := &models.User{
		Username:  "updateduser",
		Email:     "existing@example.com",
		Password:  "updatedpassword",
		FirstName: "Updated",
		LastName:  "User",
		RoleID:    uint(2),
	}
	updateUser.ID = userID

	existingUser := &models.User{
		Username:  "originaluser",
		Email:     "original@example.com",
		Password:  "hashedpassword",
		FirstName: "Original",
		LastName:  "User",
		RoleID:    uint(2),
	}
	existingUser.ID = userID

	conflictingUser := &models.User{
		Username: "anotheruser",
		Email:    "existing@example.com",
	}
	conflictingUser.ID = uint(2)

	// Set expectations for FindByID - should return existing user
	mockRepo.On("FindByID", userID).Return(existingUser, nil)

	// Set expectations for FindByUsername - should return no user for new username
	mockRepo.On("FindByUsername", updateUser.Username).Return(nil, errors.New("not found"))

	// Set expectations for FindByEmail - should return conflicting user
	mockRepo.On("FindByEmail", updateUser.Email).Return(conflictingUser, nil)

	// Call the service
	err := userService.Update(updateUser)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
	mockRepo.AssertExpectations(t)
}

// TestUserService_Update_RepositoryError tests user update when repository returns an error
func TestUserService_Update_RepositoryError(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide valid data for user update
	userID := uint(1)
	updateUser := &models.User{
		Username:  "updateduser",
		Email:     "updated@example.com",
		Password:  "updatedpassword",
		FirstName: "Updated",
		LastName:  "User",
		RoleID:    uint(2),
	}
	updateUser.ID = userID

	existingUser := &models.User{
		Username:  "originaluser",
		Email:     "original@example.com",
		Password:  "hashedpassword",
		FirstName: "Original",
		LastName:  "User",
		RoleID:    uint(2),
	}
	existingUser.ID = userID

	expectedError := errors.New("database error")

	// Set expectations for FindByID - should return existing user
	mockRepo.On("FindByID", userID).Return(existingUser, nil)

	// Set expectations for FindByUsername - should return no user for new username
	mockRepo.On("FindByUsername", updateUser.Username).Return(nil, errors.New("not found"))

	// Set expectations for FindByEmail - should return no user for new email
	mockRepo.On("FindByEmail", updateUser.Email).Return(nil, errors.New("not found"))

	// Set expectations for Update - should return error
	mockRepo.On("Update", mock.MatchedBy(func(u *models.User) bool {
		return u.ID == existingUser.ID &&
			u.Username == updateUser.Username &&
			u.Email == updateUser.Email &&
			u.FirstName == updateUser.FirstName &&
			u.LastName == updateUser.LastName &&
			u.Password == updateUser.Password &&
			u.RoleID == updateUser.RoleID
	})).Return(expectedError)

	// Call the service
	err := userService.Update(updateUser)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Delete_Success tests successful user deletion
func TestUserService_Delete_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide a valid user ID for deletion
	userID := uint(1)

	// Set expectations
	mockRepo.On("Delete", userID).Return(nil)

	// Call the service
	err := userService.Delete(userID)

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_Delete_Error tests user deletion when repository returns an error
func TestUserService_Delete_Error(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide a valid user ID for deletion
	userID := uint(1)
	expectedError := errors.New("database error")

	// Set expectations
	mockRepo.On("Delete", userID).Return(expectedError)

	// Call the service
	err := userService.Delete(userID)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

// TestUserService_List_Success tests successful user listing
func TestUserService_List_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide valid data for user listing
	page := 1
	perPage := 10
	offset := (page - 1) * perPage

	expectedUsers := []models.User{
		{
			Username:  "user1",
			Email:     "user1@example.com",
			FirstName: "User",
			LastName:  "One",
			RoleID:    uint(2),
		},
		{
			Username:  "user2",
			Email:     "user2@example.com",
			FirstName: "User",
			LastName:  "Two",
			RoleID:    uint(2),
		},
	}
	expectedCount := 2

	// Set expectations
	mockRepo.On("List", offset, perPage).Return(expectedUsers, expectedCount, nil)

	// Call the service
	users, count, err := userService.List(page, perPage)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}

// TestUserService_List_RepositoryError tests user listing when repository returns an error
func TestUserService_List_RepositoryError(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	userService := impl.NewUserService(mockRepo)

	// Test data
	// TODO: Provide valid data for user listing
	page := 1
	perPage := 10
	offset := (page - 1) * perPage

	expectedError := errors.New("database error")

	// Set expectations
	mockRepo.On("List", offset, perPage).Return([]models.User{}, 0, expectedError)

	// Call the service
	users, count, err := userService.List(page, perPage)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, users)
	assert.Equal(t, 0, count)
	mockRepo.AssertExpectations(t)
}
