package models

import "time"

type Portfolio struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Amount    float64   `json:"amount"`
	ValueUSD  float64   `json:"value_usd"`
	CreatedAt time.Time `json:"created_at"`
}
