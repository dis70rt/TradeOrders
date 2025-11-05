package main

import (
	"context"

	"github.com/dis70rt/TradeOrders/internals/database"
	"github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/dis70rt/TradeOrders/internals/routes"
)

func main() {
	db := database.ConnectPostgres()
	defer db.Close()

	rdb := database.ConnectRedis()
	defer rdb.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	logger.Init()
	routes.SetupRoutes(ctx, db, rdb).Run(":8080")
}