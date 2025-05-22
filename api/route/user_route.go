package route

import (
	"github.com/gin-gonic/gin"
	"github.com/userblog/management/api/controller"
	"github.com/userblog/management/api/middleware"
)

type UserRoute struct {
	userController controller.IUserController
	authMiddleware middleware.IAuthMiddleware
}

func NewUserRoute(userController controller.IUserController, authMiddleware middleware.IAuthMiddleware) UserRoute {
	return UserRoute{
		userController: userController,
		authMiddleware: authMiddleware,
	}
}

func (r UserRoute) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("/users")
	router.Use(r.authMiddleware.JWTAuth())

	router.GET("", r.authMiddleware.RequirePermission("user", "read"), r.userController.List)
	router.GET("/:id", r.authMiddleware.RequirePermission("user", "read"), r.userController.GetByID)
	router.POST("", r.authMiddleware.RequirePermission("user", "create"), r.userController.Create)
	router.PUT("/:id", r.authMiddleware.RequirePermission("user", "update"), r.userController.Update)
	router.DELETE("/:id", r.authMiddleware.RequirePermission("user", "delete"), r.userController.Delete)
}
