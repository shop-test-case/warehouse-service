package repository

import (
	"warehouse-service/model"

	"gorm.io/gorm"
)

type IWarehouseRepo interface {
	Create(warehouse *model.Warehouse) error
	SetActive(id uint, active bool) error
	FindAll() ([]model.Warehouse, error)
	FindByID(id uint) (*model.Warehouse, error)
}

type WarehouseRepo struct {
	DB *gorm.DB
}

func (r *WarehouseRepo) Create(warehouse *model.Warehouse) error {
	return r.DB.Create(warehouse).Error
}

func (r *WarehouseRepo) SetActive(id uint, active bool) error {
	return r.DB.Model(&model.Warehouse{}).Where("id = ?", id).Update("active", active).Error
}

func (r *WarehouseRepo) FindAll() ([]model.Warehouse, error) {
	var warehouses []model.Warehouse

	err := r.DB.Find(&warehouses).Error

	return warehouses, err
}

func (r *WarehouseRepo) FindByID(id uint) (*model.Warehouse, error) {
	var warehouse model.Warehouse

	err := r.DB.First(&warehouse, id).Error

	return &warehouse, err
}
