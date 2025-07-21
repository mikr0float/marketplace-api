package user

import (
	"marketAPI/internal/domain"
	"marketAPI/internal/storage"
)

type Service struct {
	userRepo storage.UserRepository
}

func NewService(userRepository storage.UserRepository) *Service {
	return &Service{userRepo: userRepository}
}

func (s *Service) Register(username, password string) (*domain.User, error) {
	// Валидация username и password
	// Хеширование пароля
	// Создание пользователя
	return &domain.User{}, nil // ЗАГЛУШКАА!!!!!!!!!!!
}
