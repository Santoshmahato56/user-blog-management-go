package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/userblog/management/api/controller"
	"github.com/userblog/management/api/dto"
	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/service"
	"net/http"
	"strconv"
)

// UserController implements the IUserController interface
type UserController struct {
	userService service.IUserService
}

// NewUserController creates a new user controller
func NewUserController(userService service.IUserService) controller.IUserController {
	return &UserController{
		userService: userService,
	}
}

// Create handles the create user API endpoint
func (c *UserController) Create(ctx *gin.Context) {
	var req dto.CreateUserRequest
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

	// Create the user
	if err := c.userService.Create(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove sensitive fields
	user.Password = ""

	ctx.JSON(http.StatusCreated, user)
}

// GetByID handles the get user by ID API endpoint
func (c *UserController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user
	user, err := c.userService.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Remove sensitive fields
	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

// Update handles the update user API endpoint
func (c *UserController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context to check permissions
	userInterface, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	currentUser, ok := userInterface.(models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast user from context"})
		return
	}

	// Only admin can update other users
	if currentUser.Role.Name != "admin" && currentUser.ID != uint(id) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own user information"})
		return
	}

	// Create user model from request
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	user.ID = uint(id)

	// Only admin can change roles
	if currentUser.Role.Name == "admin" && req.RoleID != 0 {
		user.RoleID = req.RoleID
	}

	// Update the user
	if err := c.userService.Update(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Remove sensitive fields
	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

// Delete handles the delete user API endpoint
func (c *UserController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Delete the user
	if err := c.userService.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// List handles the list users API endpoint
func (c *UserController) List(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	perPageStr := ctx.DefaultQuery("per_page", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// List users
	users, count, err := c.userService.List(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Remove sensitive fields
	for i := range users {
		users[i].Password = ""
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       users,
		"total":      count,
		"page":       page,
		"per_page":   perPage,
		"total_page": (count + perPage - 1) / perPage,
	})
}
