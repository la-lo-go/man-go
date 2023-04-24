package routers

import (
	"MAPIes/endpoints"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func CreateRouter() error {
	// ROUTER SET UP
	router := gin.Default()
	// err := router.SetTrustedProxies([]string{"192.168.1.2"})
	// if err != nil {
	// 	return err
	// }

	// ENDPOINTS
	addRoutes(router)

	IP := os.Getenv("IP")
	PORT := os.Getenv("PORT")

	// Run the server
	addres := IP + ":" + PORT
	err := router.Run(addres)
	if err != nil {
		log.Fatal("Error running the API")
		return err
	}

	return nil
}

// Crete the routes
func addRoutes(router *gin.Engine) {
	router.GET("/", endpoints.Ping)
	router.GET("/busqueda", endpoints.Search)
	router.GET("/manga/:mangaName/:site", endpoints.MangaPage)
}
