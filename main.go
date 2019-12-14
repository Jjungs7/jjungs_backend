package main

import (
	"fmt"
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
	config.AllowAllOrigins = true
	g.Use(cors.New(config))
	admin := g.Group("/admin")
	auth := g.Group("/auth")
	boards := g.Group("/board")
	comments := g.Group("/comment")
	posts := g.Group("/post")

	admin.Use(OnlyAdmin)

	routes.AdminRegister(admin)
	routes.AuthRegister(auth)
	routes.BoardRegister(boards)
	routes.CommentRegister(comments)
	routes.PostRegister(posts)
	g.Run()
}

func OnlyAdmin(c *gin.Context) {
	fmt.Println("test")
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