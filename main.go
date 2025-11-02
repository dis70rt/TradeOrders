package main

import (
	"github.com/dis70rt/TradeOrders/internals/database"
	"github.com/dis70rt/TradeOrders/internals/routes"
)

func main() {
	db := database.Connect()
	defer db.Close()

	routes.SetupRoutes(db).Run(":8080")
}