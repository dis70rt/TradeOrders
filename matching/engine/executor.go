package engine

import (
	"encoding/json"
	"sync"

	log "github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/dis70rt/TradeOrders/internals/trades"
	"github.com/dis70rt/TradeOrders/kafka"
)

type TradeExecutor struct {
	TradeCh chan *trades.Trade
	Producer *kafka.Producer
	wg       sync.WaitGroup
}

func NewTradeExecutor(tradeCh chan *trades.Trade, producer *kafka.Producer) *TradeExecutor {
	return &TradeExecutor{
		TradeCh: tradeCh,
		Producer: producer,
	}
}

func (te *TradeExecutor) Start() {
	log.Info("Trade executor started. Waiting for trades to publish.")
	te.wg.Add(1)
	go func() {
		defer te.wg.Done()
		for trade := range te.TradeCh {
			tradeJSON, err := json.Marshal(trade)
			if err != nil {
                log.WithError(err).Error("Failed to marshal trade for publishing")
                continue
            }

			log.Infof("Publishing trade: %+v", trade)
			te.Producer.SendMessage("TRADE_EXECUTED", trade.Instrument, tradeJSON)
		}
	}()
}

func (te *TradeExecutor) Stop() {
	te.wg.Wait()
	te.Producer.Close()
}
