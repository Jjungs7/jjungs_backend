package models

import (
	"github.com/jinzhu/gorm"
)

type ReadPermissions int
type Board struct {
	gorm.Model

	Name string `gorm:"type:varchar(40);UNIQUE_INDEX"`
	URL string `gorm:"type:varchar(40);UNIQUE_INDEX"`
	ReadPermission int
}

const (
	JJUNGS ReadPermissions = iota
	PUBLIC
)