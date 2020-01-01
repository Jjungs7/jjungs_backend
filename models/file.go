package models

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"jjungs_backend/components/database"
)

type File struct {
	ID int
	Name string `gorm:"type:varchar(256);UNIQUE_INDEX;not null"`
	Size int64

	CreatedAt time.Time
}

var baseDir string

func isOffensiveName(name string) bool {
	if strings.Contains(name, "..") {
		return true
	}
	return false
}

func GetFileNames(c *gin.Context) {
	var files []File
	database.DB.Find(&files)
	c.JSON(http.StatusOK, gin.H{
		"data": files,
	})
}

func GetFile(c *gin.Context) {
	name := c.Param("name")
	if name == "" || isOffensiveName(name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": "/file/:name should be provided",
		})
		return
	}

	fileName := baseDir + "/" + name
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": err,
		})
		return
	}

	c.File(fileName)
}

func UploadFile(c *gin.Context) {
	name := c.Param("name")
	if name == "" || isOffensiveName(name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": "/file/:name should be provided",
		})
		return
	}

	fileName := baseDir + "/" + name
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": err,
		})
		return
	}

	if _, err := os.Stat(fileName); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "FIL000",
			"message": "file exists",
		})
		return
	}

	if err := c.SaveUploadedFile(file, fileName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR404",
			"message": "incorrect file or filename. Check if file exceeds 32MiB",
		})
		return
	}

	// Filename uniqueness is checked above
	database.DB.Save(&File{Name: file.Filename, Size: file.Size})
	c.JSON(http.StatusOK, gin.H{
		"data": name,
	})
}

func RemoveFile(c *gin.Context) {
	type FileInput struct {
		Name string `json:"fileName"`
	}

	var input FileInput
	if err := binding.JSON.Bind(c.Request, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": err,
		})
		return
	}

	if input.Name == "" || isOffensiveName(input.Name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": "field \"fileName\" has to be provided",
		})
		return
	}

	fileName := baseDir + "/" + input.Name
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": "file not found",
		})
		return
	}

	if err := os.Remove(fileName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ERR400",
			"message": err,
		})
		return
	}
	database.DB.Delete(&File{Name: input.Name})
	c.JSON(http.StatusOK, gin.H{
		"data": input.Name,
	})
}

func init() {
	baseDir = os.Getenv("FILES_DIR")
}