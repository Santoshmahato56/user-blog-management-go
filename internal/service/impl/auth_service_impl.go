package impl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/repository"
	"github.com/userblog/management/internal/service"
	"github.com/userblog/management/pkg/config"
)

// AuthService implements the IAuthService interface
type AuthService struct {
	userRepo repository.IUserRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.IUserRepository) service.IAuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register registers a new user
func (s *AuthService) Register(user *models.User) error {
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

	// Set default role if not provided
	if user.RoleID == 0 {
		user.RoleID = 2 // Default to 'user' role (ID=2)
	}

	// Create the user
	err = s.userRepo.Create(user)
	if err != nil {
		return err
	}

	return nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(username, password string) (string, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	// Validate password
	if err := user.ValidatePassword(password); err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID returns a user by ID
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// ValidateToken validates a JWT token and returns the user
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	// Load environment variables
	_ = godotenv.Load()
	jwtSecret := config.GetOrDefaultString("JWT_SECRET", "")

	// Check if the token is empty
	if tokenString == "" {
		return nil, errors.New("token is required")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check if the token is expired
	expires := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expires) {
		return nil, errors.New("token has expired")
	}

	// Get user ID from claims
	userID := uint(claims["user_id"].(float64))

	// Load user from database
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// ValidatePermission checks if a user has the required permission
func (s *AuthService) ValidatePermission(user *models.User, resource, action string) bool {
	if user == nil || user.Role.Permissions == nil {
		return false
	}

	// Check if user has the required permission
	for _, permission := range user.Role.Permissions {
		if permission.Resource == resource && permission.Action == action {
			return true
		}
	}

	return false
}

// ExtractTokenFromHeader extracts the JWT token from the Authorization header
func (s *AuthService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if the header has the Bearer prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

// generateToken generates a JWT token for a user
func generateToken(user *models.User) (string, error) {
	// Load environment variables
	_ = godotenv.Load()
	jwtSecret := config.GetOrDefaultString("JWT_SECRET", "your-secret-key")
	tokenExpiry := config.GetOrDefaultString("TOKEN_EXPIRY", "24")

	// Parse token expiry duration
	expiryHours := 24 // Default to 24 hours
	_, err := fmt.Sscanf(tokenExpiry, "%d", &expiryHours)
	if err != nil {
		expiryHours = 24
	}

	// Create token claims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role_id":  user.RoleID,
		"exp":      time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
