package routes

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"

	"jjungs_backend/components/auth"
)

var JJUNGS_PASSWORD string

type InputPW struct {
	pw string `json:"pw"`
}

func AuthRegister(router *gin.RouterGroup) {
	router.POST("", AuthHandler)
}

func AuthHandler(c *gin.Context) {
	var input InputPW
	c.BindJSON(&input)
	fmt.Println(input)
	fmt.Println(c.Params)
	fmt.Println(c.Keys)

	if input.pw != JJUNGS_PASSWORD {
		c.JSON(200, gin.H{
			"message": "post",
		})

		return
	}

	token, err := auth.GenerateToken(input.pw)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"message": token,
	})
}

func init() {
	JJUNGS_PASSWORD = os.Getenv("PASSWORD")
}
