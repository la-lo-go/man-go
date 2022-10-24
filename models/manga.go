package models

type Manga struct {
	Name            string `json:"name"`
	Site            string `json:"site"`
	Link            string `json:"link"`
	ChaptersNumber  int    `json:"chapters_number"`
	Cover           string `json:"cover"`
}
