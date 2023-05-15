package endpoints

import (
	"MAPIes/sites"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"MAPIes/models"
	"MAPIes/utils"
	// "MAPIes/gorm"
)

var maxCoincidencesPerSite int

const MAX_COINCIDENCE_DEFAULT = 10
const SEARCH_EXPIRE_HOURS = 24

var jsonSearch = models.NewSearchCacheJson()
var jsonLinks = models.NewLinksMangasCacheJson()

func Search(context *gin.Context) {
	queryValues := context.Request.URL.Query()

	searchStr := strings.ToLower(queryValues["search"][0])
	searchStrFormatted, _ := utils.RemoveNonAlphanumeric(searchStr)

	// check if queryValues["max"] has a value
	if len(queryValues["max"]) > 0 && queryValues["max"][0] != "" && queryValues["max"][0] != "0" {
		requestedMaxCoincidencesPerSite, _ := strconv.Atoi(queryValues["max"][0])

		if requestedMaxCoincidencesPerSite < MAX_COINCIDENCE_DEFAULT {
			maxCoincidencesPerSite, _ = strconv.Atoi(queryValues["max"][0])
		} else {
			maxCoincidencesPerSite = MAX_COINCIDENCE_DEFAULT
		}
	} else {
		maxCoincidencesPerSite = MAX_COINCIDENCE_DEFAULT
	}

	response, searchInJsonResult := findSearchInJson(searchStrFormatted)

	// If the search in the database is not found or is partial
	switch searchInJsonResult {
	case "Partial":
		response = filterPartial(response, searchStrFormatted)
	case "NotFound":
		response = searchMangas(searchStr)
		dumpSearchToDB(searchStr, response)
	}

	if len(response) > 0 {
		dumpLinksToJson(response)

		// trim to return only de Max amount of coincidences
		context.IndentedJSON(http.StatusOK, trimMangasToMaxPerSize(response))
	} else { // If no mangas are found
		context.IndentedJSON(http.StatusNotFound, []models.Manga{})
	}
}

// Find the search in the database and return a status string.
// Possible status: "Exact", "Partial", "NotFound"
func findSearchInJson(searchStr string) ([]models.Manga, string) {
	var jsonSlice []models.Search

	// Format the search to match the json format
	searchStr, _ = utils.RemoveNonAlphanumeric(searchStr)

	jsonFile, err := jsonSearch.Read()
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(jsonFile, &jsonSlice)
	if err != nil {
		log.Println(err)
	}

	if len(jsonSlice) > 0 {
		for _, s := range jsonSlice {
			if strings.Contains(searchStr, s.Search) { // Check if the searchStr is a substr of a cached search
				if searchStr == s.Search { // Check if the searchStr is the same as a cached search
					log.Println("\n>>>> [json/searchBuffer.json]: Exact match found")
					return s.Mangas, "Exact"
				} else {
					log.Println("\n>>>> [json/searchBuffer.json]: Partial match found")
					return s.Mangas, "Partial"
				}
			}
		}
	}

	// If the searchStr is not found in the json file
	return nil, "NotFound"
}

func searchMangas(searchStr string) []models.Manga {
	var searchedMangas []models.Manga

	log.Println("New search: " + searchStr)

	// iterate through sitesTypes
	for _, s := range sites.SitesTypes {
		siteSearchMangas := searchBySite(s, searchStr, searchedMangas)
		searchedMangas = append(searchedMangas, siteSearchMangas...)
	}

	// Clean the empty mangas
	return clearMangas(searchedMangas)
}

func searchBySite(s sites.IntSite, searchStr string, searchedMangas []models.Manga) []models.Manga {
	actualSiteMangas, _ := s.GetMangas(searchStr, searchedMangas)

	// Print the site and the found mangas
	// log.Println(">>>> [%s RETURN]: %#v\n\n", reflect.TypeOf(s), actualSiteMangas)

	return actualSiteMangas
}

