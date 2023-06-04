package models

import "gorm.io/gorm"

type Manga struct {
	gorm.Model

	Name           string    `json:"name"`
	NameJoined     string    `json:"nameJoined"` // Name joined with dashes
	Site           string    `json:"site"`
	Link           string    `json:"link"`
	ChaptersNumber int       `json:"chaptersNumber"`
	Cover          string    `json:"cover"`
	WebID          string    `json:"webID"`
	Searches       []Search  `json:"Searches" gorm:"many2many:search_mangas;"`
	Chapters       []Chapter `json:"chapters"`
}
