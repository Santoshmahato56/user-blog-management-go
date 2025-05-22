package controller

import "github.com/gin-gonic/gin"

// IUserController defines the interface for user controller
type IUserController interface {
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	List(ctx *gin.Context)
}
