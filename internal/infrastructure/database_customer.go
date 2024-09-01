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

func (c Client) GetAllCustomers(ctx context.Context, email string, pageIndex string, pageSize string) ([]models.Customer, error) {
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

	var customers []models.Customer
	result := c.DB.WithContext(ctx).
		Where(models.Customer{Email: email}).
		Limit(pSize).Offset(offset).Order("first_name").
		Find(&customers)
	return customers, result.Error
}

func (c Client) GetCustomerById(ctx context.Context, customerId string) (*models.Customer, error) {
	// Parse the string into a uuid.UUID
	parsedUUID, err := uuid.Parse(customerId)

	if err != nil {
		return nil, err
	}

	// Query the Customer by id
	customer := &models.Customer{}
	result := c.DB.WithContext(ctx).Where(models.Customer{CustomerID: parsedUUID}).First(&customer)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "customer", ID: parsedUUID}
		}
		return nil, result.Error
	}

	return customer, result.Error
}

func (c Client) AddCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	customer.CustomerID = uuid.Must(uuid.NewRandom())
	result := c.DB.WithContext(ctx).
		Create(&customer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}
	return customer, nil
}

func (c Client) UpdateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {

	result := c.DB.WithContext(ctx).
		Clauses(clause.Returning{}).
		Save(&customer)

	if result.Error != nil {
		return nil, result.Error
	}

	return customer, nil
}

func (c Client) DeleteCustomer(ctx context.Context, customerId string) (int64, error) {
	parsedUUID, uuidErr := uuid.Parse(customerId)

	if uuidErr != nil {
		return 0, uuidErr
	}

	result := c.DB.WithContext(ctx).Delete(&models.Customer{}, parsedUUID)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, result.Error
}
