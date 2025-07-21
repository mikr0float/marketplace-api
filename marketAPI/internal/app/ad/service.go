package ad

import (
	"marketAPI/internal/storage"
)

type Service struct {
	adRepo storage.AdRepository
}

func NewService(adRepository storage.AdRepository) *Service {
	return &Service{adRepo: adRepository}
}
