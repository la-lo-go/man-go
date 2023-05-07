package gorm

import (
	"MAPIes/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func Init() {
	db, err := gorm.Open(sqlite.Open("my.db?_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database" + err.Error())
	}

	db.Logger = db.Logger.LogMode(3)

	err = db.AutoMigrate(
		&models.Search{},
		&models.Manga{},
		&models.Chapter{},
		&models.Page{},
	)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Database initialized")
}
