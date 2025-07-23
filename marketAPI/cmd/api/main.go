// @title Marketplace API
// @version 1.0
// @description API for marketplace application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"marketAPI/internal/app"
	"marketAPI/internal/config"
	"marketAPI/internal/transport/rest"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Инициализация приложения
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatalf("Application init error: %s", err)
	}

	// Создание Echo инстанса
	e := rest.NewServer()

	// Инициализация обработчиков
	handlers := rest.NewHandlers(
		application.Auth,
		application.User,
		application.Ad,
	)

	// Настройка роутов
	handlers.SetupRoutes(e)

	// Запуск сервера
	log.Printf("Starting server on %s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if err := e.Start(cfg.HTTPAddress()); err != nil {
		log.Fatalf("Server error: %s", err)
	}
}
