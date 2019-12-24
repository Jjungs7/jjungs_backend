package models

import (
	"jjungs_backend/components/database"
)

func init() {
	database.DB.AutoMigrate(
		&Board{},
		&Post{},
		&PostTag{},
	)

	database.DB.Model(&Post{}).AddForeignKey("board_id", "boards(id)", "SET NULL", "CASCADE")
	database.DB.Model(&PostTag{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE")
}