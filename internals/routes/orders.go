package routes

import (
	"database/sql"

	"github.com/dis70rt/TradeOrders/internals/orders"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(db *sql.DB, r *gin.RouterGroup) {
	order := r.Group("/orders")
	{	
		repo := orders.NewRepository(db)
		service := orders.NewService(repo)
		handler := orders.NewHandler(service)
		order.POST("/", handler.CreateOrder)
	}
}