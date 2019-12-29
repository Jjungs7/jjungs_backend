package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"jjungs_backend/components/auth"
)

var JjungsPW string

type InputPW struct {
	PW string `json:"pw"`
}

func AuthRegister(router *gin.RouterGroup) {
	router.POST("/val", ValidateAuth)
	router.POST("", Authenticate)
}

func ValidateAuth(c *gin.Context) {
	type InputToken struct {
		Token string `json:"token"`
	}

	var input InputToken
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		return
	}

	_, err := auth.ValidateToken(input.Token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "TKN001",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "ok",
	})
}

func Authenticate(c *gin.Context) {
	type InputPW struct {
		PW string `json:"pw"`
	}

	var input InputPW
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		return
	}

	if input.PW != JjungsPW {
		c.JSON(http.StatusOK, gin.H{
			"error": "TKN000",
		})
		return
	}

	token, err := auth.GenerateToken(input.PW)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "TKN000",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": token,
	})
}

func init() {
	JjungsPW = os.Getenv("PASSWORD")
}
