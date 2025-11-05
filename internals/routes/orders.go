package routes

import (
	"context"
	"database/sql"

	"github.com/dis70rt/TradeOrders/internals/middlewares"
	"github.com/dis70rt/TradeOrders/internals/orders"
	"github.com/dis70rt/TradeOrders/internals/trades"
	"github.com/dis70rt/TradeOrders/kafka"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterOrderRoutes(ctx context.Context, db *sql.DB, rdb *redis.Client, r *gin.RouterGroup) {
	ordersGroup := r.Group("/orders")
	producer := kafka.NewProducer()
	// defer producer.Close()
	{	
		repo := orders.NewRepository(db)
		service := orders.NewService(repo, producer)
		handler := orders.NewHandler(service)

		ordersGroup.POST("/", middlewares.IdempotencyKey(rdb), handler.CreateOrder)
		ordersGroup.GET("/", handler.GetOrders)
		ordersGroup.GET("/:id", handler.GetOrderByID)
	}

	tradesGroup := r.Group("/trades")
	{
		repo := trades.NewRepository(db)
		service := trades.NewService(repo)
		handler := trades.NewHandler(service)

		tradesGroup.GET("/", handler.GetTrades)
	}
}