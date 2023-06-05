package nyaa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	models "MAPIes/models"
	"MAPIes/utils"
)

var myClient = &http.Client{Timeout: 5 * time.Second}

var nyaaSearchJson = models.NewNyaaSearchJson()

const NYAA_DOMAIN = "https://manganyaa.com/"
const NYAA_GIT_RESPOSITORY = "https://raw.githubusercontent.com/saulabagnale/asdf-ma-jsons/master/"
const NYAA_GIT_MANGAS_NAMES_JSON = "mangaNames.json"
const NYAA_GIT_MANGA_JSON = "es.json"
const NYAA_CAP_URI = "/leer-online-gratis-espanol/capitulo/"
const NYAA_IMAGE_START_URL = "https://content.manganyaa.com/file/mnyaaa/"

type Nyaa struct{}

// Site name returns the name of the site IN LOWERCASE
func (n *Nyaa) SiteName() string {
	return "nyaa"
}

// GetMangas returns the mangas of a site that match the search
func (n *Nyaa) GetMangas(searchValue string, searchedMangas []models.Manga) ([]models.Manga, error) {
	url := NYAA_GIT_RESPOSITORY + NYAA_GIT_MANGAS_NAMES_JSON

	listaMangas, err := getMangasNamesJson(url)
	if err != nil {
		fmt.Println(err)
		return []models.Manga{}, err
	}

	filtered := filterMangas(searchValue, listaMangas)

	if len(filtered) > 0 {
		mangasReturn := []models.Manga{}
		manga := new(models.Manga)

		for _, m := range filtered {
			manga.Name = m.Name
			manga.NameJoined, _ = utils.RemoveNonAlphanumeric(m.Name)
			manga.Site = n.SiteName()
			manga.Link = "https://manganyaa.com/" + m.JoinedName + "/leer-online-gratis-espanol"
			mangaChaptersNumber := "99" //TODO: Get the number of chapters of the manga
			manga.ChaptersNumber, _ = strconv.Atoi(mangaChaptersNumber)
			manga.Cover = NYAA_IMAGE_START_URL + m.JoinedName + `/description/1.jpg`

			mangasReturn = append(mangasReturn, *manga)
		}

		return mangasReturn, nil
	} else {
		return []models.Manga{}, nil
	}
}

func getMangasNamesJson(url string) (response []NyaaSearch, err error) {
	// Check if the json file is created or is updated in the last 12 hours
	checkJson, _ := nyaaSearchJson.Check()

	if checkJson { // The json file is up to date
		byteValue, _ := nyaaSearchJson.Read()
		if err != nil {
			fmt.Println(err)
			return []NyaaSearch{}, err
		}

		err = json.Unmarshal(byteValue, &response)

		if err != nil {
			fmt.Println(err)
			return []NyaaSearch{}, err
		}

		return response, nil

	} else { // The json file is not up to date or does not exist
		r, err := myClient.Get(url)
		if err != nil {
			return response, err
		}
		defer r.Body.Close()

		response, err = formatAndConstructResponseSlice(r)
		if err != nil {
			fmt.Println(err)
			return []NyaaSearch{}, err
		}

		nyaaSearchJson.Write(response)
	}

	return response, nil
}

func formatAndConstructResponseSlice(r *http.Response) (response []NyaaSearch, err error) {
	var mangaActualSlice []string
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		return response, err
	}

	// Delete de [[, ]] from the response
	bodyStr := strings.Replace(string(body), "[[", "", -1)
	bodyStr = strings.Replace(bodyStr, "]]", "", -1)

	// Split the response into slices
	bodySplit := strings.Split(bodyStr, "],[")

	// DEBUG: Print the response
	// fmt.Printf("%v", bodySplit)

	for _, m := range bodySplit {
		// Deletes the " and then splits the response with the , (["ab-cd","ab cd"] -> [[ab-cd][ab cd]])
		mangaActualSlice = strings.Split(strings.ToLower(strings.ReplaceAll(m, `"`, ``)), ",")
		response = append(response, NyaaSearch{JoinedName: mangaActualSlice[0], Name: mangaActualSlice[1]})
	}

	return response, nil
}

func filterMangas(searchStr string, listaMangas []NyaaSearch) []NyaaSearch {
	var filteredMangas []NyaaSearch
	searchStr = strings.ToLower(searchStr)

	for _, m := range listaMangas {
		if strings.Contains(m.Name, searchStr) {
			filteredMangas = append(filteredMangas, m)
		}
	}

	return filteredMangas
}

// GetMangaPage returns the chapters of a manga avalible in a site
func (n *Nyaa) GetMangaPage(name string, url string) (mangaPage models.Manga) {
	var numberParsed string
	siteName := n.SiteName()
	jsonResponse := NyaaMangaPage{}
	urlRequest := NYAA_GIT_RESPOSITORY + "series/" + name + "/" + NYAA_GIT_MANGA_JSON

	response, err := utils.GetJsonFromGet(urlRequest)

	if err != nil {
		fmt.Println(err)
		return mangaPage
	}

	err = json.Unmarshal([]byte(response), &jsonResponse)
	if err != nil {
		fmt.Println(err, urlRequest)
		return mangaPage
	}

	mangaPage.Name = jsonResponse.MangaName
	mangaPage.NameJoined, _ = utils.RemoveNonAlphanumeric(jsonResponse.MangaName)
	mangaPage.Site = siteName
	mangaPage.Link = url
	mangaPage.Cover = "https://content.manganyaa.com/file/mnyaaa/" + mangaPage.NameJoined + "/description/1.jpg"

	// append to mangaPage.ChaptersListed only the jsonResponse.Chapters that are not 0
	for _, c := range jsonResponse.Chs {
		if c.Pages != 0 {
			numberParsed = fmt.Sprint(c.OrderNumber)

			if err != nil {
				fmt.Println(err)
				return mangaPage
			}

			mangaPage.Chapters = append(mangaPage.Chapters, models.Chapter{
				Name:         mangaPage.Name,
				Site:         siteName,
				Number:       float64(c.OrderNumber),
				LinkOriginal: NYAA_DOMAIN + name + NYAA_CAP_URI + numberParsed,
			})
		}
	}

	mangaPage.ChaptersNumber = len(mangaPage.Chapters)

	return mangaPage
}

// Returns the pages of a chapter of a manga
func (n *Nyaa) GetChapter(name string, chapterNum float64) (chapter models.Chapter) {
	siteName := n.SiteName()
	jsonResponse := NyaaMangaPage{}
	urlRequest := NYAA_GIT_RESPOSITORY + "series/" + name + "/" + NYAA_GIT_MANGA_JSON

	response, err := utils.GetJsonFromGet(urlRequest)

	if err != nil {
		fmt.Println(err)
		return chapter
	}

	err = json.Unmarshal([]byte(response), &jsonResponse)
	if err != nil {
		fmt.Println(err, urlRequest)
		return chapter
	}

	// Find the chapter in the jsonResponse.ch
	for _, c := range jsonResponse.Chs {
		if c.OrderNumber == chapterNum {
			chapter.Name = name
			chapter.Site = siteName
			chapter.Number = float64(c.OrderNumber)
			chapter.LinkOriginal = NYAA_DOMAIN + name + NYAA_CAP_URI + fmt.Sprint(c.OrderNumber)

			for i := 1; i <= c.Pages; i++ {
				chapter.Pages = append(chapter.Pages, models.Page{
					Number: i,
					Link:   NYAA_IMAGE_START_URL + name + "/es/" + fmt.Sprint(c.OrderNumber) + "/" + fmt.Sprint(i) + ".jpg",
				})
			}
		}
	}

	return chapter
}
