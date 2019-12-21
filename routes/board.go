package routes

import (
	"github.com/gin-gonic/gin"

	"jjungs_backend/models"
)

func BoardRegister(router *gin.RouterGroup) {
	router.GET("", models.GetBoards)
	router.GET("/:name", models.GetBoard)
	router.POST("", models.CreateBoard)
	router.PUT("", models.UpdateBoard)
	router.DELETE("", models.DeleteBoard)
}
