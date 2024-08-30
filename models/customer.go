package models

import (
	"github.com/google/uuid"
)

type Customer struct {
	CustomerID uuid.UUID `json:"customer_id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Address    string    `json:"address"`
}

// TableName sets the table name for Customer
func (Customer) TableName() string {
	return "wisdom.customers"
}
