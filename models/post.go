package models

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model

	BoardID int `sql:"index"`
	Board Board `gorm:"foreignkey:BoardID;association_foreignkey:ID"`
	Title string `gorm:"type:varchar(255)"`
	Body string
}