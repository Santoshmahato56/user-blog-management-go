package controller

import "github.com/gin-gonic/gin"

// IAuthController defines the interface for authentication controller
type IAuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetMe(ctx *gin.Context)
}
