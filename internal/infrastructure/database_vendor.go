package database

import (
	"context"
	"strconv"

	"github.com/johnifegwu/go-microservices/internal/models"
)

func (c Client) GetAllVendors(ctx context.Context, pageIndex string, pageSize string) ([]models.Vendor, error) {
	// Default values for page and pageSize
	page, err := strconv.Atoi(pageIndex)
	if err != nil || page < 1 {
		page = 1
	}

	pSize, err := strconv.Atoi(pageSize)
	if err != nil || pSize < 1 {
		pSize = 10
	}

	// Calculate the offset for SQL query
	offset := (page - 1) * pSize

	// Query the products table with LIMIT and OFFSET for pagination
	var vendors []models.Vendor
	result := c.DB.WithContext(ctx).Limit(pSize).Offset(offset).Order("name").Find(&vendors)
	return vendors, result.Error
}
