package app

import (
	"log"
	"marketAPI/internal/app/ad"
	"marketAPI/internal/app/auth"
	"marketAPI/internal/app/user"
	"marketAPI/internal/config"
	"marketAPI/internal/db"
	"marketAPI/internal/storage"
)

type Application struct {
	Auth *auth.Service
	User *user.Service
	Ad   *ad.Service
}

func NewApplication(cfg *config.Config) (*Application, error) {
	// Инициализация хранилища
	pgStorage, err := db.NewPostgresStorage(cfg.DB.DSN())
	if err != nil {
		log.Printf("Run Application failed")
		return nil, err
	}

	// Инициализация репозиториев
	userRepo := storage.NewUserRepository(pgStorage)
	adRepo := storage.NewAdRepository(pgStorage)

	// Инициализация сервисов
	authService := auth.NewService(
		userRepo,
		cfg.JWT.Secret,
		cfg.JWT.Lifetime,
	)
	userService := user.NewService(userRepo)
	adService := ad.NewService(adRepo, userRepo)

	return &Application{
		Auth: authService,
		User: userService,
		Ad:   adService,
	}, nil
}
