package routes

import (
	"github.com/gin-gonic/gin"

	"jjungs_backend/models"
)

func CommentRegister(router *gin.RouterGroup) {
	router.GET("", models.GetComments)
	router.POST("", models.PostComment)
	router.PUT("", models.UpdateComment)
	router.DELETE("", models.DeleteComment)
}