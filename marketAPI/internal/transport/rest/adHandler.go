package rest

import (
	/*"marketAPI/internal/app/ad"
	"marketAPI/internal/app/auth"*/
	"marketAPI/internal/domain"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

/*type adHandlers struct {
	adService   *ad.Service
	authService *auth.Service
}*/

func getUserIdFromToken(c echo.Context) (int, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return 0, echo.ErrUnauthorized
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return 0, echo.ErrUnauthorized
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, echo.ErrUnauthorized
	}

	return int(sub), nil
}

// CreateAd godoc
// @Summary Create new ad
// @Description Create new advertisement
// @Tags ads
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param input body domain.AdRequest true "Ad data"
// @Success 201 {object} domain.Ad
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /ads [post]
func (h *Handlers) CreateAd(c echo.Context) error {
	userID, err := getUserIdFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Error: domain.ErrUnauthorized.Error(),
		})
	}

	var req domain.AdRequest
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "invalid request body",
		})
	}

	ad, err := h.adService.Create(c.Request().Context(), userID, req)
	if err != nil {
		switch err {
		case domain.ErrInvalidAdTitle, domain.ErrInvalidAdDescription, domain.ErrInvalidAdPrice:
			return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
				Error: err.Error(),
			})
		default:
			return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Error: "failed to create ad",
			})
		}
	}

	return c.JSON(http.StatusCreated, ad)
}

// ListAds godoc
// @Summary List ads
// @Description Get list of advertisements with filtering and pagination
// @Tags ads
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param sort_by query string false "Sort field (date or price)" default(date)
// @Param sort_order query string false "Sort order (asc or desc)" default(desc)
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Success 200 {array} domain.Ad
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /ads [get]
func (h *Handlers) ListAds(c echo.Context) error {
	filter := domain.AdFilter{
		Page:      1,
		PageSize:  10,
		SortBy:    "date",
		SortOrder: "desc",
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err == nil && page > 0 {
		filter.Page = page
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err == nil && pageSize > 0 {
		filter.PageSize = pageSize
	}

	sortBy := c.QueryParam("sort_by")
	if sortBy != "" {
		filter.SortBy = sortBy
	}

	sortOrder := c.QueryParam("sort_order")
	if sortOrder != "" {
		filter.SortOrder = sortOrder
	}

	minPrice, err := strconv.ParseFloat(c.QueryParam("min_price"), 64)
	if err == nil {
		filter.MinPrice = &minPrice
	}

	maxPrice, err := strconv.ParseFloat(c.QueryParam("max_price"), 64)
	if err == nil {
		filter.MaxPrice = &maxPrice
	}

	ads, err := h.adService.List(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error: "failed to get ads list",
		})
	}

	// Добавляем признак принадлежности объявления текущему пользователю
	if userID, ok := c.Get("userID").(int); ok {
		for i := range ads {
			ads[i].IsOwner = ads[i].UserID == userID
		}
	}

	return c.JSON(http.StatusOK, ads)
}

// GetAd godoc
// @Summary Get ad by ID
// @Description Get advertisement details by ID
// @Tags ads
// @Accept  json
// @Produce  json
// @Param id path int true "Ad ID"
// @Success 200 {object} domain.Ad
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /ads/{id} [get]
func (h *Handlers) GetAd(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error: "invalid ad ID",
		})
	}

	ad, err := h.adService.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error: "failed to get ad",
		})
	}
	if ad == nil {
		return c.JSON(http.StatusNotFound, domain.ErrorResponse{
			Error: "ad not found",
		})
	}

	// Добавляем признак принадлежности
	if userID, ok := c.Get("userID").(int); ok {
		ad.IsOwner = ad.UserID == userID
	}

	return c.JSON(http.StatusOK, ad)
}
