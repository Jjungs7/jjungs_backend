package database

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func init() {
	_DB, err := gorm.Open("postgres",
		"host="+os.Getenv("DB_HOST")+
		" port=5432"+
		" user="+os.Getenv("DB_USER")+
		" password="+os.Getenv("DB_PW")+
		" dbname="+os.Getenv("DB_DB")+
		" sslmode=disable",
	)

	if err != nil {
		panic(err)
	}
	DB = _DB
}