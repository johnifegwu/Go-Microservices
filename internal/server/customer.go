package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/johnifegwu/go-microservices/internal/dberrors"
	"github.com/johnifegwu/go-microservices/internal/models"
	server "github.com/johnifegwu/go-microservices/internal/server/models"
	"github.com/labstack/echo/v4"
)

// GetAllCustomers godoc
// @Summary Get all customers
// @Description Get all customers with optional filtering by email, and pagination
// @Tags customers
// @Accept  json
// @Produce  json
// @Param email query string false "Email address for filtering"
// @Param pageindex query string false "Page index for pagination"
// @Param pagesize query string false "Page size for pagination"
// @Success 200 {array} models.Customer
// @Router /customers [get]
func (s *EchoServer) GetAllCustomers(ctx echo.Context) error {
	email := ctx.QueryParam("email")
	pageindex := ctx.QueryParam("pageindex")
	pagesize := ctx.QueryParam("pagesize")

	customers, err := s.DB.GetAllCustomers(ctx.Request().Context(), email, pageindex, pagesize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, customers)
}

// GetCustomerById godoc
// @Summary Get customer by ID
// @Description Get a single customer by its ID
// @Tags customers
// @Accept  json
// @Produce  json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.Customer
// @Router /customers/{id} [get]
func (s *EchoServer) GetCustomerById(ctx echo.Context) error {
	id := ctx.Param("id")
	products, err := s.DB.GetCustomerById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, products)
}

// AddCustomer godoc
// @Summary Add a new customer
// @Description Add a new customer to the database
// @Tags customers
// @Accept  json
// @Produce  json
// @Param customer body models.Customer true "Customer to add"
// @Success 201 {object} models.Customer
// @Failure 409 {object} dberrors.ConflictError
// @Failure 415 {object} string "Unsupported Media Type"
// @Router /customers [post]
func (s *EchoServer) AddCustomer(ctx echo.Context) error {
	customer := new(models.Customer)

	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	customer, err := s.DB.AddCustomer(ctx.Request().Context(), customer)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusCreated, customer)
}

// UpdateCustomer godoc
// @Summary Update an existing customer
// @Description Update a customer's details by providing its ID
// @Tags customers
// @Accept  json
// @Produce  json
// @Param customer_id path string true "Customer ID"
// @Param customer body models.Customer true "Updated customer data"
// @Success 201 {object} models.Customer
// @Failure 400 {object} string "Bad Request"
// @Failure 415 {object} string "Unsupported Media Type"
// @Router /customers/{customer_id} [put]
func (s *EchoServer) UpdateCustomer(ctx echo.Context) error {
	customer := new(models.Customer)

	if err := ctx.Bind(customer); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	ID, errUUID := uuid.Parse(ctx.Param("customer_id"))
	if errUUID == nil {
		if ID != customer.CustomerID {
			return ctx.JSON(http.StatusBadRequest, "id on path doesn't match id on body")
		}
	}

	customer, err := s.DB.UpdateCustomer(ctx.Request().Context(), customer)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, customer)
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Delete a customer from the database by its ID
// @Tags customers
// @Accept  json
// @Produce  json
// @Param id query string true "Customer ID"
// @Success 200 {object} server.Response
// @Failure 500 {object} dberrors.ZeroRowsAffectedError
// @Router /customers [delete]
func (s *EchoServer) DeleteCustomer(ctx echo.Context) error {
	var customerId = ctx.QueryParam("id")

	rowsaffected, err := s.DB.DeleteCustomer(ctx.Request().Context(), customerId)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	if rowsaffected < 1 {
		dberr := new(dberrors.ZeroRowsAffectedError)
		return ctx.JSON(http.StatusInternalServerError, dberr.Error())
	}

	response := server.Response{
		Status:  "Ok",
		Message: "Record deleted successfully",
	}
	return ctx.JSON(http.StatusOK, response)
}
