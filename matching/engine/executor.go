package engine

import (
	"encoding/json"

	log "github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/dis70rt/TradeOrders/internals/trades"
	"github.com/dis70rt/TradeOrders/kafka"
)

type TradeExecutor struct {
	TradeCh chan *trades.Trade
	Producer *kafka.Producer
}

func NewTradeExecutor(tradeCh chan *trades.Trade, producer *kafka.Producer) *TradeExecutor {
	return &TradeExecutor{
		TradeCh: tradeCh,
		Producer: producer,
	}
}

func (te *TradeExecutor) Start() {
	log.Info("Trade executor started. Waiting for trades to publish.")
	go func() {
		for trade := range te.TradeCh {
			tradeJSON, err := json.Marshal(trade)
			if err != nil {
                log.WithError(err).Error("Failed to marshal trade for publishing")
                continue
            }

			te.Producer.SendMessage("trades", trade.Instrument, tradeJSON)
		}
	}()
}
