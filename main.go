package main

import (
	"log"
	"os"
	"outdoorsy/api"
	"outdoorsy/db"

	"github.com/joho/godotenv"
)

func main() {
	// Get Config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init Database
	database, err := db.NewDatabase(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Init Router
	router, err := api.NewRouter(database)
	if err != nil {
		log.Fatal(err)
	}

	// Run Router
	appPort := os.Getenv("APP_PORT")
	if err := router.Run("0.0.0.0:" + appPort); err != nil {
		log.Fatal(err)
	}
}
