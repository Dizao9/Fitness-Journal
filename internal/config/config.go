package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DSN              string
	JWTSecret        string
	RefreshJWTSecret string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load config file")
	}

	port := os.Getenv("PORT")
	dsn := os.Getenv("DSN")
	jwtSecret := os.Getenv("JWT_SECRET")
	refreshJWTSecret := os.Getenv("REFRESH_JWT_SECRET")

	if port == "" {
		log.Fatalf("please, set PORT env")
	}

	if dsn == "" {
		log.Fatalf("please, set DSN env")
	}

	if jwtSecret == "" {
		log.Fatalf("please, set JWT_SECRET env")
	}

	if refreshJWTSecret == "" {
		log.Fatalf("please, set REFRESH_JWT_SECRET")
	}

	return &Config{
		Port:             port,
		DSN:              dsn,
		JWTSecret:        jwtSecret,
		RefreshJWTSecret: refreshJWTSecret,
	}, nil
}
