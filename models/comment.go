package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"jjungs_backend/components/database"
)

type Comment struct {
	ID int `gorm:"primary_key"`

	PostID int `sql:"index"`
	Author string `gorm:"type:varchar(40)"`
	Comment string
	Password string `gorm:"type:varchar(20)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CommentInput struct {
	ID int `json:"id"`
	PostID int `json:"postId"`
	Comment string `json:"comment"`
	Author string `json:"author"`
	Password string `json:"pw"`
}

func GetComments(c *gin.Context) {
	var input CommentInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		fmt.Println(err)
		return
	}
	if input.PostID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}

	idString := strconv.Itoa(input.PostID)
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		row := database.DB.Table("boards").Select("boards.read_permission").Joins("inner join posts on boards.id=posts.board_id").Where("posts.id="+idString).Row()
		var result string
		row.Scan(&result)
		if result == "JJUNGS" {
			c.JSON(http.StatusOK, gin.H{
				"data": nil,
			})
			return
		}
	}

	var comments []Comment
	database.DB.Order("id asc").Where("post_id="+idString).Find(&comments)
	c.JSON(http.StatusOK, gin.H{
		"data": comments,
	})
}

func PostComment(c *gin.Context) {
	var input CommentInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		fmt.Println(err)
		return
	}
	if input.PostID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	comment := Comment{
		PostID: input.PostID,
		Comment: input.Comment,
		Author: input.Author,
		Password: input.Password,
	}

	permissions, _ := c.Get("permissions")
	if permissions != "JJUNGS" && strings.ToUpper(comment.Author) == "JJUNGS" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "ERR403",
		})
		return
	}
	if permissions == "JJUNGS" {
		comment.Author = "JJUNGS"
		comment.Password = ""
	}

	if comment.PostID <= 0 || comment.Comment == "" || (comment.Author != "JJUNGS" && (comment.Author == "" || comment.Password == "")) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	errs := database.DB.Save(&comment).GetErrors()
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		fmt.Println(errs)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": comment,
	})
}

func UpdateComment(c *gin.Context) {
	var input CommentInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		fmt.Println(err)
		return
	}
	if input.PostID <= 0 || input.ID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}

	var comment Comment
	database.DB.First(&comment, "id="+strconv.Itoa(input.ID))
	if input.PostID != comment.PostID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" && comment.Author == "JJUNGS" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	if comment.Author != "JJUNGS" && input.Password != comment.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	if input.Comment == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	if input.Comment != "" {
		comment.Comment = input.Comment
	}

	errs := database.DB.Save(&comment).GetErrors()
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		fmt.Println(errs)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": comment,
	})
}

func DeleteComment(c *gin.Context) {
	var input CommentInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		fmt.Println(err)
		return
	}
	if input.PostID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": nil,
		})
		return
	}

	var comment Comment
	database.DB.First(&comment, "id="+strconv.Itoa(input.ID))
	if comment.ID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": input.ID,
		})
		return
	}

	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" && comment.Author == "JJUNGS" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	if comment.Author != "JJUNGS" && comment.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	errs := database.DB.Delete(&comment).GetErrors()
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		fmt.Println(errs)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": input.ID,
	})
}
