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

func (c Client) GetVendorById(ctx context.Context, vendorId string) (*models.Vendor, error) {
	// Parse the string into a uuid.UUID
	parsedUUID, err := uuid.Parse(vendorId)

	if err != nil {
		return nil, err
	}

	// Query the Vendor by id
	vendor := &models.Vendor{}
	result := c.DB.WithContext(ctx).Where(models.Vendor{VendorID: parsedUUID}).First(&vendor)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "vendor", ID: parsedUUID}
		}
		return nil, result.Error
	}

	return vendor, result.Error
}

func (c Client) AddVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	vendor.VendorID = uuid.Must(uuid.NewRandom())

	//Create Vendor
	result := c.DB.WithContext(ctx).Create(&vendor)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
	}

	return vendor, result.Error
}

func (c Client) UpdateVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	// Update Vendor
	result := c.DB.WithContext(ctx).
		Clauses(clause.Returning{}).
		Save(&vendor)

	if result.Error != nil {
		return nil, result.Error
	}

	return vendor, result.Error
}

func (c Client) DeleteVendor(ctx context.Context, vendorId string) (int64, error) {
	parsedUUID, uuidErr := uuid.Parse(vendorId)

	if uuidErr != nil {
		return 0, uuidErr
	}

	result := c.DB.WithContext(ctx).Delete(&models.Vendor{}, parsedUUID)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, result.Error
}
