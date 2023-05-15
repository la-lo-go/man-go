package mangaOni

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"

	models "MAPIes/models"
	generalFunctions "MAPIes/utils"
)

const MANGA_ONI_DOMAIN = "https://manga-oni.com"
const MANGA_ONI_MANGA_CLASS = "._135yj"
const MANGA_ONI_MANGA_IMAGE_CLASS = "._1-8M9"

type MangaOni struct{}

// Returns the mangas of a site that match the search
func (mmx *MangaOni) GetMangas(searchValue string, searchedMangas []models.Manga) (mangas []models.Manga, err error) {
	searchStringFormated := strings.Replace(searchValue, " ", "+", -1)

	url := MANGA_ONI_DOMAIN + "/buscar/?q=" + searchStringFormated

	// Get the HTML document
	doc, err := generalFunctions.GetHtmlFromGet(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Iterate through the mangas
	doc.Find(MANGA_ONI_MANGA_CLASS).Each(func(i int, s *goquery.Selection) {
		var mangaName string
		var mangaCover string

		s.Find(MANGA_ONI_MANGA_IMAGE_CLASS).Each(func(i int, sel *goquery.Selection) {
			tag := sel.Find("img").First()
			mangaName = strings.Trim(tag.AttrOr("alt", ""), " ")
			mangaCover = tag.AttrOr("src", "")
		})

		mangaSite := "MangaOni"

		mangaLink, boo := s.Find("a").First().Attr("href")
		if !boo {
			fmt.Println(err)
			return
		}

		mangaNameJoined, _ := generalFunctions.RemoveNonAlphanumeric(mangaName)

		// Add the manga to the list
		mangas = append(mangas, models.Manga{
			Name:       mangaName,
			NameJoined: mangaNameJoined,
			Site:       mangaSite,
			Link:       mangaLink,
			Cover:      mangaCover},
		)
	})

	return mangas, nil
}

// Returns the chapters of a manga avalible in a site
func (mmx *MangaOni) GetMangaPage(name string, url string) (mangaPage models.MangaInfo) {
	return mangaPage
}

// Returns the pages of a chapter of a manga
func (mmx *MangaOni) GetChapter(name string, chapterNum int) (chapter models.Chapter) {
	return chapter
}
