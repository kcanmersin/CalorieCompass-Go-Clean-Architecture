package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	// Import required packages
	"CalorieCompass/internal/controller/html"
	v1 "CalorieCompass/internal/controller/http/v1"
	"CalorieCompass/internal/pkg/config"
	"CalorieCompass/internal/pkg/hash"
	"CalorieCompass/internal/pkg/httpserver"
	"CalorieCompass/internal/repository/fatsecret"
	"CalorieCompass/internal/repository/postgres"
	"CalorieCompass/internal/repository/token"
	"CalorieCompass/internal/service"
	"CalorieCompass/internal/usecase"
	pg "CalorieCompass/pkg/postgres"

	// Swagger docs
	_ "CalorieCompass/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(configPath string) {
	// Configuration
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Database
	postgresDB, err := pg.New(cfg.Postgres.URL, pg.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		log.Fatalf("Postgres error: %s", err)
	}
	defer postgresDB.Close()

	// FatSecret Service
	fatSecretService := service.NewFatSecretService(
		cfg.FatSecret.ClientID,
		cfg.FatSecret.ClientSecret,
		cfg.FatSecret.ConsumerKey,
		cfg.FatSecret.ConsumerSecret,
	)

	// Repositories
	userRepo := postgres.NewUserRepo(postgresDB.DB)
	jwtRepo := token.NewJWTRepo(cfg.JWT.Secret, cfg.JWT.ExpirationHour)
	foodRepo := fatsecret.NewFoodRepository(fatSecretService)

	// Hasher
	hasher := hash.NewHasher(14)

	// Use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtRepo, hasher)
	userUseCase := usecase.NewUserUseCase(userRepo)
	foodUseCase := usecase.NewFoodUseCase(foodRepo)

	// HTTP Server
	router := gin.Default()

	// Set up static files and templates
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static", "./static")

	// Swagger - make sure URL is correct and appears before routes
	url := ginSwagger.URL("/swagger/doc.json") // The URL pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// HTTP controllers
	authController := v1.NewAuthController(authUseCase)
	userController := v1.NewUserController(userUseCase)
	foodController := v1.NewFoodController(foodUseCase)
	v1.NewRouter(router, authController, userController, foodController, jwtRepo)

	// HTML controllers
	htmlAuthController := html.NewAuthController(authUseCase)
	htmlHomeController := html.NewHomeController(authUseCase, jwtRepo)

	// HTML routes
	router.GET("/", htmlHomeController.Home)
	router.GET("/login", htmlAuthController.ShowLoginPage)
	router.POST("/login", htmlAuthController.ProcessLogin)
	router.GET("/register", htmlAuthController.ShowRegisterPage)
	router.POST("/register", htmlAuthController.ProcessRegister)
	router.GET("/logout", htmlAuthController.Logout)

	// HTTP Server
	httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("signal received: %s", s.String())
	case err := <-httpServer.Notify():
		log.Printf("server error: %s", err)
	}

	// Shutdown
	if err := httpServer.Shutdown(); err != nil {
		log.Printf("server shutdown error: %s", err)
	}
}
