package main

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/userblog/management/api/controller"
	controllerImpl "github.com/userblog/management/api/controller/impl"
	"github.com/userblog/management/api/middleware"
	middlewareImpl "github.com/userblog/management/api/middleware/impl"
	"github.com/userblog/management/api/route"
	"github.com/userblog/management/internal/repository"
	repoImpl "github.com/userblog/management/internal/repository/impl"
	"github.com/userblog/management/internal/service"
	serviceImpl "github.com/userblog/management/internal/service/impl"
	"github.com/userblog/management/pkg/config"
	"github.com/userblog/management/pkg/db"
	"github.com/userblog/management/pkg/logger"
)

func main() {

	// Create background context with debug ID
	ctx := context.Background()
	ctx = logger.AddToContext(ctx, logger.DebugIDKey, "startup")

	// Initialize database
	database := db.Connect()
	defer func(database *gorm.DB) {
		err := database.Close()
		if err != nil {
			logger.Error(ctx, "Failed to close database connection")
		}
	}(database)

	logger.Info(ctx, "Database connected")

	// Seed database with roles and permissions
	initializeDatabaseScript(database)
	logger.Info(ctx, "Database schema initialized and seeded with roles and permissions")

	// Initialize repositories with the database connection
	var userRepo repository.IUserRepository = repoImpl.NewUserRepository(database)
	var blogRepo repository.IBlogRepository = repoImpl.NewBlogRepository(database)

	// Initialize services
	var authService service.IAuthService = serviceImpl.NewAuthService(userRepo)
	var userService service.IUserService = serviceImpl.NewUserService(userRepo)
	var blogService service.IBlogService = serviceImpl.NewBlogService(blogRepo)

	// Initialize middleware
	var authMiddleware middleware.IAuthMiddleware = middlewareImpl.NewAuthMiddleware(authService)

	// Initialize controllers
	var authController controller.IAuthController = controllerImpl.NewAuthController(authService)
	var userController controller.IUserController = controllerImpl.NewUserController(userService)
	var blogController controller.IBlogController = controllerImpl.NewBlogController(blogService)

	// Initialize routes
	authRoute := route.NewAuthRoute(authController, authMiddleware)
	userRoute := route.NewUserRoute(userController, authMiddleware)
	blogRoute := route.NewBlogRoute(blogController, authMiddleware)

	// Initialize router
	router := gin.Default()

	// Apply middlewares
	router.Use(middleware.GlobalExceptionHandler())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Create API router group
	api := router.Group("/api")

	// Register routes
	authRoute.AuthRoute(api)
	userRoute.UserRoute(api)
	blogRoute.BlogRoute(api)

	// Start server
	port := config.GetOrDefaultString("PORT", "8080")

	logger.InfoF(ctx, "Server starting on port: %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Fatal(ctx, "Failed to start server")
	}
}
