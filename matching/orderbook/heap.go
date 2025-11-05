package orderbook

import (
	"github.com/dis70rt/TradeOrders/internals/orders"
)

type OrderHeap struct {
	data []*orders.MatchOrder
	Price  int
	IsBuy  bool
}

func (h *OrderHeap) Len() int { return len(h.data) }

func (h *OrderHeap) Less(i, j int) bool {
	if h.IsBuy {
		return h.data[i].Price > h.data[j].Price
	}
	return h.data[i].Price < h.data[j].Price 
}

func (h *OrderHeap) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *OrderHeap) Push(x any) {
	order := x.(*orders.MatchOrder)
	h.data = append(h.data, order)
}

func (h *OrderHeap) Pop() any {
	old := h.data
	n := len(old)
	item := old[n-1]
	h.data = old[0 : n-1]
	return item
}

func (h *OrderHeap) Peek() *orders.MatchOrder {
    if len(h.data) == 0 {
        return nil
    }
    return h.data[0]
}
