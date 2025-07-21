package config

import (
	"os"
	"time"
)

type Config struct {
	HTTP struct {
		Host string
		Port string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	JWT struct {
		Secret   string
		Lifetime time.Duration
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	// HTTP
	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.HTTP.Port = os.Getenv("HTTP_PORT")

	// DB
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.Name = os.Getenv("DB_NAME")

	// JWT
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	cfg.JWT.Lifetime = 24 * time.Hour

	return cfg, nil
}
