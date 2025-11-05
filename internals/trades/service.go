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

func (s *Service) UpdateDatabase(ctx context.Context) {
	handler := kafka.ConsumerHandler{
		Process: func(msg *sarama.ConsumerMessage) {
			var trade Trade
			if err := json.Unmarshal(msg.Value, &trade); err != nil {
				log.WithError(err).Error("failed to unmarshal order")
				return
			}
			s.repo.ApplyTrade(ctx, &trade)
			log.Info("Trade applied to database")
		},
	}

	consumer := kafka.NewConsumer("TRADE_EXECUTED", "persistence-trades", handler)
	defer consumer.Close()
	consumer.Start()
}
