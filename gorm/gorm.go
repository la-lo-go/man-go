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
func SearchMangaLink(siteName string, mangaName string) string {
	manga, err := SearchManga(siteName, mangaName)
	if err != nil {
		log.Println("Error searching manga in database: ", err)
	}

	if manga.Name != "" {
		return manga.Link
	}

	// If there is no coincidence, return an empty string
	return ""
}

func SearchManga(site string, name string) (manga models.Manga, err error) {
	site = strings.ToLower(site)
	nameJoined, err := utils.RemoveNonAlphanumeric(name)
	if err != nil {
		return manga, err
	}

	result := db.Where("site = ? AND name_joined = ?", site, nameJoined).First(&manga)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return manga, nil // No error, but no manga found
		}
		return manga, result.Error // Any other error occurred
	}

	return manga, nil
}

func DumpMangaToDB(site string, page models.Manga) error {
	// Search for the manga if is already there
	manga, err := SearchManga(site, page.Name)
	if err != nil {
		return err
	}

	if manga.Name == "" { // If the manga is not there, add it
		manga.Name, _ = utils.RemoveNonAlphanumeric(page.Name)
		manga.Chapters = page.Chapters

		err = AddManga(manga)
		if err != nil {
			return err
		}
	} else { // If the manga is there, confirm that the chapters are the same
		
		// If the chapters are the same, return
		if len(manga.Chapters) == len(page.Chapters) {
			return nil
		}

		manga.Chapters = page.Chapters
		err = UpdateManga(manga)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddManga(manga models.Manga) (err error) {
	result := db.Create(&manga)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateManga(manga models.Manga) (err error) {
	result := db.Save(&manga)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FindChapterWebID(site string, name string, chapterNum float64) (chapter models.Chapter) {
    manga, err := SearchManga(site, name)
    if err != nil {
        return chapter
    }

    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    result := tx.Where("manga_id = ? AND number = ?", manga.ID, chapterNum).First(&chapter)

    if result.Error != nil {
        tx.Rollback()
        return chapter
    }

    tx.Commit()
    return chapter
}
