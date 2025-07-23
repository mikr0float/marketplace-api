package ad

import (
	"context"
	"marketAPI/internal/domain"
	"marketAPI/internal/storage"
)

type Service struct {
	adRepo   *storage.AdRepository
	userRepo *storage.UserRepository
}

func NewService(adRepo *storage.AdRepository, userRepo *storage.UserRepository) *Service {
	return &Service{
		adRepo:   adRepo,
		userRepo: userRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID int, req domain.AdRequest) (*domain.Ad, error) {
	// Валидация данных
	if len(req.Title) < 5 || len(req.Title) > 100 {
		return nil, domain.ErrInvalidAdTitle
	}
	if len(req.Description) < 10 || len(req.Description) > 1000 {
		return nil, domain.ErrInvalidAdDescription
	}
	if req.Price < 0 {
		return nil, domain.ErrInvalidAdPrice
	}

	ad := &domain.Ad{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		UserID:      userID,
	}

	if err := s.adRepo.Create(ctx, ad); err != nil {
		return nil, err
	}

	return ad, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*domain.Ad, error) {
	return s.adRepo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error) {
	// Установка значений по умолчанию
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.SortBy == "" {
		filter.SortBy = "date"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	ads, err := s.adRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return ads, nil
}
