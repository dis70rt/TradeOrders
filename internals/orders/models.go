package orders

import (
	"strings"
	"time"
)

type Order struct {
	ID         string    `json:"id"`
	ClientID   string    `json:"client_id"`
	Instrument string    `json:"instrument"`
	Side       string    `json:"side"`
	Type       string    `json:"type"`
	Price      float64   `json:"price"`
	Quantity   float64   `json:"quantity"`
	Remaining  float64   `json:"remaining"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderRequest struct {
	ClientID   string    `json:"client_id"`
	Instrument string    `json:"instrument"`
	Side       string    `json:"side"`
	Type       string    `json:"type"`
	Price      float64   `json:"price"`
	Quantity   float64   `json:"quantity"`
	Remaining  float64   `json:"remaining"`
}

func (o *OrderRequest) Validate() bool {
	validSides := map[string]bool{"buy": true, "sell": true}
	validTypes := map[string]bool{"limit": true, "market": true}

	oside := strings.ToLower(o.Side)
	otype := strings.ToLower(o.Type)

	if o.ClientID == "" || o.Instrument == "" {
		return false
	}

	if !validSides[oside] || !validTypes[otype] {
		return false
	}

	if o.Price <= 0 || o.Quantity <= 0 {
		return false
	}

	return true
}
