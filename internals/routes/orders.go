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
	order := r.Group("/orders")
	producer := kafka.NewProducer()
	// defer producer.Close()
	{	
		repo := orders.NewRepository(db, producer)
		service := orders.NewService(repo)
		handler := orders.NewHandler(service)

		order.POST("/", middlewares.IdempotencyKey(rdb), handler.CreateOrder)
		order.GET("/", handler.GetOrders)
	}

	{
		repo := trades.NewRepository(db)
		service := trades.NewService(repo)
		go service.UpdateDatabase(ctx)
	}
}