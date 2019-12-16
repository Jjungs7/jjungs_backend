package routes

import (
	"github.com/gin-gonic/gin"

	"jjungs_backend/models"
)

func PostRegister(router *gin.RouterGroup) {
	router.GET("", models.GetPosts)
	router.GET("/:id", models.GetPost)
	router.POST("", models.CreatePost)
	router.PUT("", models.UpdatePost)
	router.DELETE("", models.DeletePost)
}