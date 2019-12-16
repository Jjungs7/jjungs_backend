package models

import "time"

type Comment struct {
	ID int `gorm:"primary_key"`

	PostID int `sql:"index"`
	Post Post `gorm:"foreignkey:PostID;association_foreignkey:ID;"`
	Author string `gorm:"type:varchar(40)"`
	Comment string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}