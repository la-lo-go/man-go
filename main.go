package main

import (
	"MAPIes/routers"
	"MAPIes/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalln("Error loading .env")
	}

	// Create the database connection
	gorm.Init()

	// Create the router
	err := routers.CreateRouter()
	if err != nil {
		log.Fatal("Error creating the router")
	}
}
