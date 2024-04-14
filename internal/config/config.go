package config

import (
	"avito/internal/db"
	"avito/internal/server/http/router"
	"avito/internal/service"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer router.ConfigHTTPServer `yaml:"http_server"`
	Database   db.ConfigDatabase       `yaml:"database"`
	Cache      service.ConfigCache     `yaml:"cache"`
}

func MustLoad(curDir string) *Config {
	configPath := curDir + os.Getenv("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}
