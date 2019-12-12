package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model

	PostID int `sql:"index"`
	Post Post `gorm:"foreignkey:PostID;association_foreignkey:ID;"`
	Author string `gorm:"type:varchar(40)"`
	Comment string
}