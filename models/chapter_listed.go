package models

import (
	"gorm.io/gorm"
)

type ChapterListed struct {
	gorm.Model

	Name         string  `json:"name"`
	Number       float64 `json:"number"`
	LinkOriginal string  `json:"linkOriginal"`
}
