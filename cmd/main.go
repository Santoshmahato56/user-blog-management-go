package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	controllerImpl "github.com/userblog/management/api/controller/impl"
	"github.com/userblog/management/api/middleware"
	middlewareImpl "github.com/userblog/management/api/middleware"
	"github.com/userblog/management/api/route"
	repoImpl "github.com/userblog/management/internal/repository/impl"
	serviceImpl "github.com/userblog/management/internal/service/impl"
	"github.com/userblog/management/pkg/config"
	"github.com/userblog/management/pkg/db"
	"github.com/userblog/management/pkg/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	startTime := time.Now()

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
	initializeDatabaseScript(ctx, database)
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
	startServerWithGracefulShutdown(ctx, router, startTime)
}

func startServerWithGracefulShutdown(ctx context.Context, router *gin.Engine, startTime time.Time) {
	port := config.GetOrDefaultString("PORT", "8080")

	logger.InfoF(ctx, "Server starting on port: %s", port)

	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router.Handler(),
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.FatalF(ctx, "❌ Failed to bind to %s: %v", addr, err)
	}

	elapsedTime := time.Since(startTime)
	logger.InfoF(ctx, "✅ Server is ready to accept connections on %s (started in %.2f seconds)", addr, elapsedTime.Seconds())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start a server in a separate goroutine
	go func() {
		if err = srv.Serve(listener); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logger.FatalF(ctx, "❌ Server failed: %v", err)
		}
	}()

	// Wait for the interrupt signal
	<-quit
	logger.InfoF(ctx, "⚠️ Shutdown signal received, shutting down server gracefully...")
	shutdownStart := time.Now()

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.FatalF(ctx, "❌ Server forced to shutdown: %v", err)
	}

	shutdownDuration := time.Since(shutdownStart).Seconds()
	logger.InfoF(ctx, "✅ Server exited gracefully in %.2f seconds", shutdownDuration)
}
