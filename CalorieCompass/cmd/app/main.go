package main

import (
	"log"

	"CalorieCompass/internal/app"
)

// @title CalorieCompass API
// @version 1.0
// @description A calorie tracking API with authentication
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	log.Println("Starting CalorieCompass application")
	app.Run("./config/config.yml")
}
