package models

import (
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
	Order int `gorm:"INDEX"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetBoards(c *gin.Context) {
	var whereClause = ""
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		whereClause = "read_permission <> 'JJUNGS'"
	}

	var boards []Board
	database.DB.Order("\"order\", id").Find(&boards, whereClause)
	c.JSON(200, gin.H{
		"data": boards,
	})
}

func GetBoard(c *gin.Context) {
	var board Board

	var whereClause = ""
	if permissions, _ := c.Get("permissions"); permissions != "JJUNGS" {
		whereClause = " and read_permission <> 'JJUNGS'"
	}

	boardID := c.Param("boardID")
	database.DB.First(&board, "boards.id=" + boardID + whereClause)
	if board.Name == "" {
		c.JSON(200, gin.H{
			"data": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"data": board,
	})
}

type BoardInput struct {
	ID int `json:"id"`
	Name string `json:"name"`
	URL string `json:"url"`
	ReadPermission string `json:"read"`
}

func CreateBoard(c *gin.Context) {
	var input BoardInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
			"message": err,
		})
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
			"message": errs[0],
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": board,
	})
}

func UpdateBoard(c *gin.Context) {
	var input BoardInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": err,
		})
		return
	}

	var board Board
	database.DB.First(&board, &Board{ID: input.ID})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": errs[0],
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": board,
	})
}

func DeleteBoard(c *gin.Context) {
	var boardInput BoardInput
	if err := binding.JSON.Bind(c.Request, &boardInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": err,
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ERR500",
			"message": errs[0],
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": boardInput.ID,
	})
}