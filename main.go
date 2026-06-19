package main

import (
	"log"
	"os"

	"sport-venue-rental-api/config"
	"sport-venue-rental-api/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[WARN] .env file not found, using system environment variables")
	}

	config.ConnectDB()

	e := echo.New()

	routes.InitRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("[INFO] server running on port " + port)

	if err := e.Start(":" + port); err != nil {
		log.Fatal("[ERROR] failed to start server: ", err)
	}
}
