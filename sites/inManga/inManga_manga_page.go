package inManga

type InMangaMangaPage struct {
	Data struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
		Result  []struct {
			PagesCount               int           `json:"pagesCount"`
			Watched                  bool          `json:"watched"`
			MangaIdentification      string        `json:"mangaIdentification"`
			MangaName                string        `json:"mangaName"`
			FriendlyMangaName        string        `json:"friendlyMangaName"`
			ID                       int           `json:"id"`
			MangaID                  int           `json:"mangaId"`
			Number                   float64       `json:"number"`
			RegistrationDate         string        `json:"registrationDate"`
			Description              string        `json:"description"`
			Pages                    []interface{} `json:"pages"`
			Identification           string        `json:"identification"`
			FeaturedChapter          bool          `json:"featuredChapter"`
			FriendlyChapterNumber    string        `json:"friendlyChapterNumber"`
			FriendlyChapterNumberURL string        `json:"friendlyChapterNumberUrl"`
		} `json:"result"`
		StatusCode int         `json:"statusCode"`
		Errors     interface{} `json:"errors"`
	} `json:"data"`
}
