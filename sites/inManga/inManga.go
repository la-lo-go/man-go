package inManga

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"sort"
	"strconv"
	"strings"

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
		mangaID := strings.Split(mangaCoverLinkRelative, "/thumbnails/manga/")[1]

		// Get manga attributes
		mangaName := strings.Trim(s.Find(".m0").First().Text(), " ")
		mangaNameJoined, _ := utils.RemoveNonAlphanumeric(strings.Trim(mangaName, " "))
		mangaSite := in.SiteName()
		mangaLink := "https://inmanga.com/ver/manga/" + mangaID
		mangaChaptersNumber, _ := strconv.Atoi(s.Find(".icon-info text-muted").First().Text())
		mangaCover := INMANGA_THUMBNAIL_URL + mangaID

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

// GetMangaPage Returns the chapters of a manga avalible in a site
func (in *InManga) GetMangaPage(name string, url string) (mangaPage models.MangaInfo) {
	jsonResponse := InMangaMangaPage{}

	mangaID := extractInMangaID(url)

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

	mangaPage.Name = name
	mangaPage.Site = in.SiteName()
	mangaPage.Cover = INMANGA_THUMBNAIL_URL + mangaID // TODO: Ver si esta bien

	for _, chapter := range jsonResponse.Data.Result {
		// Join the manga name with the chapter number
		chapterName := "Capítulo: " + chapter.FriendlyChapterNumber

		mangaPage.ChaptersListed = append(mangaPage.ChaptersListed, models.ChapterListed{
			Number:       chapter.Number,
			Name:         chapterName,
			LinkOriginal: "https://inmanga.com/ver/manga/" + name + "/" + chapter.FriendlyChapterNumberURL + "/" + chapter.Identification,
		})
	}

	mangaPage.ChaptersNumber = len(mangaPage.ChaptersListed)

	// Sort the chapters
	sort.Slice(mangaPage.ChaptersListed, func(i, j int) bool {
		return mangaPage.ChaptersListed[i].Number < mangaPage.ChaptersListed[j].Number
	})

	err = dumpMangaToDB(mangaPage)
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

func dumpMangaToDB(page models.MangaInfo) error {
	// Search for the manga in the json if is already there
	manga, err := searchManga(page.Name)
	if err != nil {
		return err
	}

	if manga.MangaName == "" { // If the manga is not there, add it
		manga.MangaName, _ = utils.RemoveNonAlphanumeric(page.Name)
		manga.Chapters = fromChapterListedToInMangaChapter(page.ChaptersListed)

		err = gorm.AddInManga(manga)
		if err != nil {
			return err
		}
	} else { // If the manga is there, update it
		manga.Chapters = fromChapterListedToInMangaChapter(page.ChaptersListed)
		err = gorm.UpdateInManga(manga)
		if err != nil {
			return err
		}
	}

	return nil
}

// Search for the manga in inMangaChaptersJson if is already there by nameJoined
func searchManga(name string) (manga models.InMangaManga, err error) {
	nameJoined, err := utils.RemoveNonAlphanumeric(name)
	if err != nil {
		return manga, err
	}

	return gorm.SearchInManga(nameJoined)
}

func fromChapterListedToInMangaChapter(listed []models.ChapterListed) (chapters []models.InMangaChapter) {
	for _, chapter := range listed {
		chapters = append(chapters, models.InMangaChapter{
			Number: chapter.Number,
			ID:     extractInMangaID(chapter.LinkOriginal),
		})
	}

	return chapters
}

// Returns the pages of a chapter of a manga
func (in *InManga) GetChapter(name string, chapterNum float64) (chapter models.Chapter) {
	chapter.Name = name
	chapter.Site = in.SiteName()
	chapter.Number = chapterNum

	chapterDB := gorm.FindInMangaChapterID(name, chapterNum)
	if chapterDB.ID == "" {
		return chapter
	}

	doc, err := utils.GetHtmlFromGet(INMANGA_CHAPTERS_INDEX_URL + chapterDB.ID)
	if err != nil {
		fmt.Println(err)
		return chapter
	}

	log.Println(INMANGA_CHAPTERS_INDEX_URL + chapterDB.ID)

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
