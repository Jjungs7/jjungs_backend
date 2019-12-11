package models

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model

	BoardId int `gorm:"foreignkey:id;association_foreignkey:board_id;INDEX"`
	Title string
	Body string
}
