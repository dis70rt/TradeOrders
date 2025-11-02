package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func IdempotencyKey(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header is required"})
			return
		}

		ctx := context.Background()
		exists, err := redisClient.Exists(ctx, idempotencyKey).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			return
		}

		if exists > 0 {
			val, _ := redisClient.Get(ctx, idempotencyKey).Result()
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Duplicate request", "result": val})
			return
		}

		c.Next()
		status := c.Writer.Status()

		if status == http.StatusOK || status == http.StatusCreated {
			orderID := c.GetString("order_id")
			redisClient.Set(ctx, idempotencyKey, orderID, time.Hour)
		}
	}
}