package endpoints

import (
    "log"
	"net/http"
	"strconv"
	"strings"
    
	"github.com/gin-gonic/gin"
    
	"MAPIes/sites"
	"MAPIes/gorm"
	"MAPIes/models"
	"MAPIes/utils"
)

var maxCoincidencesPerSite int

const MAX_COINCIDENCE_DEFAULT = 10
const SEARCH_EXPIRE_HOURS = 24

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

	response, searchInDBResult := gorm.FindSearch(searchStrFormatted)

	// If the search in the database is not found or is partial
	switch searchInDBResult {
	case "Partial":
		response = filterPartial(response, searchStrFormatted)
	case "NotFound":
		response = searchMangas(searchStr)
		gorm.DumpSearchToDB(searchStr, response)
	}

	if len(response) > 0 {
		// trim to return only de Max amount of coincidences
		context.IndentedJSON(http.StatusOK, trimMangasToMaxPerSize(response))
	} else { // If no mangas are found
		context.IndentedJSON(http.StatusNotFound, []models.Manga{})
	}
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
