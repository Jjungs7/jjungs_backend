package models

import (
	//"github.com/jinzhu/gorm"
)

type PostTag struct {
	PostID int `gorm:"INDEX"`
	Post Post `gorm:"foreignkey:PostID;association_foreignkey:ID;"`
	Keyword string `gorm:"type:varchar(20)"`
}