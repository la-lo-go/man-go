package sites

import (
	"MAPIes/models"
	"MAPIes/sites/inManga"
	"MAPIes/sites/mangaOni"
	"MAPIes/sites/nyaa"
	// "MAPIes/sites/tuMangaNet"
)

type IntSite interface {
	// Site name returns the name of the site IN LOWERCASE
	SiteName() string

	// GetMangas returns the mangas of a site that match the search
	GetMangas(searchStr string, searchedMangas []models.Manga) ([]models.Manga, error)

	// GetMangaPage returns the info of a manga
	GetMangaPage(name string, url string) models.MangaInfo

	// GetChapter returns the info of a chapter
	GetChapter(name string, chapterNum float64) models.Chapter
}

// SitesTypes is a slice of all the supported sites.
var SitesTypes = []IntSite{
	&inManga.InManga{},
	&nyaa.Nyaa{},
	// &tuMangaNet.TuMangaNet{},
	&mangaOni.MangaOni{},
}
