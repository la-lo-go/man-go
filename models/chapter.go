package models

type Chapter struct {
	Name         string   `json:"name"`
	Site         string   `json:"site"`
	Number       int      `json:"number"`
	Pages_number int      `json:"pagesNumber"`
	Pages        []string `json:"pages"`
}
