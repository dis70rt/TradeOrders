package main

import (
	"github.com/dis70rt/TradeOrders/internals/database"
	"github.com/dis70rt/TradeOrders/internals/routes"
)

func main() {
	db := database.ConnectPostgres()
	defer db.Close()

	rdb := database.ConnectRedis()
	defer rdb.Close()

	routes.SetupRoutes(db, rdb).Run(":8080")
}