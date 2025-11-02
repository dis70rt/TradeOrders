package orders

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

func (h *Handler) CreateOrder(c *gin.Context) {
	var order OrderRequest
	if err := c.ShouldBindBodyWithJSON(&order); err != nil {
		log.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	orderID, err := h.service.CreateOrder(ctx, &order)
	if err != nil {
		log.WithError(err).Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order_id": orderID})
}

func (h *Handler) GetOrders(c *gin.Context) {
	limit, page := utils.ParsePaginationParams(c)

	ctx := c.Request.Context()
	orders, err := h.service.GetOrders(ctx, limit, page)
	if err != nil {
		log.WithError(err).Error("Failed to retrieve orders")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}