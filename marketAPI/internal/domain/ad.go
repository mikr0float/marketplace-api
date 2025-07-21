package domain

import (
	"time"
)

type Ad struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	Price       float64   `json:"price" db:"price"`
	UserID      int       `json:"-" db:"user_id"` // Мб сделать через id google
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	//Для ответа:
	Username string `json:"username,omitempty"`
	IsOwner  bool   `json:"is_owner,omitempty"`
}

type AdRequest struct {
	Title       string  `json:"title" validate:"required,min=5,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=1000"`
	ImageURL    string  `json:"image_url" validate:"required,url"`
	Price       float64 `json:"price" validate:"required,min=0"`
}

type AdFilter struct {
	Page      int
	PageSize  int
	SortBy    string // "date" or "price"
	SortOrder string // "asc" or "desc"
	MinPrice  *float64
	MaxPrice  *float64
}
