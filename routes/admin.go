package routes

import (
	"github.com/gin-gonic/gin"

	"jjungs_backend/models"
)

func AdminRegister(router *gin.RouterGroup) {
	router.POST("/post", models.CreatePost)
	router.PUT("/post", models.UpdatePost)
	router.DELETE("/post", models.DeletePost)

	router.POST("/board", models.CreateBoard)
	router.PUT("/board", models.UpdateBoard)
	router.DELETE("/board", models.DeleteBoard)

	router.GET("/file", models.GetFileNames)
	router.POST("/file/:name", models.UploadFile)
	router.DELETE("/file", models.RemoveFile)
}
