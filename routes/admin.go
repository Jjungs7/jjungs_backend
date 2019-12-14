package routes

import "github.com/gin-gonic/gin"

func AdminRegister(router *gin.RouterGroup) {
	router.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Some useful functionalities",
		})
	})
}