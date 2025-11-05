package orderbook

import (
	"github.com/dis70rt/TradeOrders/internals/orders"
	"github.com/dis70rt/TradeOrders/internals/trades"
	// log "github.com/dis70rt/TradeOrders/internals/logger"
)

func (ob *OrderBook) canMatch(order *orders.MatchOrder) bool {
	switch order.Side {
	case "BUY":
		if ob.SellOrders.Len() == 0 {
			return false
		}
		bestSell := ob.SellOrders.Peek()
		return order.Price >= bestSell.Price
	
	case "SELL":
		if ob.BuyOrders.Len() == 0 {
			return false
		}
		bestBuy := ob.BuyOrders.Peek()
		return order.Price <= bestBuy.Price
	}
	return false
}

func (ob *OrderBook) TryMatch(order *orders.MatchOrder) {
    if order.Type == "MARKET" {
        ob.matchMarketOrder(order)
    } else {
        ob.matchLimitOrder(order)
    }
}

func (ob *OrderBook) matchLimitOrder(order *orders.MatchOrder) {
    for ob.canMatch(order) {
        if order.Quantity == 0 {
            break
        }

        switch order.Side {
        case "BUY":
            bestSell := ob.SellOrders.Peek()
            tradeQty := min(order.Quantity, bestSell.Quantity)

            trade := trades.New(bestSell, order, tradeQty, bestSell.Price)
            ob.TradeOut <- trade

            order.Quantity -= tradeQty
            bestSell.Quantity -= tradeQty
            if bestSell.Quantity == 0 {
                ob.SellOrders.Pop()
            }

        case "SELL":
            bestBuy := ob.BuyOrders.Peek()
            tradeQty := min(order.Quantity, bestBuy.Quantity)

            trade := trades.New(bestBuy, order, tradeQty, bestBuy.Price)
            ob.TradeOut <- trade

            order.Quantity -= tradeQty
            bestBuy.Quantity -= tradeQty
            if bestBuy.Quantity == 0 {
                ob.BuyOrders.Pop()
            }
        }
    }
}

func (ob *OrderBook) matchMarketOrder(order *orders.MatchOrder) {
    switch order.Side {
    case "BUY":
        for order.Quantity > 0 && ob.SellOrders.Len() > 0 {
            bestSell := ob.SellOrders.Peek()
            tradeQty := min(order.Quantity, bestSell.Quantity)

            trade := trades.New(bestSell, order, tradeQty, bestSell.Price)
            ob.TradeOut <- trade

            order.Quantity -= tradeQty
            bestSell.Quantity -= tradeQty
            if bestSell.Quantity == 0 {
                ob.SellOrders.Pop()
            }
        }

    case "SELL":
        for order.Quantity > 0 && ob.BuyOrders.Len() > 0 {
            bestBuy := ob.BuyOrders.Peek()
            tradeQty := min(order.Quantity, bestBuy.Quantity)

            trade := trades.New(bestBuy, order, tradeQty, bestBuy.Price)
            ob.TradeOut <- trade

            order.Quantity -= tradeQty
            bestBuy.Quantity -= tradeQty
            if bestBuy.Quantity == 0 {
                ob.BuyOrders.Pop()
            }
        }
    }
}