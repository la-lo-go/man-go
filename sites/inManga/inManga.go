package inManga

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"MAPIes/gorm"
	"MAPIes/models"
	"MAPIes/utils"
)

const INMANGA_THUMBNAIL_URL = "https://pack-yak.intomanga.com/thumbnails/manga/"
const INMANGA_GET_ALL_URL = "https://inmanga.com/chapter/getall?mangaIdentification="
const INMANGA_CHAPTERS_INDEX_URL = "https://inmanga.com/chapter/chapterIndexControls?identification="
const INMANGA_PAGE_URL = "https://pack-yak.intomanga.com/images/manga/Name/chapter/NumberChapter/page/NumberPage/"

type InManga struct{}

func (in *InManga) SiteName() string {
	return "inmanga"
}

// GetMangas Returns the mangas of a site that match the search
func (in *InManga) GetMangas(searchValue string, searchedMangas []models.Manga) (mangas []models.Manga, err error) {
	searchStringFormated := strings.Replace(searchValue, " ", "+", -1)
	url := "https://inmanga.com/manga/getMangasConsultResult"
	payload := strings.NewReader(`filter%5Bgeneres%5D%5B%5D=-1&filter%5BqueryString%5D=+` + searchStringFormated + `+&filter%5Bskip%5D=0&filter%5Btake%5D=10&filter%5Bsortby%5D=1&filter%5BbroadcastStatus%5D=0&filter%5BonlyFavorites%5D=false&d=`)

	// Consigue el HTML de la pagina
	doc, err := utils.GetHtmlFromPost(url, payload)
	if err != nil {
		fmt.Println(err)
		return []models.Manga{}, err
	}

	doc.Find(".manga-result").Each(func(i int, s *goquery.Selection) {
		// Get the manga name with ID
		mangaCoverLinkRelative, _ := s.Find("img").Attr("data-src")
		mangaNameAndMangaID := strings.Split(mangaCoverLinkRelative, "/thumbnails/manga/")[1]
		mangaID := strings.Split(mangaNameAndMangaID, "/")[1]

		// Get manga attributes
		mangaName := strings.Trim(s.Find(".m0").First().Text(), " ")
		mangaNameJoined, _ := utils.RemoveNonAlphanumeric(strings.Trim(mangaName, " "))
		mangaSite := in.SiteName()
		mangaLink := "https://inmanga.com/ver/manga/" + mangaNameAndMangaID
		mangaChaptersNumber, _ := strconv.Atoi(s.Find(".icon-info text-muted").First().Text())
		mangaCover := INMANGA_THUMBNAIL_URL + mangaNameAndMangaID

		mangas = append(mangas, models.Manga{
			Name:           mangaName,
			NameJoined:     mangaNameJoined,
			Site:           mangaSite,
			Link:           mangaLink,
			ChaptersNumber: mangaChaptersNumber,
			Cover:          mangaCover,
			WebID:          mangaID,
		})
	})

	return mangas, nil
}

// GetMangaPage Returns the chapters of a manga avalible in a site
func (in *InManga) GetMangaPage(name string, url string) (mangaPage models.Manga) {
	jsonResponse := InMangaMangaPage{}
	siteName := in.SiteName()

	mangaDBInfo, err := gorm.SearchManga(siteName, name)
	if err == nil {
		return mangaDBInfo
	}

	mangaID := mangaDBInfo.WebID

	urlRequest := INMANGA_GET_ALL_URL + mangaID
	response, err := utils.GetJsonFromGet(urlRequest)

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

	mangaPage.Name = mangaDBInfo.Name
	mangaPage.NameJoined = name
	mangaPage.Site = siteName
	mangaPage.Cover = mangaDBInfo.Cover
	mangaPage.WebID = mangaID
	mangaPage.Link = url

	for _, chapter := range jsonResponse.Data.Result {
		// Join the manga name with the chapter number
		chapterName := "Capítulo: " + chapter.FriendlyChapterNumber

		mangaPage.Chapters = append(mangaPage.Chapters, models.Chapter{
			Number:       chapter.Number,
			Name:         chapterName,
			Site:         siteName,
			LinkOriginal: "https://inmanga.com/ver/manga/" + name + "/" + chapter.FriendlyChapterNumberURL + "/" + chapter.Identification,
			WebID:        chapter.Identification,
		})
	}

	mangaPage.ChaptersNumber = len(mangaPage.Chapters)

	// Sort the chapters
	sort.Slice(mangaPage.Chapters, func(i, j int) bool {
		return mangaPage.Chapters[i].Number < mangaPage.Chapters[j].Number
	})

	err = gorm.DumpMangaToDB(in.SiteName(), mangaPage)
	if err != nil {
		return mangaPage
	}

	return mangaPage
}

func extractInMangaID(url string) (ID string) {
	urlSplit := strings.Split(url, "/")
	ID = urlSplit[len(urlSplit)-1]

	return ID
}

// Returns the pages of a chapter of a manga
func (in *InManga) GetChapter(name string, chapterNum float64) (chapter models.Chapter) {
	chapter.Name = name
	chapter.Site = in.SiteName()
	chapter.Number = chapterNum

	chapterDB := gorm.FindChapterWebID(in.SiteName(), name, chapterNum)
	if chapterDB.WebID == "" {
		return chapter
	}

	doc, err := utils.GetHtmlFromGet(INMANGA_CHAPTERS_INDEX_URL + chapterDB.WebID)
	if err != nil {
		fmt.Println(err)
		return chapter
	}

	doc.Find(".PageListClass").Each(func(i int, s *goquery.Selection) {
		s.Find("option").Each(func(i int, s *goquery.Selection) { // get the value of each option
			pageID, _ := s.Attr("value")
			chapter.Pages = append(chapter.Pages, models.Page{
				Number: i + 1,
				Link:   INMANGA_PAGE_URL + pageID,
			})
		})
	})

	return chapter
}
