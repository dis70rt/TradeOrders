package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParsePaginationParams(c *gin.Context) (limit, page int) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	var err error
	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err = strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	return limit, page
}