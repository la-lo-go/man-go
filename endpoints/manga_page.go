package endpoints

import (
	"MAPIes/gorm"
	"MAPIes/sites"
	"MAPIes/sites/inManga"
	"MAPIes/sites/mangaOni"
	"MAPIes/sites/nyaa"
	"MAPIes/sites/tuMangaNet"
	"strings"

	"github.com/gin-gonic/gin"

	"MAPIes/models"
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

	case "mangaOni":
		siteObj = &mangaOni.MangaOni{}

	case "nyaa":
		siteObj = &nyaa.Nyaa{}

	case "tumanganet":
		siteObj = &tuMangaNet.TuMangaNet{}

	default:
		siteObj = nil
	}

	if siteObj != nil {
		// Search the url from the site in JSON
		url := gorm.SearchMangaURL(mangaName, site)

		if url != "" {
			context.IndentedJSON(http.StatusOK, siteObj.GetMangaPage(mangaName, url))
		} else {
			context.IndentedJSON(http.StatusNotFound, []models.Manga{})
		}

	} else {
		context.IndentedJSON(http.StatusNotFound, []models.MangaInfo{})
	}
}
