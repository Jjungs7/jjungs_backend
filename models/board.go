package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"jjungs_backend/components/database"
)

type Board struct {
	ID int
	Name string `gorm:"type:varchar(40);UNIQUE_INDEX;not null"`
	URL string `gorm:"type:varchar(40);UNIQUE_INDEX;not null"`
	ReadPermission string `gorm:"type:varchar(20);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetBoards(c *gin.Context) {
	var whereClause = ""
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		whereClause = " read_permission <> 'JJUNGS'"
	}

	var boards []Board
	database.DB.Order("id asc").Find(&boards, whereClause)
	c.JSON(200, gin.H{
		"data": boards,
	})
}

func GetBoard(c *gin.Context) {
	type Posts struct {
		Board Board
		Posts []Post
	}

	var whereClause = ""
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		whereClause = " and read_permission <> 'JJUNGS'"
	}

	var posts Posts
	id := c.Param("id")
	database.DB.First(&posts.Board, "boards.id=" + id + whereClause)
	if posts.Board.Name == "" {
		c.JSON(200, gin.H{
			"data": nil,
		})
		return
	}

	database.DB.Order("id asc").Find(&posts.Posts, "posts.board_id=?", posts.Board.ID)
	c.JSON(200, gin.H{
		"data": posts,
	})
}

type BoardInput struct {
	ID int `json:"id"`
	Name string `json:"name"`
	URL string `json:"url"`
	ReadPermission string `json:"read"`
}

func CreateBoard(c *gin.Context) {
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	var input BoardInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
		fmt.Println(err)
		return
	}

	board := Board{
		Name: input.Name,
		URL: input.URL,
		ReadPermission: input.ReadPermission,
	}

	if board.Name == "" || board.URL == "" || board.ReadPermission == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	errs := database.DB.Save(&board).GetErrors()
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		fmt.Println(errs)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": board,
	})
}

func UpdateBoard(c *gin.Context) {
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	var input BoardInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	var board Board
	database.DB.First(&board, Board{ID: input.ID})
	if board.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	if input.Name != "" {
		board.Name = input.Name
	}

	if input.URL != "" {
		board.URL = input.URL
	}

	if input.ReadPermission != "" {
		board.ReadPermission = input.ReadPermission
	}

	errs := database.DB.Save(&board).GetErrors()
	if len(errs) > 0 {
		fmt.Println(errs)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": board,
	})
}

func DeleteBoard(c *gin.Context) {
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ERR401",
		})
		return
	}

	var boardInput BoardInput
	if err := binding.JSON.Bind(c.Request, &boardInput); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
		})
		return
	}

	if boardInput.ID <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": boardInput.ID,
		})
		return
	}

	errs := database.DB.Delete(&Board{ID: boardInput.ID}).GetErrors()
	if len(errs) > 0 {
		fmt.Println(errs)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": boardInput.ID,
	})
}