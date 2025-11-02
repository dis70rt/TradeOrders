package orders

import (
	"context"
	"errors"

	log "github.com/dis70rt/TradeOrders/internals/logger"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(ctx context.Context, order *OrderRequest) (string, error) {
	if !order.Validate() {
		log.Info("Invalid order data")
		return "", errors.New("Invalid to order details")
	}
	return s.repo.CreateOrder(ctx, order)
}