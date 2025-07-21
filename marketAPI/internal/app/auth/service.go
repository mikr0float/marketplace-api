package auth

import (
	"marketAPI/internal/storage"
)

type Service struct {
	userRepo    storage.UserRepository
	tokenSecret string
}

func NewService(userRepository storage.UserRepository, secret string) *Service {
	return &Service{
		userRepo:    userRepository,
		tokenSecret: secret,
	}
}

func (s *Service) Login(username, password string) (string, error) {
	// Проверка учетных данных
	// Генерация JWT токена
	return "", nil // ЗАГЛУШКА!!!!!!!!
}

func (s *Service) ValidetaToken(tokenString string) (int, error) {
	// Валидация токена и возврат userID
	return 0, nil // ЗАГЛУШКА!!!!!!!!
}
