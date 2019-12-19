package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"jjungs_backend/components/database"
)

type Post struct {
	ID int `gorm:"primary_key"`

	BoardID int `sql:"index"`
	Board *Board `gorm:"foreignkey:BoardID;association_foreignkey:ID"`
	Title string `gorm:"type:varchar(255);not null"`
	Body string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type PostTag struct {
	PostID int `gorm:"INDEX"`
	Post *Post `gorm:"foreignkey:PostID;association_foreignkey:ID;"`
	Keyword string `gorm:"type:varchar(20)"`
}

func GetPosts(c *gin.Context) {
	var posts []Post
	database.DB.Order("id asc").Find(&posts)
	for idx, _ := range posts {
		posts[idx].Board = new(Board)
		database.DB.First(&posts[idx].Board, "boards.id=?", posts[idx].BoardID)
	}

	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		for i := len(posts)-1; i>=0; i-- {
			if posts[i].Board.ReadPermission == "JJUNGS" {
				posts = append(posts[:i], posts[i+1:]...)
			}
		}
	}

	c.JSON(200, gin.H{
		"data": posts,
	})
}

func GetPost(c *gin.Context) {
	var post Post
	id := c.Param("id")
	database.DB.First(&post, "posts.id=?", id)
	if post.ID == 0 {
		c.JSON(200, gin.H{
			"data": nil,
		})
		return
	}

	post.Board = &Board{}
	database.DB.First(&post.Board, "boards.id=?", post.BoardID)
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" && post.Board.ReadPermission == "JJUNGS" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": post,
	})
}

type PostInput struct {
	ID int `json:"id"`
	BoardID int `json:"boardId"`
	Title string `json:"title"`
	Body string `json:"body"`
}

func CreatePost(c *gin.Context) {
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	var input PostInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "ERR500",
		})
		fmt.Println(err)
		return
	}

	post := Post{
		Title: input.Title,
		Body: input.Body,
		BoardID: input.BoardID,
	}

	if post.Title == "" || post.BoardID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": "ERR400",
		})
		return
	}

	errs := database.DB.Save(&post).GetErrors()
	if len(errs) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": "ERR400",
		})
		fmt.Println(errs)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

func UpdatePost(c *gin.Context) {
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	var input PostInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		fmt.Println(err)
		return
	}

	var post Post
	database.DB.First(&post, "posts.id=?", input.ID)
	if post.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	if input.BoardID > 0 {
		post.BoardID = input.BoardID
	}

	if input.Title != "" {
		post.Title = input.Title
	}

	if input.Body != "" {
		post.Body = input.Body
	}

	errs := database.DB.Save(&post).GetErrors()
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		fmt.Println(errs)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

func DeletePost(c *gin.Context) {
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	var postInput PostInput
	if err := binding.JSON.Bind(c.Request, &postInput); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	if postInput.ID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": postInput.ID,
		})
		return
	}

	errs := database.DB.Delete(&Post{ID: postInput.ID}).GetErrors()
	if len(errs) > 0 {
		fmt.Println(errs)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": postInput.ID,
	})
}