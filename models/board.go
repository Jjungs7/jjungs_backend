package models

import (
	"github.com/jinzhu/gorm"
)

type ReadPermissions string
type Board struct {
	gorm.Model

	Name string `gorm:"type:varchar(40);UNIQUE_INDEX"`
	URL string `gorm:"type:varchar(40);UNIQUE_INDEX"`
	ReadPermission string `gorm:"type:varchar(20)"`
}

const (
	JJUNGS ReadPermissions = "JJUNGS"
	PUBLIC ReadPermissions = "PUBLIC"
)