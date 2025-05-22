package controller

import "github.com/gin-gonic/gin"

// IBlogController defines the interface for blog controller
type IBlogController interface {
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	List(ctx *gin.Context)
	ListByUser(ctx *gin.Context)
}
