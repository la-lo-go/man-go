package models

import (
	"gorm.io/gorm"
)

type InMangaManga struct {
	gorm.Model

	MangaName string           `json:"mangaName"`
	Chapters  []InMangaChapter `json:"chapters"`
}

type InMangaChapter struct {
	gorm.Model

	Number         float64       `json:"number"`
	PagesCount     int           `json:"pagesCount"`
	ID             string        `json:"id"`
	Pages          []InMangaPage `json:"pages"`
	InMangaMangaID int           `json:"inMangaMangaID"`
}

type InMangaPage struct {
	gorm.Model

	Number           int    `json:"number"`
	ID               string `json:"id"`
	InMangaChapterID int    `json:"inMangaChapterID"`
}
