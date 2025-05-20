package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"warehouse-service/model"
	"warehouse-service/repository"
)

type StockController struct {
	StockRepo     repository.IStockRepo
	WarehouseRepo repository.IWarehouseRepo
}

func (c *StockController) GetProductStock(productID uint) (int, []model.WarehouseStockDTO, error) {
	stocks, err := c.StockRepo.FindStockByProductID(productID)
	if err != nil {
		return 0, nil, err
	}

	total := 0
	var result []model.WarehouseStockDTO
	for _, stock := range stocks {
		warehouse, err := c.WarehouseRepo.FindByID(stock.WarehouseID)
		if err != nil || !warehouse.Active {
			continue
		}
		total += stock.Quantity - stock.ReservedQuantity
		result = append(result, model.WarehouseStockDTO{
			WarehouseID: warehouse.ID,
			Quantity:    stock.Quantity - stock.ReservedQuantity,
		})
	}

	return total, result, nil
}

func (c *StockController) TransferProductStock(fromWarehouseID, toWarehouseID, productID uint, quantity int) error {
	from, err := c.WarehouseRepo.FindByID(fromWarehouseID)
	if err != nil || !from.Active {
		return errors.New("source warehouse not available")
	}

	to, err := c.WarehouseRepo.FindByID(toWarehouseID)
	if err != nil || !to.Active {
		return errors.New("target warehouse not available")
	}

	return c.StockRepo.TransferStock(fromWarehouseID, toWarehouseID, productID, quantity)
}

func (c *StockController) ReserveProductStock(warehouseID, productID uint, quantity int) error {
	return c.StockRepo.ReserveStock(warehouseID, productID, quantity)
}

func (c *StockController) ReleaseProductStock(warehouseID, productID uint, quantity int) error {
	return c.StockRepo.ReleaseReservedStock(warehouseID, productID, quantity)
}

func (c *StockController) GetProductStockWithDetail(productID uint) (int, []model.StockWithProductDTO, error) {
	stocks, err := c.StockRepo.FindStockByProductID(productID)
	if err != nil {
		return 0, nil, err
	}

	product, err := fetchProduct(productID)
	if err != nil {
		return 0, nil, err
	}

	total := 0
	var result []model.StockWithProductDTO
	for _, stock := range stocks {
		warehouse, err := c.WarehouseRepo.FindByID(stock.WarehouseID)
		if err != nil || !warehouse.Active {
			continue
		}
		available := stock.Quantity - stock.ReservedQuantity
		total += available
		result = append(result, model.StockWithProductDTO{
			WarehouseID: warehouse.ID,
			Quantity:    available,
			Product:     *product,
		})
	}

	return total, result, nil
}

func fetchProduct(productID uint) (*model.ProductDTO, error) {
	url := os.Getenv("PRODUCT_SERVICE_URL") + "/products/" + strconv.Itoa(int(productID))
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch product")
	}
	defer resp.Body.Close()

	var product model.ProductDTO
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}
