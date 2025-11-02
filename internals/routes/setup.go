package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	routerGroup := router.Group("/api/v1")
	{
		RegisterOrderRoutes(db, routerGroup)
	}
	return router
}