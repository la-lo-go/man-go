package models

type MangaLinksClustered struct {
	Name       string     `json:"name"`
	SitesLinks []SiteLink `json:"sitesLinks"`
}

type SiteLink struct {
	Site string `json:"site"`
	Link string `json:"link"`
}
