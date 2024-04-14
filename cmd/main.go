package main

import (
	"avito/internal/config"
	"avito/internal/db"
	"avito/internal/repository"
	"avito/internal/server/http/handler"
	"avito/internal/server/http/middleware"
	"avito/internal/server/http/router"
	"avito/internal/service"
	"avito/internal/usecase"
	"flag"
	"log"
)

var configPath = flag.String("config", "./config/config.yaml", "config path")

func main() {
	flag.Parse()

	cfg := config.MustLoad(*configPath)

	dbConn, err := db.Connect(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	bannerRepo := repository.NewBannerRepository(dbConn)
	userRepo := repository.NewUserRepository(dbConn)

	bannerCache, err := service.NewBannerCache(&cfg.Cache)
	if err != nil {
		log.Fatalf("Failed to create service: %s", err)
	}
	jwtAuth := service.NewJWTAuth(&cfg.JWT)

	authJWT := middleware.NewAuthJWT(jwtAuth)

	bannerUsecase := usecase.NewBannerUsecase(bannerRepo, bannerCache)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtAuth)

	bannerHandler := handler.NewBannerHandler(bannerUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	router := router.NewHttpRouter()
	router.Register(bannerHandler, userHandler, authJWT.JWTAuth)

	if err := router.Run(&cfg.HTTPServer); err != nil {
		log.Fatal(err)
	}
}
