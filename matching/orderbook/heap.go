package orderbook

import (
	"github.com/dis70rt/TradeOrders/internals/orders"
)

type OrderHeap struct {
	Orders []orders.MatchOrder
	Price  int
	IsBuy  bool
}

func (h *OrderHeap) Len() int { return len(h.Orders) }

func (h *OrderHeap) Less(i, j int) bool {
	if h.IsBuy {
		return h.Orders[i].Price > h.Orders[j].Price
	}
	return h.Orders[i].Price < h.Orders[j].Price 
}

func (h *OrderHeap) Swap(i, j int) {
	h.Orders[i], h.Orders[j] = h.Orders[j], h.Orders[i]
}

func (h *OrderHeap) Push(x any) {
	order := x.(orders.MatchOrder)
	h.Orders = append(h.Orders, order)
}

func (h *OrderHeap) Pop() any {
	old := h.Orders
	n := len(old)
	item := old[n-1]
	h.Orders = old[0 : n-1]
	return item
}

func (h *OrderHeap) Peek() *orders.MatchOrder {
    if len(h.Orders) == 0 {
        return nil
    }
    return &h.Orders[0]
}
