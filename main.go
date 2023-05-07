package main

import (
	"MAPIes/routers"
	"MAPIes/gorm"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create the database connection
	gorm.Init()

	// Create the router
	err = routers.CreateRouter()
	if err != nil {
		log.Fatal("Error creating the router")
	}
}
