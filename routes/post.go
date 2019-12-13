package routes

import "github.com/gin-gonic/gin"

func PostRegister(router *gin.RouterGroup) {
	router.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get",
		})
	})

	router.GET("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get11",
		})
	})

	router.POST("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "post",
		})
	})

	router.PUT("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "put",
		})
	})

	router.DELETE("/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "del",
		})
	})
}