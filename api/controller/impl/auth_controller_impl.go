package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/userblog/management/api/controller"
	"github.com/userblog/management/api/dto"
	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/service"
	"net/http"
)

// AuthController implements the IAuthController interface
type AuthController struct {
	authService service.IAuthService
}

// NewAuthController creates a new authentication controller
func NewAuthController(authService service.IAuthService) controller.IAuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register handles the register API endpoint
func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user model from request
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleID:    req.RoleID,
	}

	// Register the user
	if err := c.authService.Register(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	})
}

// Login handles the login API endpoint
func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate the user
	token, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
	})
}

// GetMe handles the get current user API endpoint
func (c *AuthController) GetMe(ctx *gin.Context) {
	// Get user from context
	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast user from context"})
		return
	}

	user.Password = ""
	ctx.JSON(http.StatusOK, user)
}
