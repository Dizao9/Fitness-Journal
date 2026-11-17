package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DSN              string
	AccessJWTSecret  string
	RefreshJWTSecret string
	CookieSecure     bool
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("config file is not loaded")
	}

	port := os.Getenv("PORT")
	dsn := os.Getenv("DSN")
	jwtSecret := os.Getenv("JWT_SECRET")
	refreshJWTSecret := os.Getenv("REFRESH_JWT_SECRET")
	cookieSecureRaw := os.Getenv("COOKIE_SECURE")

	cookieSecure, err := strconv.ParseBool(cookieSecureRaw)
	if err != nil {
		cookieSecure = true
		log.Printf("COOKIE_SECURE is failed. set default value = TRUE")
	}

	var missing []string

	if port == "" {
		missing = append(missing, "PORT")
	}

	if dsn == "" {
		missing = append(missing, "DSN")
	}

	if jwtSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}

	if refreshJWTSecret == "" {
		missing = append(missing, "REFRESH_JWT_SECRET")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return &Config{
		Port:             port,
		DSN:              dsn,
		AccessJWTSecret:  jwtSecret,
		RefreshJWTSecret: refreshJWTSecret,
		CookieSecure:     cookieSecure,
	}, nil
}
