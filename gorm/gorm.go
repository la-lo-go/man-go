package gorm

import (
	"MAPIes/models"
	"MAPIes/utils"
	"errors"
	"log"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
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

// Find the search in the database and return the search and a status string.
// Possible status: "Exact", "Partial", "NotFound"
func FindSearch(searchStr string) ([]models.Manga, string) {
	var search models.Search

	// Format the search to match the database format
	searchStr, _ = utils.RemoveNonAlphanumeric(searchStr)

	// Search for exact match
	err := db.Model(&models.Search{}).Preload("Mangas").Where("search = ?", searchStr).First(&search).Error
	if err == nil {
		log.Println("\n>>>> [gorm/searches]: Exact match found")
		return search.Mangas, "Exact"
	}

	// Search for partial match
	var searches []models.Search
	err = db.Model(&models.Search{}).Preload("Mangas").Find(&searches).Error
	if err == nil {
		for _, s := range searches {
			if strings.Contains(searchStr, s.Search) {
				log.Println("\n>>>> [gorm/searches]: Partial match found")
				return s.Mangas, "Partial"
			}
		}
	}

	// If the searchStr is not found in the database
	return nil, "NotFound"
}

func DumpSearchToDB(searchStr string, mangas []models.Manga) {
	searchStrFormatted, _ := utils.RemoveNonAlphanumeric(searchStr)

	search := models.Search{
		Search: searchStrFormatted,
		Mangas: mangas,
	}

	err := db.Create(&search).Error
	if err != nil {
		log.Println(err)
	}
}

// Search the link of the manga in the database based on the name and site and return the link
func SearchMangaURL(mangaName string, siteName string) string {
	var manga models.Manga

	// Format the manga name and site name to match the database format
	mangaNameFormatted, _ := utils.RemoveNonAlphanumeric(mangaName)
	siteNameFormatted := strings.ToLower(siteName)

	err := db.Where("name_joined = ? AND site = ?", mangaNameFormatted, siteNameFormatted).First(&manga).Error
	if err == nil {
		return manga.Link
	}

	// If there is no coincidence, return an empty string
	return ""
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

func FindInMangaChapterID(name string, chapterNum float64) (chapter models.InMangaChapter) {
	manga, err := SearchInManga(name)
	if err != nil {
		return chapter
	}

	db.Where("in_manga_manga_id = ? AND number = ?", manga.ID, chapterNum).First(&chapter)

	return chapter
}
