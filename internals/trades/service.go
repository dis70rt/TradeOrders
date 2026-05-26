package trades

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	log "github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/dis70rt/TradeOrders/kafka"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetTrades(ctx context.Context, limit, page int) ([]Trade, error) {
    offset := (page - 1) * limit
    return s.repo.GetTrades(ctx, limit, offset)
}

func (s *Service) UpdateDatabase(ctx context.Context) {
	handler := kafka.ConsumerHandler{
		Process: func(msg *sarama.ConsumerMessage) error {
			var trade Trade
			if err := json.Unmarshal(msg.Value, &trade); err != nil {
				log.WithError(err).Error("failed to unmarshal trade")
				return err
			}
			if err := s.repo.ApplyTrade(ctx, &trade); err != nil {
				log.WithError(err).Error("failed to apply trade to database")
				return err
			}
			log.Info("Trade applied to database")
			return nil
		},
	}

	consumer := kafka.NewConsumer("TRADE_EXECUTED", "persistence-trades", handler)
	defer consumer.Close()
	consumer.Start()
}
