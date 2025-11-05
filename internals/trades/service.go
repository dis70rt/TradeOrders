package trades

import (
	"encoding/json"

	"github.com/dis70rt/TradeOrders/kafka"
)

type Service struct {
	producer *kafka.Producer
}

func NewService(producer *kafka.Producer) *Service {
	return &Service{
		producer: producer,
	}
}

func (s *Service) Exec(trade *Trade) error {
	tradeJSON, err := json.Marshal(trade)
	if err != nil {
		return err
	}

	s.producer.SendMessage("trades", trade.Instrument, tradeJSON)
	return nil
}
