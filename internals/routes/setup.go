package routes

import (
	"database/sql"
	"net/http"

	"github.com/dis70rt/TradeOrders/internals/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(db *sql.DB, rdb *redis.Client) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.Use(middlewares.IdempotencyKey(rdb))

	routerGroup := router.Group("/api/v1")
	{
		RegisterOrderRoutes(db, routerGroup)
	}
	return router
}