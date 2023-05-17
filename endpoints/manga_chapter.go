package endpoints

import (
	"MAPIes/models"
	"MAPIes/sites/inManga"
	"MAPIes/sites/mangaOni"
	"MAPIes/sites/nyaa"
	"MAPIes/sites/tuMangaNet"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func MangaChapter(context *gin.Context) {

	site := strings.ToLower(context.Param("site"))
	mangaName := strings.ToLower(context.Param("mangaName"))
	chapterNumber := context.Param("chapterNumber")

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

	numberFloat, err := strconv.ParseFloat(chapterNumber, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, []models.MangaInfo{})
		return
	}

	if siteObj != nil {
		context.IndentedJSON(http.StatusOK, siteObj.GetChapter(mangaName, numberFloat))
	} else {
		context.IndentedJSON(http.StatusNotFound, []models.MangaInfo{})
	}
}
