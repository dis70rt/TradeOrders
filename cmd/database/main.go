package main

import (
	"context"

	"github.com/dis70rt/TradeOrders/internals/database"
	log "github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/dis70rt/TradeOrders/internals/orders"
	"github.com/dis70rt/TradeOrders/internals/trades"
)

func main() {
	log.Init()

	log.Infof("Starting database service...")
	db := database.ConnectPostgres()
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	{
		repo := trades.NewRepository(db)
		service := trades.NewService(repo)
		go service.UpdateDatabase(ctx)

	}
	{
		repo := orders.NewRepository(db)
		service := orders.NewConsumerService(repo)
		service.InsertOrderDatabase(ctx)
	}
}