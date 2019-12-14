package models

import (
	"github.com/jinzhu/gorm"
)

type Board struct {
	gorm.Model

	Name string `gorm:"type:varchar(40);UNIQUE_INDEX;not null"`
	URL string `gorm:"type:varchar(40);UNIQUE_INDEX;not null"`
	ReadPermission string `gorm:"type:varchar(20);not null"`
}