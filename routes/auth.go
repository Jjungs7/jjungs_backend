package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"jjungs_backend/components/auth"
)

var JjungsPassword string

type InputPW struct {
	PW string `json:"pw"`
}

func AuthRegister(router *gin.RouterGroup) {
	router.POST("", AuthHandler)
}

func AuthHandler(c *gin.Context) {
	var input InputPW
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "ERR500",
		})
		return
	}

	if input.PW != JjungsPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	token, err := auth.GenerateToken(input.PW)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": token,
	})
}

func init() {
	JjungsPassword = os.Getenv("PASSWORD")
}
