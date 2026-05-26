package orders

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID    `json:"id"`
	ClientID   uuid.UUID    `json:"client_id"`
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
	ClientID   uuid.UUID    `json:"client_id"`
	Instrument string    `json:"instrument"`
	Side       string    `json:"side"`
	Type       string    `json:"type"`
	Price      float64   `json:"price"`
	Quantity   float64   `json:"quantity"`
	Remaining  float64   `json:"remaining"`
}

type MatchOrder struct {
	ID         uuid.UUID  `json:"id"`
	ClientID   uuid.UUID    `json:"client_id"`
	Instrument string  `json:"instrument"`
	Side       string  `json:"side"`
	Type       string  `json:"type"`
	Price      float64 `json:"price"`
	Quantity   float64 `json:"quantity"`
	Timestamp  int64   `json:"timestamp"`
}

func (o *OrderRequest) Validate() bool {
	validSides := map[string]bool{"buy": true, "sell": true}
	validTypes := map[string]bool{"limit": true, "market": true}

	oside := strings.ToLower(o.Side)
	otype := strings.ToLower(o.Type)

	if o.ClientID == uuid.Nil || o.Instrument == "" {
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
