package engine

import (
	"encoding/json"
	"sync"

	"github.com/IBM/sarama"
	log "github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/dis70rt/TradeOrders/internals/orders"
	"github.com/dis70rt/TradeOrders/internals/trades"
	"github.com/dis70rt/TradeOrders/kafka"
	"github.com/dis70rt/TradeOrders/matching/orderbook"
)

type Engine struct {
	OrderChannels map[string]chan *orders.MatchOrder
	TradeChannels chan *trades.Trade
	mutex sync.RWMutex
}

func (engine *Engine) getOrderChannel(instrument string) chan *orders.MatchOrder {
	engine.mutex.RLock()
	ch, exists := engine.OrderChannels[instrument]
	engine.mutex.RUnlock()

	if exists {
		return ch
	}

	engine.mutex.Lock()
	defer engine.mutex.Unlock()

	if ch, exists := engine.OrderChannels[instrument]; exists {
        return ch
    }

	newCh := make(chan *orders.MatchOrder, 100)
	engine.OrderChannels[instrument] = newCh
	
	go engine.runInstrumentProcessor(instrument, newCh)
	return newCh
}

func (engine *Engine) runInstrumentProcessor(instrument string, orderChan <- chan *orders.MatchOrder)  {
	book := orderbook.NewOrderBook(instrument)
	book.TradeOut = engine.TradeChannels

	for order := range orderChan {
		book.TryMatch(order)
		if order.Quantity > 0 {
			book.AddOrder(order)
		}
	}
}


func StartMatchingEngine() {
	producer := kafka.NewProducer()
	// defer producer.Close()

	engine := &Engine{
		OrderChannels: make(map[string]chan *orders.MatchOrder),
		TradeChannels: make(chan *trades.Trade, 256),
	}

	tradeExec := NewTradeExecutor(engine.TradeChannels, producer)
	tradeExec.Start()

	handler := kafka.ConsumerHandler{
		Process: func(msg *sarama.ConsumerMessage) {
			var order orders.MatchOrder
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.WithError(err).Error("failed to unmarshal order")
				return
			}

			log.Infof("Received order: %+v", order)
			orderCh := engine.getOrderChannel(order.Instrument)
			orderCh <- &order
		},
	}

	log.Info("Matching engine dispatcher started. Waiting for orders...")
	consumer := kafka.NewConsumer("orders.inbound", "matching-engine-group", handler)
	defer consumer.Close()
	consumer.Start()
}