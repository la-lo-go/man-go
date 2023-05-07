package models

import (
	"gorm.io/gorm"
)

type MangaInfo struct {
	gorm.Model

	Name           string          `json:"name"`
	Site           string          `json:"site"`
	Cover          string          `json:"cover"`
	ChaptersNumber int             `json:"chaptersNumber"`
	ChaptersListed []ChapterListed `json:"chaptersListed"`
}
