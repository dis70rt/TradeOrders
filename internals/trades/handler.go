package trades

import (
    "net/http"

    log "github.com/dis70rt/TradeOrders/internals/logger"
    "github.com/dis70rt/TradeOrders/internals/utils"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) GetTrades(c *gin.Context) {
    limit, page := utils.ParsePaginationParams(c)

    ctx := c.Request.Context()
    trades, err := h.service.GetTrades(ctx, limit, page)
    if err != nil {
        log.WithError(err).Error("Failed to retrieve trades")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve trades"})
        return
    }

    c.JSON(http.StatusOK, trades)
}