package main

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"jjungs_backend/components/auth"
	"jjungs_backend/routes"
)

func main() {
	g := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080", "https://api.jjungscope.co.kr", "http://api.jjungscope.co.kr"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowCredentials = true
	g.Use(cors.New(config))
	g.Use(PreHandler)

	admin := g.Group("/admin")
	authRoute := g.Group("/auth")
	boards := g.Group("/board")
	posts := g.Group("/post")

	admin.Use(OnlyAdmin)

	routes.AdminRegister(admin)
	routes.AuthRegister(authRoute)
	routes.BoardRegister(boards)
	routes.PostRegister(posts)
	g.Run()
}

func PreHandler(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	if authToken != "" {
		v, _ := auth.ValidateToken(strings.Split(authToken, " ")[1])
		if v == "JJUNGS" {
			c.Set("permissions", v)
		}
	}
	c.Next()
}

func OnlyAdmin(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		c.Abort()
		return
	}

	if _, err := auth.ValidateToken(strings.Split(authToken, " ")[1]); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		c.Abort()
		return
	}
	c.Next()
}