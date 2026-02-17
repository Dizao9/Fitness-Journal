package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DSN       string
	JWTSecret string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load config file")
	}

	port := os.Getenv("PORT")
	dsn := os.Getenv("DSN")
	jwtSecret := os.Getenv("JWT_SECRET")

	if port == "" {
		log.Fatalf("please, set PORT env")
	}

	if dsn == "" {
		log.Fatalf("please, set DSN env")
	}

	if jwtSecret == "" {
		log.Fatalf("please, set JWT_SECRET env")
	}

	return &Config{
		Port:      port,
		DSN:       dsn,
		JWTSecret: jwtSecret,
	}, nil
}
