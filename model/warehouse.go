package model

import "gorm.io/gorm"

type Warehouse struct {
	gorm.Model
	Name   string `json:"name"`
	Active bool   `json:"active"`
}
