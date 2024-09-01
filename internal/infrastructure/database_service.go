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

func (c Client) GetAllServices(ctx context.Context, pageIndex string, pageSize string) ([]models.Service, error) {
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
	var services []models.Service
	result := c.DB.WithContext(ctx).Limit(pSize).Offset(offset).Order("name").Find(&services)
	return services, result.Error
}

func (c Client) GetServiceById(ctx context.Context, serviceId string) (*models.Service, error) {
	// Parse the string into a uuid.UUID
	parsedUUID, err := uuid.Parse(serviceId)

	if err != nil {
		return nil, err
	}

	// Query the Customer by id
	service := &models.Service{}
	result := c.DB.WithContext(ctx).Where(models.Service{ServiceID: parsedUUID}).First(&service)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "service", ID: parsedUUID}
		}
		return nil, result.Error
	}

	return service, result.Error
}

func (c Client) AddService(ctx context.Context, service *models.Service) (*models.Service, error) {
	service.ServiceID = uuid.Must(uuid.NewRandom())

	//Create product
	result := c.DB.WithContext(ctx).Create(&service)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
	}

	return service, result.Error
}

func (c Client) UpdateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	// Update product
	result := c.DB.WithContext(ctx).
		Clauses(clause.Returning{}).
		Save(&service)

	if result.Error != nil {
		return nil, result.Error
	}

	return service, result.Error
}

func (c Client) DeleteService(ctx context.Context, serviceid string) (int64, error) {
	parsedUUID, uuidErr := uuid.Parse(serviceid)

	if uuidErr != nil {
		return 0, uuidErr
	}

	result := c.DB.WithContext(ctx).Delete(&models.Service{}, parsedUUID)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, result.Error
}
