package models

import "gorm.io/gorm"

type Manga struct {
	gorm.Model

	Name           string      `json:"name"`
	NameJoined     string      `json:"name_joined"` // Name joined with dashes
	Site           string      `json:"site"`
	Link           string      `json:"link"`
	ChaptersNumber int         `json:"chapters_number" `
	Cover          string      `json:"cover"`
	Searches       []Search `json:"Searches" gorm:"many2many:search_mangas;"`
	Chapters       []Chapter   `json:"chapters"`
}
