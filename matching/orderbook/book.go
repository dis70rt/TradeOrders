package orderbook

import (
	"container/heap"

	"github.com/dis70rt/TradeOrders/internals/orders"
	"github.com/dis70rt/TradeOrders/internals/trades"
)

type OrderBookManager struct {
    Books map[string]*OrderBook
}

type OrderBook struct {
	Instrument	string
	BuyOrders 	*OrderHeap
	SellOrders 	*OrderHeap
	TradeOut 	chan *trades.Trade
	OrderIn 	chan *orders.MatchOrder
}

func NewManager() *OrderBookManager {
    return &OrderBookManager{
        Books: make(map[string]*OrderBook),
    }
}

func NewOrderBook(instrument string) *OrderBook {
	buy := &OrderHeap{IsBuy: true}
	sell := &OrderHeap{IsBuy: false}
	heap.Init(buy)
	heap.Init(sell)
	return &OrderBook{
		Instrument: instrument,
		BuyOrders:  buy,
		SellOrders: sell,
		TradeOut:   make(chan *trades.Trade, 1024),
        OrderIn:    make(chan *orders.MatchOrder, 1024),
	}
}

func (m *OrderBookManager) GetOrCreate(instrument string) *OrderBook {
    ob, ok := m.Books[instrument]
    if !ok {
        ob = NewOrderBook(instrument)
        m.Books[instrument] = ob
    }
    return ob
}

func (ob *OrderBook) AddOrder(order *orders.MatchOrder) {
	if order.Side == "BUY" {
		heap.Push(ob.BuyOrders, order)
	} else {
		heap.Push(ob.SellOrders, order)
	}
}