package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DSN  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load config file")
	}

	port := os.Getenv("PORT")
	dsn := os.Getenv("DSN")

	if port == "" {
		log.Fatalf("please, set PORT env")
	}

	if dsn == "" {
		log.Fatalf("please, set DSN env")
	}

	return &Config{
		Port: port,
		DSN:  dsn,
	}, nil
}
