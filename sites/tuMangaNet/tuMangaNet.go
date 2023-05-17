package tuMangaNet

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"

	models "MAPIes/models"
	generalFunctions "MAPIes/utils"
)

const MANGA_NET_DOMAIN = "https://tumanga.net"
const MANGA_NET_MANGA_CLASS = ".post-title"

type TuMangaNet struct{}

// Returns the mangas of a site that match the search
func (tmn *TuMangaNet) GetMangas(searchValue string, searchedMangas []models.Manga) (mangas []models.Manga, err error) {
	searchStringFormated := strings.Replace(searchValue, " ", "+", -1)

	url := MANGA_NET_DOMAIN + "/?s=" + searchStringFormated + "&post_type=wp-manga&m_orderby=views"

	// Get the HTML document
	doc, err := generalFunctions.GetHtmlFromGet(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Iterate through the mangas
	doc.Find(MANGA_NET_MANGA_CLASS).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the mangas
		mangaName := s.Find("a").Text()
		mangaNameJoined, _ := generalFunctions.RemoveNonAlphanumeric(mangaName)
		mangaSite := "TuManga.net"

		// Construct the URL
		url, err := generalFunctions.RemoveNonAlphanumeric(mangaName)
		if err != nil {
			fmt.Println(err)
			return
		}

		mangaLink := MANGA_NET_DOMAIN + "/manga/" + url
		mangaChaptersNumber := 99 // TODO: hacer que no sean por defecto
		mangaCover, _ := s.Find(".img-responsive").First().Attr("src")

		// Add the manga to the list
		mangas = append(mangas, models.Manga{
			Name:           mangaName,
			NameJoined:     mangaNameJoined,
			Site:           mangaSite,
			Link:           mangaLink,
			ChaptersNumber: mangaChaptersNumber,
			Cover:          mangaCover},
		)
	})

	return mangas, nil
}

// Returns the chapters of a manga avalible in a site
func (tmn *TuMangaNet) GetMangaPage(name string, url string) (mangaPage models.MangaInfo) {
	return mangaPage
}

// Returns the pages of a chapter of a manga
func (tmn *TuMangaNet) GetChapter(name string, chapterNum float64) (chapter models.Chapter) {
	return chapter
}
