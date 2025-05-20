package handler

import (
	"net/http"
	"warehouse-service/controller"
	"warehouse-service/model"

	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	WarehouseController *controller.WarehouseController
}

func (h *WarehouseHandler) AddWarehouse(c *gin.Context) {
	var warehouse model.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if err := h.WarehouseController.AddWarehouse(&warehouse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add warehouse"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "warehouse added"})
}

func (h *WarehouseHandler) ToggleActive(c *gin.Context) {
	var body struct {
		ID     uint `json:"id"`
		Active bool `json:"active"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.WarehouseController.ToggleActive(body.ID, body.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}
