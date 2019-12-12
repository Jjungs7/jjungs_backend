package models

import (
	"jjungs_backend/components/database"
)

func init() {
	database.DB.AutoMigrate(
		&Post{},
		&Board{},
		&Comment{},
		&PostTag{},
	)

	database.DB.Model(&Post{}).AddForeignKey("board_id", "boards(id)", "CASCADE", "CASCADE")
	database.DB.Model(&Comment{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE")
	database.DB.Model(&PostTag{}).AddForeignKey("post_id", "posts(id)", "CASCADE", "CASCADE")
}