package user

import (
	"context"
	"marketAPI/internal/domain"
	"marketAPI/internal/storage"
)

type Service struct {
	userRepo *storage.UserRepository
}

func NewService(userRepository *storage.UserRepository) *Service {
	return &Service{userRepo: userRepository}
}

func (s *Service) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *Service) Exists(ctx context.Context, username string) (bool, error) {
	return s.userRepo.Exists(ctx, username)
}
