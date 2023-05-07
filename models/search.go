package models

import "gorm.io/gorm"

type Search struct {
	gorm.Model

	Search string  `json:"search"`
	Mangas []Manga `json:"mangas" gorm:"many2many:search_mangas;"`
}
