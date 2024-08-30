package models

import (
	"github.com/google/uuid"
)

type Service struct {
	ServiceID uuid.UUID `json:"service_id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	Price     float64   `json:"price" gorm:"type:numeric(12,2)"`
}

// TableName sets the table name for Service
func (Service) TableName() string {
	return "wisdom.services"
}
