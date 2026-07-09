package models

import "time"

type Source struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	APIName   string    `json:"api_name"`
	BaseURL   string    `json:"base_url"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}
