package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"jjungs_backend/components/database"
	"jjungs_backend/models"
)

func BoardRegister(router *gin.RouterGroup) {
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

	router.POST("", CreateBoard)

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

type BoardInput struct {
	Name string `json:"name"`
	URL string `json:"url"`
	ReadPermission string `json:"read"`
}

func CreateBoard(c *gin.Context) {
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR403",
		})
		return
	}

	var input BoardInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "ERR500",
		})
		return
	}
	board := models.Board{
		Name: input.Name,
		URL: input.URL,
		ReadPermission: input.ReadPermission,
	}
	if board.Name == "" || board.URL == "" || board.ReadPermission == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "ERR400",
		})
		return
	}
	errs := database.DB.Save(&board).GetErrors()
	if len(errs) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": "ERR400",
		})
		fmt.Println(errs)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": []models.Board{
			board,
		},
	})
}