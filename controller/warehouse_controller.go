package controller

import (
	"warehouse-service/model"
	"warehouse-service/repository"
)

type WarehouseController struct {
	WarehouseRepo repository.IWarehouseRepo
}

func (c *WarehouseController) AddWarehouse(warehouse *model.Warehouse) error {
	return c.WarehouseRepo.Create(warehouse)
}

func (c *WarehouseController) ToggleActive(id uint, active bool) error {
	return c.WarehouseRepo.SetActive(id, active)
}

func (c *WarehouseController) List() ([]model.Warehouse, error) {
	return c.WarehouseRepo.FindAll()
}
