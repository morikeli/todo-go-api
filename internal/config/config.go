package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port string
	SecretKey string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()

	if err != nil {
		log.Println("[WARNING]: Using environment variable but .env file not found!")
	}

	var config *Config = &Config {
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port: os.Getenv("APP_PORT"),
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	return config, nil
}