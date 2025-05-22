package impl

import (
	"github.com/userblog/management/api/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/userblog/management/api/controller/interfaces"
	"github.com/userblog/management/internal/models"
	"github.com/userblog/management/internal/service"
)

// BlogController implements the IBlogController interface
type BlogController struct {
	blogService service.IBlogService
}

// NewBlogController creates a new blog controller
func NewBlogController(blogService service.IBlogService) interfaces.IBlogController {
	return &BlogController{
		blogService: blogService,
	}
}

// Create handles the create blog API endpoint
func (c *BlogController) Create(ctx *gin.Context) {
	var req dto.CreateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	// Create blog model from request
	blog := models.Blog{
		Title:     req.Title,
		Content:   req.Content,
		Published: req.Published,
	}

	// Create the blog
	if err := c.blogService.Create(&blog, user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, blog)
}

// GetByID handles the get blog by ID API endpoint
func (c *BlogController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	// Get the blog
	blog, err := c.blogService.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// Check if the blog is published or if the user is the owner
	userInterface, exists := ctx.Get("user")
	if exists {
		user, ok := userInterface.(models.User)
		if ok && (user.ID == blog.UserID || user.Role.Name == "admin") {
			ctx.JSON(http.StatusOK, blog)
			return
		}
	}

	// If the blog is not published and the user is not the owner or admin
	if !blog.Published {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Blog is not published"})
		return
	}

	ctx.JSON(http.StatusOK, blog)
}

// Update handles the update blog API endpoint
func (c *BlogController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var req dto.UpdateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	// Create blog model from request
	blog := models.Blog{
		Title:     req.Title,
		Content:   req.Content,
		Published: req.Published,
	}
	blog.ID = uint(id)

	// Update the blog
	if err := c.blogService.Update(&blog, user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, blog)
}

// Delete handles the delete blog API endpoint
func (c *BlogController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

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

	// Delete the blog
	if err := c.blogService.Delete(uint(id), user.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}

// List handles the list blogs API endpoint
func (c *BlogController) List(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	perPageStr := ctx.DefaultQuery("per_page", "10")
	publishedOnly := ctx.DefaultQuery("published_only", "false")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	pubOnly, err := strconv.ParseBool(publishedOnly)
	if err != nil {
		pubOnly = false
	}

	// List blogs
	blogs, count, err := c.blogService.List(page, perPage, pubOnly)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       blogs,
		"total":      count,
		"page":       page,
		"per_page":   perPage,
		"total_page": (count + perPage - 1) / perPage,
	})
}

// ListByUser handles the list blogs by user API endpoint
func (c *BlogController) ListByUser(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	pageStr := ctx.DefaultQuery("page", "1")
	perPageStr := ctx.DefaultQuery("per_page", "10")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// List blogs by user
	blogs, count, err := c.blogService.ListByUser(uint(userID), page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       blogs,
		"total":      count,
		"page":       page,
		"per_page":   perPage,
		"total_page": (count + perPage - 1) / perPage,
	})
}
