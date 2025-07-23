package config

import (
	"fmt"
	"os"
	"time"
)

type HTTP struct {
	Host string
	Port string
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWT struct {
	Secret   string
	Lifetime time.Duration
}

type Config struct {
	HTTP HTTP
	DB   DB
	JWT  JWT
}

func (c *Config) DataBase() *DB {
	return &c.DB
}

// Address возвращает строку адреса для HTTP сервера в формате "host:port"
func (c *Config) HTTPAddress() string {
	return fmt.Sprintf("%s:%s", c.HTTP.Host, c.HTTP.Port)
}

func (db *DB) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)
}

func Load() (*Config, error) {
	cfg := &Config{}

	// HTTP
	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	if cfg.HTTP.Host == "" {
		cfg.HTTP.Host = "0.0.0.0" // default value
	}

	cfg.HTTP.Port = os.Getenv("HTTP_PORT")
	if cfg.HTTP.Port == "" {
		cfg.HTTP.Port = "8080"
	}

	// DB
	cfg.DB.Host = os.Getenv("DB_HOST")
	if cfg.DB.Host == "" {
		cfg.DB.Host = "localhost"
	}

	cfg.DB.Port = os.Getenv("DB_PORT")
	if cfg.DB.Port == "" {
		cfg.DB.Port = "5432"
	}

	cfg.DB.User = os.Getenv("DB_USER")
	if cfg.DB.User == "" {
		cfg.DB.User = "postgres"
	}

	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	if cfg.DB.Password == "" {
		cfg.DB.Password = "postgres"
	}

	cfg.DB.Name = os.Getenv("DB_NAME")
	if cfg.DB.Name == "" {
		cfg.DB.Name = "marketplace"
	}

	// JWT
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	cfg.JWT.Lifetime = 24 * time.Hour

	return cfg, nil
}
