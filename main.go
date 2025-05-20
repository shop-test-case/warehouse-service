package main

import (
	"warehouse-service/config"
	"warehouse-service/controller"
	"warehouse-service/database"
	"warehouse-service/handler"
	"warehouse-service/middleware"
	"warehouse-service/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)

	warehouseRepo := &repository.WarehouseRepo{DB: db}
	stockRepo := &repository.StockRepo{DB: db}

	warehouseController := &controller.WarehouseController{WarehouseRepo: warehouseRepo}
	stockController := &controller.StockController{
		StockRepo:     stockRepo,
		WarehouseRepo: warehouseRepo,
	}

	warehouseHandler := &handler.WarehouseHandler{WarehouseController: warehouseController}
	stockHandler := &handler.StockHandler{StockController: stockController}

	r := gin.Default()
	r.Use(cors.Default())

	auth := r.Group("/")
	auth.Use(middleware.JWT(cfg.JWTSecret))

	auth.POST("/warehouse", warehouseHandler.AddWarehouse)
	auth.PUT("/warehouse/active", warehouseHandler.ToggleActive)
	auth.GET("/stock/product/:product_id", stockHandler.GetStockByProduct)
	auth.GET("/stock/product-with-detail/:product_id", stockHandler.GetStockWithProduct)
	auth.POST("/stock/transfer", stockHandler.TransferStock)
	auth.POST("/stock/reserve", stockHandler.ReserveStock)
	auth.POST("/stock/release", stockHandler.ReleaseStock)

	r.Run(":" + cfg.Port)
}
