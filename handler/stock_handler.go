package handler

import (
	"net/http"
	"strconv"
	"warehouse-service/controller"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	StockController *controller.StockController
}

func (h *StockHandler) GetStockByProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	total, list, err := h.StockController.GetProductStock(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product_id":  productID,
		"total_stock": total,
		"warehouses":  list,
	})
}

func (h *StockHandler) TransferStock(c *gin.Context) {
	var body struct {
		FromWarehouseID uint `json:"from_warehouse_id"`
		ToWarehouseID   uint `json:"to_warehouse_id"`
		ProductID       uint `json:"product_id"`
		Quantity        int  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.StockController.TransferProductStock(body.FromWarehouseID, body.ToWarehouseID, body.ProductID, body.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock transferred successfully"})
}

func (h *StockHandler) ReserveStock(c *gin.Context) {
	var body struct {
		WarehouseID uint `json:"warehouse_id"`
		ProductID   uint `json:"product_id"`
		Quantity    int  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	if err := h.StockController.ReserveProductStock(body.WarehouseID, body.ProductID, body.Quantity); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "stock reserved"})
}

func (h *StockHandler) ReleaseStock(c *gin.Context) {
	var body struct {
		WarehouseID uint `json:"warehouse_id"`
		ProductID   uint `json:"product_id"`
		Quantity    int  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	if err := h.StockController.ReleaseProductStock(body.WarehouseID, body.ProductID, body.Quantity); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "reserved stock released"})
}

func (h *StockHandler) GetStockWithProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	total, list, err := h.StockController.GetProductStockWithDetail(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"product_id":  productID,
		"total_stock": total,
		"warehouses":  list,
	})
}
