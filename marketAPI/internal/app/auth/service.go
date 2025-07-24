package auth

import (
	"context"
	"marketAPI/internal/domain"
	"marketAPI/internal/storage"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo    *storage.UserRepository
	tokenSecret string
	tokenExpiry time.Duration
}

func NewService(userRepository *storage.UserRepository, secret string, expiry time.Duration) *Service {
	return &Service{
		userRepo:    userRepository,
		tokenSecret: secret,
		tokenExpiry: expiry,
	}
}

func (s *Service) Register(ctx context.Context, username, password string) (*domain.User, error) {
	// Валидация username (3-50 символов, только буквы/цифры)
	if len(username) < 3 || len(username) > 50 {
		return nil, domain.ErrInvalidUsername
	}

	// Валидация password (минимум 8 символов)
	if len(password) < 8 {
		return nil, domain.ErrInvalidPassword
	}

	// Проверка существования пользователя
	exists, err := s.userRepo.Exists(ctx, username)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrUserExists
	}

	user := &domain.User{
		Username: username,
		Password: password,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", domain.ErrInvalidCredentials
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(s.tokenExpiry).Unix(),
	})

	return token.SignedString([]byte(s.tokenSecret))
}

func (s *Service) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, domain.ErrInvalidToken
		}
		return []byte(s.tokenSecret), nil
	})

	if err != nil {
		return 0, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		sub, ok := claims["sub"].(float64)
		if ok {
			return int(sub), nil
		}
	}

	return 0, domain.ErrInvalidToken
}

func (s *Service) GetSecret() string {
	return s.tokenSecret
}

func (s *Service) CheckDB(ctx context.Context) error {
	return s.userRepo.PingDB(ctx)
}