func clearMangas(searchedMangas []models.Manga) (clearedMangas []models.Manga) {

	for _, m := range searchedMangas {
		if m.Name != "" {
			clearedMangas = append(clearedMangas, m)
		}
	}

	return clearedMangas
}

// Trim the mangas to return only the max amount of coincidences PER SITE
func trimMangasToMaxPerSize(searchedMangas []models.Manga) (trimmedMangas []models.Manga) {

	// count the times that the same Manga.Site is found
	mangasCount := map[string]int{}
	for _, m := range searchedMangas {
		mangasCount[m.Site]++

		if mangasCount[m.Site] <= maxCoincidencesPerSite {
			trimmedMangas = append(trimmedMangas, m)
		}
	}

	return trimmedMangas
}

func filterPartial(enterSlice []models.Manga, searchStr string) (exitSlice []models.Manga) {
	// Filters the mangas that contains the search string
	for _, m := range enterSlice {
		formatedName, _ := utils.RemoveNonAlphanumeric(m.Name)

		if strings.Contains(formatedName, searchStr) {
			exitSlice = append(exitSlice, m)
		}
	}

	// Trim results and return them
	return exitSlice
}

func dumpSearchToDB(searchStr string, mangas []models.Manga) {
	var jsonSlice []models.Search
	searchStrFormatted, _ := utils.RemoveNonAlphanumeric(searchStr)

	jsonFile, err := jsonSearch.Read()
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(jsonFile, &jsonSlice)
	if err != nil {
		log.Println("error:", err)
	}

	jsonSlice = append(
		jsonSlice,
		models.Search{
			Search: searchStrFormatted,
			Mangas: mangas,
		},
	)

	err = jsonSearch.Write(jsonSlice)
	if err != nil {
		return
	}

	// err = gorm.UploadSearch(jsonSlice[0])
	// if err != nil {
	// 	return
	// }
}

func dumpLinksToJson(response []models.Manga) {
	var jsonLinksCopy map[string]json.RawMessage
	var mangaLinks models.MangaLinksClustered
	var siteLink models.SiteLink

	jsonFile, err := jsonLinks.Read()
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(jsonFile, &jsonLinksCopy)
	if err != nil {
		log.Println(err)
	}

	// Initialize the jsonLinksCopy map if is empty
	if jsonLinksCopy == nil {
		jsonLinksCopy = make(map[string]json.RawMessage)
	}

	for _, m := range response {
		mangaName, _ := utils.RemoveNonAlphanumeric(m.Name)
		mangaName = strings.ToLower(mangaName)

		// Find if there is a coincided with the same name
		err = json.Unmarshal(jsonLinksCopy[mangaName], &mangaLinks)
		if err == nil { // There is a coincidence, add the new link if it does not already exist
			// check if the site is already in the list
			siteFound := false
			for _, sl := range mangaLinks.SitesLinks {
				if sl.Site == m.Site {
					siteFound = true
					break
				}
			}

			if !siteFound { // If the site is not in the list of the manga, add it
				mangaLinks.SitesLinks = append(
					mangaLinks.SitesLinks,
					models.SiteLink{
						Site: m.Site,
						Link: m.Link,
					},
				)

				jsonLinksCopy[mangaName], err = json.Marshal(mangaLinks)
				if err != nil {
					log.Println(err)
				}
			}

		} else { // There is no coincidence, create a new entry
			siteLink = models.SiteLink{
				Site: m.Site,
				Link: m.Link,
			}
			mangaLinks = models.MangaLinksClustered{
				Name:       m.Name,
				SitesLinks: []models.SiteLink{siteLink},
			}

			jsonLinksCopy[mangaName], err = json.Marshal(mangaLinks)
			if err != nil {
				log.Println(err)
			}
		}
	}

	err = jsonLinks.Write(jsonLinksCopy)
	if err != nil {
		return
	}
}
