package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// LoadEnv reads configuration from the environment (and .env if present)
// and fails fast when a required value is missing.
func LoadEnv() (*Config, error) {
	// .env is optional: in production the platform injects real env vars.
	err := godotenv.Load()
	if err != nil  {
		log.Fatal("Failed loading .env file")
	}

	cfg := &Config{
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}
	return cfg, nil
}
