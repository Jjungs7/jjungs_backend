package models

import (
	"jjungs_backend/components/database"
)

func init() {
	database.DB.AutoMigrate(
		&Post{},
		&Board{},
	)
}