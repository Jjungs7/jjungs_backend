package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID int `gorm:"primary_key"`

	BoardID int `sql:"index"`
	Board Board `gorm:"foreignkey:BoardID;association_foreignkey:ID"`
	Title string `gorm:"type:varchar(255)"`
	Body string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func GetPosts(c *gin.Context) {

}

func GetPost(c *gin.Context) {

}

func CreatePost(c *gin.Context) {

}

func DeletePost(c *gin.Context) {

}

func UpdatePost(c *gin.Context) {

}