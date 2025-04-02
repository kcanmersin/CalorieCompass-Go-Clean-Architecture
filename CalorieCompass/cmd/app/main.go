package main

import (
	"log"

	"CalorieCompass/internal/app"
)

func main() {
	log.Println("Starting CalorieCompass application")
	app.Run("./config/config.yml")
}