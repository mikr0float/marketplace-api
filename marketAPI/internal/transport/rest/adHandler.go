package rest

import (
	"marketAPI/internal/app/ad"
	"marketAPI/internal/app/auth"
	"net/http"
)

type adHandler struct {
	adService   *ad.Service
	authService *auth.Service
}

func (h *adHandler) CreateAd(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	// Парсинг запроса
	// Валидация данных
	// Создание объявления
}

func (h *adHandler) ListAds(w http.ResponseWriter, r *http.Request) {
	// Парсинг параметров запроса (фильтры, сортировка)
	// Получение списка объявлений
}
