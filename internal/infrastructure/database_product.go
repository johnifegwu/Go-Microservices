package database

import (
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/johnifegwu/go-microservices/internal/dberrors"
	"github.com/johnifegwu/go-microservices/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c Client) SearchProducts(ctx context.Context, searchterm string, pageindex string, pagesize string) ([]models.Product, error) {
	// Convert pageindex and pagesize to integers with proper error handling
	pIndex, err := strconv.Atoi(pageindex)
	if err != nil || pIndex < 1 {
		pIndex = 1 // default to first page
	}

	pSize, err := strconv.Atoi(pagesize)
	if err != nil || pSize < 1 {
		pSize = 10 // default page size
	} else if pSize > 100 {
		pSize = 100 // maximum page size
	}

	// Construct the search pattern to match any part of the name
	searchPattern := "%" + searchterm + "%"

	// Calculate the offset for pagination
	offset := (pIndex - 1) * pSize

	var products []models.Product

	// Perform the database query with the LIKE pattern and pagination
	result := c.DB.WithContext(ctx).
		Where("name LIKE ?", searchPattern).
		Limit(pSize).
		Offset(offset).
		Find(&products)

	// Return the products and any error from the query
	return products, result.Error
}

func (c Client) GetAllProducts(ctx context.Context, pageIndex string, pageSize string) ([]models.Product, error) {
	// Default values for page and pageSize
	page, err := strconv.Atoi(pageIndex)
	if err != nil || page < 1 {
		page = 1
	}

	pSize, err := strconv.Atoi(pageSize)
	if err != nil || pSize < 1 {
		pSize = 10
	} else if pSize > 100 {
		pSize = 100
	}

	// Calculate the offset for SQL query
	offset := (page - 1) * pSize

	// Query the products table with LIMIT and OFFSET for pagination
	var products []models.Product
	result := c.DB.WithContext(ctx).Limit(pSize).Offset(offset).Order("name").Find(&products)
	return products, result.Error
}

func (c Client) GetProductById(ctx context.Context, productId string) (*models.Product, error) {
	// Parse the string into a uuid.UUID
	parsedUUID, err := uuid.Parse(productId)

	if err != nil {
		return nil, err
	}

	// Query the product by product_id
	product := &models.Product{}
	result := c.DB.WithContext(ctx).Where(models.Product{ProductID: parsedUUID}).First(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "product", ID: parsedUUID}
		}
		return nil, result.Error
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
	} else if pSize > 100 {
		pSize = 100
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

func (c Client) AddProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	product.ProductID = uuid.Must(uuid.NewRandom())

	//Create product
	result := c.DB.WithContext(ctx).Create(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
	}

	return product, result.Error
}

func (c Client) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	// Update product
	result := c.DB.WithContext(ctx).
		Clauses(clause.Returning{}).
		Save(&product)

	if result.Error != nil {
		return nil, result.Error
	}

	return product, result.Error
}

func (c Client) DeleteProduct(ctx context.Context, productId string) (int64, error) {
	// Parse the string into a uuid.UUID
	parsedUUID, errUUID := uuid.Parse(productId)

	if errUUID != nil {
		return 0, errUUID
	}
	// Update product
	result := c.DB.WithContext(ctx).Delete(&models.Product{}, parsedUUID)

	return result.RowsAffected, result.Error
}
