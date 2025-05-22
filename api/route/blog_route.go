package route

import (
	"github.com/gin-gonic/gin"
	"github.com/userblog/management/api/controller"
	"github.com/userblog/management/api/middleware"
)

type BlogRoute struct {
	blogController controller.IBlogController
	authMiddleware middleware.IAuthMiddleware
}

func NewBlogRoute(blogController controller.IBlogController, authMiddleware middleware.IAuthMiddleware) BlogRoute {
	return BlogRoute{
		blogController: blogController,
		authMiddleware: authMiddleware,
	}
}

func (r BlogRoute) BlogRoute(rg *gin.RouterGroup) {
	router := rg.Group("/blogs")

	// Public routes
	router.GET("", r.blogController.List)
	router.GET("/:id", r.blogController.GetByID)
	router.GET("/user/:user_id", r.blogController.ListByUser)

	// Protected routes
	authRouter := router.Group("")
	authRouter.Use(r.authMiddleware.JWTAuth())

	authRouter.POST("", r.authMiddleware.RequirePermission("blog", "create"), r.blogController.Create)
	authRouter.PUT("/:id", r.authMiddleware.RequirePermission("blog", "update"), r.blogController.Update)
	authRouter.DELETE("/:id", r.authMiddleware.RequirePermission("blog", "delete"), r.blogController.Delete)
}
