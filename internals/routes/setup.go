package routes

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(ctx context.Context, db *sql.DB, rdb *redis.Client) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	routerGroup := router.Group("/api/v1")
	{
		RegisterOrderRoutes(ctx, db, rdb, routerGroup)
	}
	return router
}