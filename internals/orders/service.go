package orders

import (
	"context"
	"encoding/json"
	"errors"

	log "github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/dis70rt/TradeOrders/kafka"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	producer *kafka.Producer
}

func NewService(repo *Repository, producer *kafka.Producer) *Service {
	return &Service{repo: repo, producer: producer}
}

func (s *Service) CreateOrder(ctx context.Context, order *OrderRequest) (string, error) {
	if !order.Validate() {
		log.Info("Invalid order data")
		return "", errors.New("Invalid to order details")
	}
	orderID := uuid.New()
	matchOrder := MatchOrder{
		ID: orderID,
		Instrument: order.Instrument,
		Side: order.Side,
		Type: order.Type,
		Price: order.Price,
		Quantity: order.Quantity,
	}

	orderJSON, _ := json.Marshal(matchOrder) 
	s.producer.SendMessage("orders", order.Instrument, orderJSON)

	return s.repo.CreateOrder(ctx, orderID, order)
}

func (s *Service) GetOrders(ctx context.Context, limit, page int) ([]Order, error) {
	return s.repo.GetOrders(ctx, limit, page)
}	