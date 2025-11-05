package routes

import (
	"database/sql"

	"github.com/dis70rt/TradeOrders/internals/middlewares"
	"github.com/dis70rt/TradeOrders/internals/orders"
	"github.com/dis70rt/TradeOrders/kafka"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterOrderRoutes(db *sql.DB, rdb *redis.Client, r *gin.RouterGroup) {
	order := r.Group("/orders")
	producer := kafka.NewProducer()
	defer producer.Close()
	{	
		repo := orders.NewRepository(db)
		service := orders.NewService(repo, producer)
		handler := orders.NewHandler(service)

		order.POST("/", middlewares.IdempotencyKey(rdb), handler.CreateOrder)
		order.GET("/", handler.GetOrders)
	}
}