package impl_test

import (
	"errors"
	"fmt"
	"github.com/userblog/management/internal/repository/mocks"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/service/impl"
)

func TestMain(m *testing.M) {
	// Setup environment variables for testing
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("TOKEN_EXPIRY", "24")

	// Run tests
	code := m.Run()

	// Exit
	os.Exit(code)
}

// TestUser extends models.User to allow for test mocking
type TestUser struct {
	models.User
	MockValidatePassword func(password string) error
}

// ValidatePassword overrides the original method for testing
func (u *TestUser) ValidatePassword(password string) error {
	if u.MockValidatePassword != nil {
		return u.MockValidatePassword(password)
	}
	return nil
}

// TestAuthServiceSuite executes all test cases for AuthService
func TestAuthServiceSuite(t *testing.T) {
	for name, test := range AuthServiceTestBed {
		t.Run(name, test)
	}
}

var AuthServiceTestBed = map[string]func(*testing.T){
	"TestRegister_Success":               TestAuthService_Register_Success,
	"TestRegister_UsernameExists":        TestAuthService_Register_UsernameExists,
	"TestRegister_EmailExists":           TestAuthService_Register_EmailExists,
	"TestRegister_RepositoryError":       TestAuthService_Register_RepositoryError,
	"TestLogin_Success":                  TestAuthService_Login_Success,
	"TestLogin_UserNotFound":             TestAuthService_Login_UserNotFound,
	"TestLogin_InvalidPassword":          TestAuthService_Login_InvalidPassword,
	"TestGetUserByID_Success":            TestAuthService_GetUserByID_Success,
	"TestGetUserByID_NotFound":           TestAuthService_GetUserByID_NotFound,
	"TestValidatePermission_Granted":     TestAuthService_ValidatePermission_Granted,
	"TestValidatePermission_Denied":      TestAuthService_ValidatePermission_Denied,
	"TestExtractTokenFromHeader_Valid":   TestAuthService_ExtractTokenFromHeader_Valid,
	"TestExtractTokenFromHeader_Invalid": TestAuthService_ExtractTokenFromHeader_Invalid,
	"TestExtractTokenFromHeader_Missing": TestAuthService_ExtractTokenFromHeader_Missing,
}

// TestAuthService_Register_Success tests successful user registration
func TestAuthService_Register_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide valid data for user registration
	user := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(0), // Should default to 2
	}

	// Set expectations for FindByUsername - should return no user
	mockRepo.On("FindByUsername", user.Username).Return(nil, errors.New("not found"))

	// Set expectations for FindByEmail - should return no user
	mockRepo.On("FindByEmail", user.Email).Return(nil, errors.New("not found"))

	// Set expectations for Create
	mockRepo.On("Create", mock.MatchedBy(func(u *models.User) bool {
		return u.Username == user.Username &&
			u.Email == user.Email &&
			u.Password == user.Password &&
			u.RoleID == uint(2) // Default role
	})).Return(nil)

	// Call the service
	err := authService.Register(user)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, uint(2), user.RoleID) // Default role should be set
	mockRepo.AssertExpectations(t)
}

// TestAuthService_Register_UsernameExists tests user registration when username already exists
func TestAuthService_Register_UsernameExists(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide valid data for user registration
	user := &models.User{
		Username:  "existinguser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(0),
	}

	existingUser := &models.User{
		Username: "existinguser",
		Email:    "existing@example.com",
	}
	existingUser.ID = uint(1)

	// Set expectations for FindByUsername - should return existing user
	mockRepo.On("FindByUsername", user.Username).Return(existingUser, nil)

	// Call the service
	err := authService.Register(user)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username already exists")
	mockRepo.AssertExpectations(t)
}

