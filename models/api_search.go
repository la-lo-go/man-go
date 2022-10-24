package models

import (
	"time"
)

type ApiSearch struct {
	Search   string    `json:"search"`
	Date     time.Time `json:"date"`
	Response []Manga   `json:"response"`
}