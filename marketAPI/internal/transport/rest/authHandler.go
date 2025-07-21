package rest

import (
	"marketAPI/internal/app/auth"
	"net/http"
)

type authHandler struct {
	authServide *auth.Service
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) { // Сделать через echo!!!!
	// Парсинг запроса
	// Вызов authService.Login
	// Возврат токена
}
