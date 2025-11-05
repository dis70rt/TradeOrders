package orderbook

import (
	"container/heap"
	"sync"

	"github.com/dis70rt/TradeOrders/internals/orders"
	"github.com/dis70rt/TradeOrders/internals/trades"
)

type OrderBookManager struct {
    Books map[string]*OrderBook
    mu    sync.Mutex
}

type OrderBook struct {
	mu sync.Mutex
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
    m.mu.Lock()
    defer m.mu.Unlock()
    ob, ok := m.Books[instrument]
    if !ok {
        ob = NewOrderBook(instrument)
        m.Books[instrument] = ob
    }
    return ob
}

func (ob *OrderBook) AddOrder(order *orders.MatchOrder) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if order.Side == "BUY" {
		heap.Push(ob.BuyOrders, order)
	} else {
		heap.Push(ob.SellOrders, order)
	}
}