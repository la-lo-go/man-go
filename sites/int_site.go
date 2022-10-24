package sites

import "MAPIes/models"

type IntSite interface {
	GetMangas(searchStr string, searchedMangas []models.Manga) ([]models.Manga, error)
	GetMangaPage(name string, url string) models.MangaPage
	GetChapter(name string, chapterNum int) models.Chapter
}
