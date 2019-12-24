package routes

import (
	"github.com/gin-gonic/gin"

	"jjungs_backend/models"
)

func BoardRegister(router *gin.RouterGroup) {
	router.GET("", models.GetBoards)
	router.GET("/:url", models.GetBoard)
}
