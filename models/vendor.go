package models

import (
	"github.com/google/uuid"
)

type Vendor struct {
	VendorID uuid.UUID `json:"vendor_id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `json:"name" gorm:"not null"`
	Contact  string    `json:"contact"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Address  string    `json:"address"`
}

// TableName sets the table name for Vendor
func (Vendor) TableName() string {
	return "wisdom.vendors"
}
