package domain

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AdResponse struct {
	Ads      []Ad `json:"ads"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Total    int  `json:"total"`
}

type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
	Timestamp time.Time         `json:"timestamp"`
}
