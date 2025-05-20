package model

type WarehouseStockDTO struct {
	WarehouseID uint `json:"warehouse_id"`
	Quantity    int  `json:"quantity"`
}

type StockWithProductDTO struct {
	WarehouseID uint       `json:"warehouse_id"`
	Quantity    int        `json:"quantity"`
	Product     ProductDTO `json:"product"`
}
