package inManga

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sort"
	// "io/ioutil"
	// "net/http"
	"strconv"
	"strings"

	"MAPIes/general_functions"
	"MAPIes/models"
)

type InManga struct{}

const INMANGA_THUMBNAIL_URL = "https://pack-yak.intomanga.com/thumbnails/manga/"
const INMANGA_GET_ALL_URL = "https://inmanga.com/chapter/getall?mangaIdentification="

// GetMangas Returns the mangas of a site that match the search
func (in *InManga) GetMangas(searchValue string, searchedMangas []models.Manga) (mangas []models.Manga, err error) {
	searchStringFormated := strings.Replace(searchValue, " ", "+", -1)
	url := "https://inmanga.com/manga/getMangasConsultResult"
	payload := strings.NewReader(`filter%5Bgeneres%5D%5B%5D=-1&filter%5BqueryString%5D=+` + searchStringFormated + `+&filter%5Bskip%5D=0&filter%5Btake%5D=10&filter%5Bsortby%5D=1&filter%5BbroadcastStatus%5D=0&filter%5BonlyFavorites%5D=false&d=`)

	fmt.Println("\nsearchStringFormated=" + searchStringFormated)

	// Consigue el HTML de la pagina
	doc, err := general_functions.GetHtmlFromPost(url, payload)
	if err != nil {
		fmt.Println(err)
		return []models.Manga{}, err
	}

	doc.Find(".manga-result").Each(func(i int, s *goquery.Selection) {
		// Get the manga name with ID
		mangaCoverLinkRelative, _ := s.Find("img").Attr("data-src")
		mangaID := strings.Split(mangaCoverLinkRelative, "/thumbnails/manga/")[1]

		// Get manga attributes
		mangaName, _ := general_functions.RemoveNonAlphanumeric(strings.Trim(s.Find(".m0").First().Text(), " "))
		mangaSite := "InManga"
		mangaLink := "https://inmanga.com/ver/manga/" + mangaID
		mangaChaptersNumber, _ := strconv.Atoi(s.Find(".icon-info text-muted").First().Text())
		mangaCover := INMANGA_THUMBNAIL_URL + mangaID

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

// GetMangaPage Returns the chapters of a manga avalible in a site
func (in *InManga) GetMangaPage(name string, url string) (mangaPage models.MangaPage) {
	jsonResponse := InMangaMangaPage{}

	urlSplit := strings.Split(url, "/")
	mangaID := urlSplit[len(urlSplit)-1]

	urlRequest := INMANGA_GET_ALL_URL + mangaID
	response, err := general_functions.GetJsonFromGet(urlRequest)

	if err != nil {
		fmt.Println(err)
		return mangaPage
	}

	response = strings.Replace(response, `"{`, `{`, 1)
	response = strings.Replace(response, `}"`, `}`, 1)
	response = strings.Replace(response, `\`, ``, -1)

	err = json.Unmarshal([]byte(response), &jsonResponse)
	if err != nil {
		fmt.Println(err)
		return mangaPage
	}

	mangaPage.Name = name
	mangaPage.Site = "InManga"
	mangaPage.Cover = INMANGA_THUMBNAIL_URL + mangaID // TODO: Ver si esta bien

	for _, chapter := range jsonResponse.Data.Result {
		mangaPage.ChaptersListed = append(mangaPage.ChaptersListed, models.ChapterListed{
			Number:       chapter.Number,
			LinkOriginal: "https://inmanga.com/ver/manga/" + name + "/" + chapter.FriendlyChapterNumberURL + "/" + chapter.Identification,
		})
	}

	mangaPage.ChaptersNumber = len(mangaPage.ChaptersListed)

	// Sort the chapters
	sort.Slice(mangaPage.ChaptersListed, func(i, j int) bool {
		return mangaPage.ChaptersListed[i].Number < mangaPage.ChaptersListed[j].Number
	})

	return mangaPage
}

// Returns the pages of a chapter of a manga
func (in *InManga) GetChapter(name string, chapterNum int) (chapter models.Chapter) {
	return chapter
}
