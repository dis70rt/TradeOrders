package trades

import (
	"time"

	"github.com/dis70rt/TradeOrders/internals/orders"
	"github.com/google/uuid"
)

type Trade struct {
	ID          uuid.UUID `json:"id"`
	BuyOrderID  uuid.UUID `json:"buy_order_id"`
	SellOrderID uuid.UUID `json:"sell_order_id"`
	Instrument  string    `json:"instrument"`
	Price       float64   `json:"price"`
	Quantity    float64   `json:"quantity"`
	ExecutedAt  time.Time `json:"executed_at"`
}

func New(buy, sell *orders.MatchOrder, quantity float64, price float64) *Trade {
	tradeID := uuid.New()
	return &Trade{
		ID: tradeID,
		BuyOrderID: buy.ID,
		SellOrderID: sell.ID,
		Instrument: sell.Instrument,
		Price: price,
		Quantity: quantity,
		ExecutedAt: time.Now(),
	}
}