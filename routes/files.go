package routes

import (
	"github.com/gin-gonic/gin"

	"jjungs_backend/models"
)

func FileRegister(router *gin.RouterGroup) {
	router.GET("/:name", models.GetFile)
}
