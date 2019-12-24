package routes

import (
	"github.com/gin-gonic/gin"

	"jjungs_backend/models"
)

func PostRegister(router *gin.RouterGroup) {
	router.GET("", models.GetPosts)
	router.GET("/:input", models.GetPosts)
}