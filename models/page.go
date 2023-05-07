package models

import (
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model

	Number    int    `json:"number"`
	Link      string `json:"link"`
	ChapterID int    `json:"chapterID"`
}