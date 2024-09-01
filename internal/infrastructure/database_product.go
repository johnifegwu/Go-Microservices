package database

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/johnifegwu/go-microservices/internal/models"
)

func (c Client) GetAllProducts(ctx context.Context, pageIndex string, pageSize string) ([]models.Product, error) {
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
	var products []models.Product
	result := c.DB.WithContext(ctx).Limit(pSize).Offset(offset).Order("name").Find(&products)
	return products, result.Error
}

func (c Client) GetProductById(ctx context.Context, productId string) (models.Product, error) {
	// Parse the string into a uuid.UUID
	parsedUUID, err := uuid.Parse(productId)

	// Query the product by product_id
	var product models.Product
	result := c.DB.WithContext(ctx).Where(models.Product{ProductID: parsedUUID}).First(&product)

	if err != nil {
		return product, err
	}

	return product, result.Error
}

func (c Client) GetAllProductsByVendor(ctx context.Context, vendorID string, pageIndex string, pageSize string) ([]models.Product, error) {
	// Parse the string into a uuid.UUID
	parsedUUID, errUUID := uuid.Parse(vendorID)

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
	var products []models.Product
	result := c.DB.WithContext(ctx).Where("vendor_id = ?", parsedUUID).Limit(pSize).Offset(offset).Order("name").Find(&products)

	if errUUID != nil {
		return products, errUUID
	}

	return products, result.Error
}
