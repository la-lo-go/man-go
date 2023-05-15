package endpoints

import (
	"MAPIes/sites"
	"MAPIes/sites/inManga"
	"MAPIes/sites/mangaOni"
	"MAPIes/sites/nyaa"
	"MAPIes/sites/tuMangaNet"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"MAPIes/models"
	"fmt"
	"net/http"
)

// Variable of the interface 'IntSite' that is implemented by the different sites structs
var siteObj sites.IntSite

func MangaPage(context *gin.Context) {

	mangaName := strings.ToLower(context.Param("mangaName"))
	site := strings.ToLower(context.Param("site"))

	switch site {
	case "inmanga":
		siteObj = &inManga.InManga{}

	case "mangamx":
		siteObj = &mangaOni.MangaOni{}

	case "nyaa":
		siteObj = &nyaa.Nyaa{}

	case "tumanganet":
		siteObj = &tuMangaNet.TuMangaNet{}

	default:
		siteObj = nil
	}

	fmt.Println(reflect.TypeOf(siteObj))

	if siteObj != nil {
		// Search the url from the site in JSON
		url := searchMangaURL(mangaName, site)

		if url != "" {
			context.IndentedJSON(http.StatusOK, siteObj.GetMangaPage(mangaName, url))
		} else {
			context.IndentedJSON(http.StatusNotFound, []models.Manga{})
		}

	} else {
		context.IndentedJSON(http.StatusNotFound, []models.MangaInfo{})
	}
}

// search the url of the manga in the API_linksMangas.json
// file based on the name and the site and return the url
func searchMangaURL(mangaName string, siteName string) string {
	var jsonLinksCopy map[string]json.RawMessage
	var mangaLinks models.MangaLinksClustered

	jsonFile, err := jsonLinks.Read()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(jsonFile, &jsonLinksCopy)
	if err != nil {
		fmt.Println(err)
	}

	if jsonLinksCopy == nil {
		jsonLinksCopy = make(map[string]json.RawMessage)
	}

	// Find if there is a coincidence with the same name
	err = json.Unmarshal(jsonLinksCopy[mangaName], &mangaLinks)
	if err == nil { // There is a coincidence, get the link from the site
		for _, sl := range mangaLinks.SitesLinks {
			if strings.ToLower(sl.Site) == siteName {
				return sl.Link
			}
		}
	}

	// If there is no coincidence, return an empty string
	return ""
}
