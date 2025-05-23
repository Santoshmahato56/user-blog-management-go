package main

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	controllerImpl "github.com/userblog/management/api/controller/impl"
	"github.com/userblog/management/api/middleware"
	middlewareImpl "github.com/userblog/management/api/middleware"
	"github.com/userblog/management/api/route"
	repoImpl "github.com/userblog/management/internal/repository/impl"
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
	//initializeDatabaseScript(database)
	logger.Info(ctx, "Database schema initialized and seeded with roles and permissions")

	// Initialize repositories with the database connection
	var userRepo = repoImpl.NewUserRepository(database)
	var blogRepo = repoImpl.NewBlogRepository(database)

	// Initialize services
	var authService = serviceImpl.NewAuthService(userRepo)
	var userService = serviceImpl.NewUserService(userRepo)
	var blogService = serviceImpl.NewBlogService(blogRepo)

	// Initialize middleware
	var authMiddleware middleware.IAuthMiddleware = middlewareImpl.NewAuthMiddleware(authService)

	// Initialize controllers
	var authController = controllerImpl.NewAuthController(authService)
	var userController = controllerImpl.NewUserController(userService)
	var blogController = controllerImpl.NewBlogController(blogService)

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
