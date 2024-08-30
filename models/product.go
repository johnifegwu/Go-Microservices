package models

import (
	"github.com/google/uuid"
)

type Product struct {
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	Price     float64   `json:"price" gorm:"type:numeric(12,2)"`
	VendorID  uuid.UUID `json:"vendor_id" gorm:"type:uuid;not null"`
	Vendor    Vendor    `gorm:"foreignKey:VendorID;references:VendorID"`
}

// TableName sets the table name for Product
func (Product) TableName() string {
	return "wisdom.products"
}
