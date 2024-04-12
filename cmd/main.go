package main

import (
	"avito/internal/db"
	"avito/internal/repository"
	"avito/internal/server/http/handler"
	"avito/internal/server/http/router"
	"avito/internal/usecase"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	LoadEnvironment()

	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	bannerRepo := repository.NewBannerRepository(dbConn)

	bannerUsecase := usecase.NewBannerUsecase(bannerRepo)

	bannerHandler := handler.NewBannerHandler(bannerUsecase)

	router := router.NewHttpRouter()
	router.Register(bannerHandler)

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}

func LoadEnvironment() {
	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if err := godotenv.Load(curDir + "/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}
