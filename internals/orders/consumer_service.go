package orders

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	log "github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/dis70rt/TradeOrders/kafka"
)

type ConsumerService struct {
	repo *Repository
}

func NewConsumerService(repo *Repository) *ConsumerService {
	return &ConsumerService{repo: repo}
}

func (s *ConsumerService) InsertOrderDatabase(ctx context.Context) {
	handler := kafka.ConsumerHandler{
		Process: func(msg *sarama.ConsumerMessage) {
			var order MatchOrder
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.WithError(err).Error("failed to unmarshal order")
				return
			}

			s.repo.CreateOrder(ctx, &order)
			log.Info("Order inserted into database")
		},
	}

	consumer := kafka.NewConsumer("ORDER_ACCEPTED", "persistence-trades", handler)
	defer consumer.Close()
	consumer.Start()
}