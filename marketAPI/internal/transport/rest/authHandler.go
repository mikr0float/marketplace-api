package rest

import (
	"marketAPI/internal/app/ad"
	"marketAPI/internal/app/auth"
	"marketAPI/internal/app/user"
	"marketAPI/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	authService *auth.Service
	userService *user.Service
	adService   *ad.Service
}

func NewHandlers(auth *auth.Service, user *user.Service, ad *ad.Service) *Handlers {
	return &Handlers{
		authService: auth,
		userService: user,
		adService:   ad,
	}
}

// Register godoc
// @Summary Register new user
// @Description Register new user with username and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body domain.AuthRequest true "User credentials"
// @Success 201 {object} domain.User
// @Failure 400 {object} domain.ErrorResponse
// @Failure 409 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /auth/register [post]
func (h *Handlers) Register(c echo.Context) error {
	var req domain.AuthRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "invalid request body",
		})
	}

	user, err := h.authService.Register(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		switch err {
		case domain.ErrUserExists:
			return c.JSON(http.StatusConflict, domain.ErrorResponse{
				Error: "user already exists",
			})
		case domain.ErrInvalidUsername, domain.ErrInvalidPassword:
			return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
				Error: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Error: "registration failed",
			})
		}
	}

	return c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login user
// @Description Login user with username and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body domain.AuthRequest true "User credentials"
// @Success 200 {object} domain.AuthResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /auth/login [post]
func (h *Handlers) Login(c echo.Context) error {
	var req domain.AuthRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "invalid request body",
		})
	}

	token, err := h.authService.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			return c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
				Error: "invalid credentials",
			})
		default:
			return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Error: "login failed",
			})
		}
	}

	return c.JSON(http.StatusOK, domain.AuthResponse{
		Token: token,
	})
}
