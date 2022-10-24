package nyaa

type NyaaMangaPage struct {
	MangaName string `json:"mangaName"`
	Completed string `json:"completed"`
	Chs       []struct {
		OrderNumber float64 `json:"orderNumber"`
		Pages       int     `json:"pages"`
	} `json:"chs"`
}
