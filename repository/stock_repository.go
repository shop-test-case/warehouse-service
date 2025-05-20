package repository

import (
	"errors"
	"warehouse-service/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IStockRepo interface {
	FindStockByProductID(productID uint) ([]model.WarehouseStock, error)
	TransferStock(fromWarehouseID, toWarehouseID, productID uint, quantity int) error
	ReserveStock(warehouseID, productID uint, quantity int) error
	ReleaseReservedStock(warehouseID, productID uint, quantity int) error
}

type StockRepo struct {
	DB *gorm.DB
}

func (r *StockRepo) FindStockByProductID(productID uint) ([]model.WarehouseStock, error) {
	var stocks []model.WarehouseStock
	err := r.DB.Where("product_id = ?", productID).Find(&stocks).Error
	return stocks, err
}

func (r *StockRepo) TransferStock(fromWarehouseID, toWarehouseID, productID uint, quantity int) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var from model.WarehouseStock
		err := tx.Where("warehouse_id = ? AND product_id = ?", fromWarehouseID, productID).First(&from).Error
		if err != nil || from.Quantity-from.ReservedQuantity < quantity {
			return errors.New("insufficient available stock")
		}

		from.Quantity -= quantity
		if err := tx.Save(&from).Error; err != nil {
			return err
		}

		var to model.WarehouseStock
		err = tx.Where("warehouse_id = ? AND product_id = ?", toWarehouseID, productID).First(&to).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				to = model.WarehouseStock{WarehouseID: toWarehouseID, ProductID: productID, Quantity: quantity}
				return tx.Create(&to).Error
			}
			return err
		}
		to.Quantity += quantity

		return tx.Save(&to).Error
	})
}

func (r *StockRepo) ReserveStock(warehouseID, productID uint, quantity int) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var stock model.WarehouseStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).
			First(&stock).Error; err != nil {
			return err
		}

		if stock.Quantity-stock.ReservedQuantity < quantity {
			return errors.New("not enough available stock")
		}
		stock.ReservedQuantity += quantity

		return tx.Save(&stock).Error
	})
}

func (r *StockRepo) ReleaseReservedStock(warehouseID, productID uint, quantity int) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var stock model.WarehouseStock
		if err := tx.Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).First(&stock).Error; err != nil {
			return err
		}

		if stock.ReservedQuantity < quantity {
			return errors.New("reserved quantity too low")
		}
		stock.ReservedQuantity -= quantity

		return tx.Save(&stock).Error
	})
}
