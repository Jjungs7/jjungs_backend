package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"jjungs_backend/routes"
)

func main() {
	g := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	g.Use(cors.New(config))
	auth := g.Group("/auth")
	posts := g.Group("/post")
	boards := g.Group("/board")
	comments := g.Group("/comment")

	routes.AuthRegister(auth)
	routes.PostRegister(posts)
	routes.BoardRegister(boards)
	routes.CommentRegister(comments)
	g.Run()
}
