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
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"}
	g.Use(cors.New(config))
	g.Use(PreHandler)

	admin := g.Group("/admin")
	auth := g.Group("/auth")
	boards := g.Group("/board")
	comments := g.Group("/comment")
	posts := g.Group("/post")

	admin.Use(OnlyAdmin)

	// database.DB.Create(&models.Post{
	// 	Model:   gorm.Model{},
	// 	BoardID: 0,
	// 	Board:   models.Board{},
	// 	Title:   "",
	// 	Body:    "",
	// })

	routes.AdminRegister(admin)
	routes.AuthRegister(auth)
	routes.BoardRegister(boards)
	routes.CommentRegister(comments)
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