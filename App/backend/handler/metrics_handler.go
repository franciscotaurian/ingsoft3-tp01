package handler

import (
	"net/http"

	"realico-comidas-backend/service"

	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	orderService service.OrderService
}

func NewMetricsHandler(orderService service.OrderService) *MetricsHandler {
	return &MetricsHandler{orderService: orderService}
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	metrics, err := h.orderService.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}
