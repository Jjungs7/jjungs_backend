package routes

import (
	"github.com/gin-gonic/gin"

	"jjungs_backend/models"
)

func BoardRegister(router *gin.RouterGroup) {
	router.GET("", models.GetBoards)
	router.GET("/:id", models.GetBoard)
	router.POST("", models.CreateBoard)
	router.PUT("/:id", models.UpdateBoard)
	router.DELETE("/:id", models.DeleteBoard)
}
