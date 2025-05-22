package route

import (
	"github.com/gin-gonic/gin"
	"github.com/userblog/management/api/controller"
	"github.com/userblog/management/api/middleware"
)

type AuthRoute struct {
	authController controller.IAuthController
	authMiddleware middleware.IAuthMiddleware
}

func NewAuthRoute(authController controller.IAuthController, authMiddleware middleware.IAuthMiddleware) AuthRoute {
	return AuthRoute{
		authController: authController,
		authMiddleware: authMiddleware,
	}
}

func (r AuthRoute) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", r.authController.Register)
	router.POST("/login", r.authController.Login)
	router.GET("/me", r.authMiddleware.JWTAuth(), r.authController.GetMe)
}
