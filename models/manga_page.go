package models

type MangaPage struct {
	Name           string          `json:"name"`
	Site           string          `json:"site"`
	Cover          string          `json:"cover"`
	ChaptersNumber int             `json:"chaptersNumber"`
	ChaptersListed []ChapterListed `json:"chaptersListed"`
}
