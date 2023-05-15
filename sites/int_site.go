package sites

import (
	"MAPIes/models"
	"MAPIes/sites/inManga"
	"MAPIes/sites/mangaOni"
	"MAPIes/sites/nyaa"
	"MAPIes/sites/tuMangaNet"
)

type IntSite interface {
	GetMangas(searchStr string, searchedMangas []models.Manga) ([]models.Manga, error)
	GetMangaPage(name string, url string) models.MangaInfo
	GetChapter(name string, chapterNum int) models.Chapter
}

// SitesTypes is a slice of all the supported sites.
var SitesTypes = []IntSite{
	&inManga.InManga{},
	&nyaa.Nyaa{},
	&tuMangaNet.TuMangaNet{},
	&mangaOni.MangaOni{},
}
