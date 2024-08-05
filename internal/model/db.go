package model

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var MyDB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("data/gorm.db"), &gorm.Config{})
	log.Println("db: {}", db.Name())
	if err != nil {
		log.Fatalln("db error: {}", err)
	}
	MyDB = db
}
