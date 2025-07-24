package rest

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/health" {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing auth header")
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth header format")
			}

			userID, err := h.authService.ValidateToken(tokenStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			c.Set("userID", userID)
			return next(c)
		}
	}
}
