package main

import (
	"avito/internal/config"
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
	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	LoadEnvironment(curDir)

	cfg := config.MustLoad(curDir)

	dbConn, err := db.Connect(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	bannerRepo := repository.NewBannerRepository(dbConn)
	userRepo := repository.NewUserRepository(dbConn)

	bannerUsecase := usecase.NewBannerUsecase(bannerRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	bannerHandler := handler.NewBannerHandler(bannerUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	router := router.NewHttpRouter()
	router.Register(bannerHandler, userHandler)

	if err := router.Run(&cfg.HTTPServer); err != nil {
		log.Fatal(err)
	}
}

func LoadEnvironment(curDir string) {
	if err := godotenv.Load(curDir + "/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}
