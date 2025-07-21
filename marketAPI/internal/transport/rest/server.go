package rest

import (
	"marketAPI/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
}

func NewServer(cfg *config.Config) *Server {
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return &Server{
		echo: e,
		cfg:  cfg,
	}
}

func (s *Server) Start() error {
	return s.echo.Start(":" + s.cfg.HTTP.Port)
}
