package mangaMx

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"

	generalFunctions "MAPIes/general_functions"
	models "MAPIes/models"
)

const MANGA_MX_DOMAIN = "https://manga-mx.com"
const MANGA_MX_MANGA_CLASS = "._135yj"

type MangaMX struct{}

// Returns the mangas of a site that match the search
func (mmx *MangaMX) GetMangas(searchValue string, searchedMangas []models.Manga) (mangas []models.Manga, err error) {
	searchStringFormated := strings.Replace(searchValue, " ", "+", -1)

	url := MANGA_MX_DOMAIN + "/buscar/?q=" + searchStringFormated

	// Get the HTML document
	doc, err := generalFunctions.GetHtmlFromGet(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Iterate through the mangas
	doc.Find(MANGA_MX_MANGA_CLASS).Each(func(i int, s *goquery.Selection) {
		var mangaName string
		var mangaCover string

		doc.Find("._1-8M9").Each(func(i int, s *goquery.Selection) {
			tag := s.Find("img").First()
			mangaName, _ = generalFunctions.RemoveNonAlphanumeric(tag.AttrOr("alt", ""))
			mangaCover = tag.AttrOr("src", "")
		})

		mangaSite := "MangaFox"

		mangaLink, boo := s.Find("a").First().Attr("href")
		if !boo {
			fmt.Println(err)
			return
		}

		mangaChaptersNumber := 99 // TODO: hacer que no sean por defecto

		// Add the manga to the list
		mangas = append(mangas, models.Manga{
			Name:           mangaName,
			Site:           mangaSite,
			Link:           mangaLink,
			ChaptersNumber: mangaChaptersNumber,
			Cover:          mangaCover},
		)
	})

	return mangas, nil
}

// Returns the chapters of a manga avalible in a site
func (mmx *MangaMX) GetMangaPage(name string, url string) (mangaPage models.MangaPage) {
	return mangaPage
}

// Returns the pages of a chapter of a manga
func (mmx *MangaMX) GetChapter(name string, chapterNum int) (chapter models.Chapter) {
	return chapter
}
