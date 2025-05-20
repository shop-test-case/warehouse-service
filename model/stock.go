package model

type WarehouseStock struct {
	ID               uint `gorm:"primaryKey" json:"id"`
	WarehouseID      uint `json:"warehouse_id"`
	ProductID        uint `json:"product_id"`
	Quantity         int  `json:"quantity"`
	ReservedQuantity int  `json:"reserved_quantity"`
}
