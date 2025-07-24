package rest

import (
	_ "marketAPI/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (h *Handlers) SetupRoutes(e *echo.Echo) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Public routes
	e.GET("/health", h.HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Auth routes
	authGroup := e.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)

	// Ad routes
	adGroup := e.Group("/ads")
	adGroup.GET("", h.ListAds)
	adGroup.GET("/:id", h.GetAd)

	// Protected routes
	protectedAdGroup := adGroup.Group("")
	protectedAdGroup.Use(h.AuthMiddleware())
	protectedAdGroup.POST("", h.CreateAd)
}
