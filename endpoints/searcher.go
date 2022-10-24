package endpoints

import (
	"MAPIes/sites"
	"MAPIes/sites/inManga"
	"MAPIes/sites/mangaMx"
	"MAPIes/sites/nyaa"
	"MAPIes/sites/tuMangaNet"
	"encoding/json"
	"fmt"
	"net/http"

	// "reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"MAPIes/general_functions"
	"MAPIes/models"
)

var maxCoincidencesPerSite int

const MAX_COINCIDENCE_DEFAULT = 1
const SEARCH_EXPIRE_HOURS = 24

var jsonSearch = models.NewSearchCacheJson()
var jsonLinks = models.NewLinksMangasCacheJson()

func Search(context *gin.Context) {
	queryValues := context.Request.URL.Query()

	searchStr := strings.ToLower(queryValues["search"][0])

	// check if queryValues["max"] has a value
	if len(queryValues["max"]) > 0 {
		maxCoincidencesPerSite, _ = strconv.Atoi(queryValues["max"][0])
	} else {
		maxCoincidencesPerSite = MAX_COINCIDENCE_DEFAULT
	}

	jsonReturn, searchInJson := findSearchInJson(searchStr)

	switch searchInJson {
	case "NotFound":
		jsonReturn = searchMangas(searchStr)
		dumpSearchToJson(searchStr, jsonReturn)
	case "Partial":
		jsonReturn = filterPartial(jsonReturn, searchStr)
	}

	if len(jsonReturn) > 0 {
		dumpLinksToJson(jsonReturn)
		// trim to return only de Max amount of coincidences
		context.IndentedJSON(http.StatusOK, trimMangas(jsonReturn))
	} else {
		context.IndentedJSON(http.StatusNotFound, []models.Manga{})
	}
}

func findSearchInJson(searchStr string) ([]models.Manga, string) {
	var jsonSlice []models.ApiSearch

	// Format the search to match the json format
	searchStr, _ = general_functions.RemoveNonAlphanumeric(searchStr)

	jsonFile, err := jsonSearch.Read()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(jsonFile, &jsonSlice)
	if err != nil {
		fmt.Println(err)
	}

	if len(jsonSlice) > 0 {
		for _, s := range jsonSlice {
			if strings.Contains(searchStr, s.Search) && time.Since(s.Date) < time.Duration(SEARCH_EXPIRE_HOURS)*time.Hour { // Check if the searchStr is a substr of a cached search
				if searchStr == s.Search { // Check if the searchStr is the same as a cached search
					fmt.Println("\n\n>>>> [json/searchBuffer.json]: Exact match found")
					return s.Response, "Exact"
				} else {
					fmt.Println("\n\n>>>> [json/searchBuffer.json]: Partial match found")
					return s.Response, "Partial"
				}
			}
		}
	}

	return nil, "NotFound" // If the searchStr is not found in the json file
}

func searchMangas(searchStr string) []models.Manga {

	// Classes of different sites to search on
	sitesClasses := []sites.IntSite{
		&inManga.InManga{},
		&nyaa.Nyaa{},
		&tuMangaNet.TuMangaNet{},
		&mangaMx.MangaMX{},
	}

	var searchedMangas []models.Manga

	fmt.Printf("\n\nNueva busqueda: " + searchStr)

	// iterate through sitesClasses
	for _, s := range sitesClasses {
		siteSearchMangas := searchBySite(s, searchStr, searchedMangas)
		searchedMangas = append(searchedMangas, siteSearchMangas...)
	}

	// Clean the empty mangas
	return clearMangas(searchedMangas)
}

func searchBySite(s sites.IntSite, searchStr string, searchedMangas []models.Manga) []models.Manga {
	actualSiteMangas, _ := s.GetMangas(searchStr, searchedMangas)

	// Print the site and the found mangas
	// fmt.Printf("\n>>>> [%s RETURN]: %#v\n\n", reflect.TypeOf(s), actualSiteMangas)

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

func trimMangas(searchedMangas []models.Manga) (trimmedMangas []models.Manga) {

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
		if strings.Contains(strings.ToLower(m.Name), searchStr) {
			exitSlice = append(exitSlice, m)
		}
	}

	// Trim results and return them
	return trimMangas(exitSlice)
}

func dumpSearchToJson(searchStr string, response []models.Manga) {
	var jsonSlice []models.ApiSearch

	jsonFile, err := jsonSearch.Read()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(jsonFile, &jsonSlice)
	if err != nil {
		fmt.Println("error:", err)
	}

	jsonSlice = append(jsonSlice, models.ApiSearch{Search: searchStr, Date: time.Now(), Response: response})

	err = jsonSearch.Write(jsonSlice)
	if err != nil {
		return
	}
}

func dumpLinksToJson(response []models.Manga) {
	var jsonLinksCopy map[string]json.RawMessage
	var mangaLinks models.MangaLinksClustered
	var siteLink models.SiteLink

	jsonFile, err := jsonLinks.Read()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(jsonFile, &jsonLinksCopy)
	if err != nil {
		fmt.Println(err)
	}

	// Initialize the jsonLinksCopy map if is empty
	if jsonLinksCopy == nil {
		jsonLinksCopy = make(map[string]json.RawMessage)
	}

	for _, m := range response {
		mangaName, _ := general_functions.RemoveNonAlphanumeric(m.Name)
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
				mangaLinks.SitesLinks = append(mangaLinks.SitesLinks, models.SiteLink{Site: m.Site, Link: m.Link})

				jsonLinksCopy[mangaName], err = json.Marshal(mangaLinks)
				if err != nil {
					fmt.Println(err)
				}
			}

		} else { // There is no coincidence, create a new entry
			siteLink = models.SiteLink{Site: m.Site, Link: m.Link}
			mangaLinks = models.MangaLinksClustered{Name: m.Name, SitesLinks: []models.SiteLink{siteLink}}

			jsonLinksCopy[mangaName], err = json.Marshal(mangaLinks)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	err = jsonLinks.Write(jsonLinksCopy)
	if err != nil {
		return
	}
}
