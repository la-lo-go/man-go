package models

import (
	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model

	Name         string  `json:"name"`
	Site         string  `json:"site"`
	Number       float64 `json:"number"`
	PagesNumber  int     `json:"pagesNumber"`
	LinkOriginal string  `json:"linkOriginal"`
	WebID        string  `json:"webID"`
	Pages        []Page  `json:"pages"`
	MangaID      int     `json:"mangaID"` // Foreign key
}
