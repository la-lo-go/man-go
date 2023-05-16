package gorm

import (
	"MAPIes/models"
	"errors"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(sqlite.Open("mango.db?_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database" + err.Error())
	}

	db.Logger = db.Logger.LogMode(3)

	err = db.AutoMigrate(
		&models.Search{},
		&models.Manga{},
		&models.Chapter{},
		&models.Page{},
		&models.InMangaManga{},
		&models.InMangaChapter{},
		&models.InMangaPage{},
	)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Database initialized")
}

func getDB() *gorm.DB {
	return db
}

func SearchInManga(name string) (manga models.InMangaManga, err error) {
	result := db.Where("manga_name = ?", name).First(&manga)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return manga, nil // No error, but no manga found
		}
		return manga, result.Error // Any other error occurred
	}

	return manga, nil
}

func AddInManga(manga models.InMangaManga) (err error) {
	result := db.Create(&manga)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateInManga(manga models.InMangaManga) (err error) {
	result := db.Save(&manga)

	if result.Error != nil {
		return result.Error
	}

	db.Commit()
	return nil
}
