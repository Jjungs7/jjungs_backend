package models

import (
	"jjungs_backend/components/database"
)

func init() {
	database.DB.AutoMigrate(
		&Board{},
		&Post{},
		&PostTag{},
		&File{},
	)

	database.DB.Model(&Post{}).AddForeignKey("board_id", "boards(id)", "CASCADE", "CASCADE")
	database.DB.Model(&PostTag{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE")
}