// TestAuthService_Register_EmailExists tests user registration when email already exists
func TestAuthService_Register_EmailExists(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide valid data for user registration
	user := &models.User{
		Username:  "testuser",
		Email:     "existing@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(0),
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
	err := authService.Register(user)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
	mockRepo.AssertExpectations(t)
}

// TestAuthService_Register_RepositoryError tests user registration when repository returns an error
func TestAuthService_Register_RepositoryError(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide valid data for user registration
	user := &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		RoleID:    uint(0),
	}

	expectedError := errors.New("database error")

	// Set expectations for FindByUsername - should return no user
	mockRepo.On("FindByUsername", user.Username).Return(nil, errors.New("not found"))

	// Set expectations for FindByEmail - should return no user
	mockRepo.On("FindByEmail", user.Email).Return(nil, errors.New("not found"))

	// Set expectations for Create - should return error
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(expectedError)

	// Call the service
	err := authService.Register(user)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

// TestAuthService_Login_Success tests successful user login
func TestAuthService_Login_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide valid username and password for login
	username := "testuser"
	password := "password123"

	// Create a test user that can be properly mocked
	testUser := &TestUser{}
	testUser.User.Username = username
	testUser.User.Password = "$2a$10$pSSL5DKKN8hVpFdU15.5..a8LRHxnEZ9hMpMjPGGi.qSn1yJj.lPe" // hashed "password123"
	testUser.User.RoleID = uint(2)
	testUser.User.ID = uint(1)

	// Mock passwordValidation method
	testUser.MockValidatePassword = func(pwd string) error {
		return nil // Simulate successful validation
	}

	// Set expectations for FindByUsername - should return user
	mockRepo.On("FindByUsername", username).Return(&testUser.User, nil)

	// Call the service
	token, err := authService.Login(username, password)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

// TestAuthService_Login_UserNotFound tests login when user is not found
func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide non-existent username and password for login
	username := "nonexistentuser"
	password := "password123"

	// Set expectations for FindByUsername - should return not found error
	mockRepo.On("FindByUsername", username).Return(nil, errors.New("record not found"))

	// Call the service
	token, err := authService.Login(username, password)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid username or password")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

// TestAuthService_Login_InvalidPassword tests login with invalid password
func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide username and invalid password for login
	username := "testuser"
	password := "wrongpassword"

	// Create a test user that can be properly mocked
	testUser := &TestUser{}
	testUser.User.Username = username
	testUser.User.Password = "$2a$10$pSSL5DKKN8hVpFdU15.5..a8LRHxnEZ9hMpMjPGGi.qSn1yJj.lPe" // hashed "password123"
	testUser.User.RoleID = uint(2)
	testUser.User.ID = uint(1)

	// Mock passwordValidation method
	testUser.MockValidatePassword = func(pwd string) error {
		return errors.New("invalid password")
	}

	// Set expectations for FindByUsername - should return user
	mockRepo.On("FindByUsername", username).Return(&testUser.User, nil)

	// Call the service
	token, err := authService.Login(username, password)

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid username or password")
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

// TestAuthService_GetUserByID_Success tests successful user retrieval by ID
func TestAuthService_GetUserByID_Success(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

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
	user, err := authService.GetUserByID(userID)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

// TestAuthService_GetUserByID_NotFound tests user retrieval when user is not found
func TestAuthService_GetUserByID_NotFound(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide a non-existent user ID for this test case
	userID := uint(999)
	expectedError := errors.New("record not found")

	// Set expectations
	mockRepo.On("FindByID", userID).Return(nil, expectedError)

	// Call the service
	user, err := authService.GetUserByID(userID)

	// Assert expectations
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// TestAuthService_ValidatePermission_Granted tests permission validation when permission is granted
func TestAuthService_ValidatePermission_Granted(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide valid user with permissions for this test case
	permission := models.Permission{
		Resource:    "blogs",
		Action:      "create",
		Description: "Can create blogs",
	}

	role := models.Role{
		Name:        "user",
		Description: "Regular user role",
		Permissions: []models.Permission{permission},
	}

	user := &models.User{
		Username: "testuser",
		Role:     role,
	}

	// Call the service
	hasPermission := authService.ValidatePermission(user, "blogs", "create")

	// Assert expectations
	assert.True(t, hasPermission)
}

// TestAuthService_ValidatePermission_Denied tests permission validation when permission is denied
func TestAuthService_ValidatePermission_Denied(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide valid user without specific permission for this test case
	permission := models.Permission{
		Resource:    "blogs",
		Action:      "read",
		Description: "Can read blogs",
	}

	role := models.Role{
		Name:        "user",
		Description: "Regular user role",
		Permissions: []models.Permission{permission},
	}

	user := &models.User{
		Username: "testuser",
		Role:     role,
	}

	// Call the service - user has "read" but not "create" permission
	hasPermission := authService.ValidatePermission(user, "blogs", "create")

	// Assert expectations
	assert.False(t, hasPermission)

	// Also test with nil user or nil permissions
	assert.False(t, authService.ValidatePermission(nil, "blogs", "create"))

	userWithoutPermissions := &models.User{
		Username: "testuser",
		Role: models.Role{
			Name:        "limited",
			Description: "Limited role with no permissions",
		},
	}
	assert.False(t, authService.ValidatePermission(userWithoutPermissions, "blogs", "create"))
}

// TestAuthService_ExtractTokenFromHeader_Valid tests token extraction from a valid header
func TestAuthService_ExtractTokenFromHeader_Valid(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide a valid Authorization header for this test case
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.ZQgFDnAFKA4hZCNMfopNGtrl_TUAcnFf2vvvKGq2z6E"
	authHeader := fmt.Sprintf("Bearer %s", token)

	// Call the service
	extractedToken, err := authService.ExtractTokenFromHeader(authHeader)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, token, extractedToken)
}

// TestAuthService_ExtractTokenFromHeader_Invalid tests token extraction from an invalid header
func TestAuthService_ExtractTokenFromHeader_Invalid(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Test data
	// TODO: Provide an invalid Authorization header for this test case
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.ZQgFDnAFKA4hZCNMfopNGtrl_TUAcnFf2vvvKGq2z6E"

	// Invalid headers
	invalidHeaders := []string{
		token,              // No "Bearer" prefix
		"Token " + token,   // Wrong prefix
		"Bearer",           // No token
		"Bearer  " + token, // Extra space
	}

	// Test all invalid headers
	for _, header := range invalidHeaders {
		extractedToken, err := authService.ExtractTokenFromHeader(header)

		// Assert expectations
		assert.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "authorization header format")
		assert.Empty(t, extractedToken)
	}
}

// TestAuthService_ExtractTokenFromHeader_Missing tests token extraction when header is missing
func TestAuthService_ExtractTokenFromHeader_Missing(t *testing.T) {
	// Initialize mock repository
	mockRepo := new(mocks.MockUserRepository)

	// Create the service with mock repository
	authService := impl.NewAuthService(mockRepo)

	// Call the service with empty header
	extractedToken, err := authService.ExtractTokenFromHeader("")

	// Assert expectations
	assert.Error(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "authorization header is required")
	assert.Empty(t, extractedToken)
}
