package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/johnifegwu/go-microservices/internal/dberrors"
	"github.com/johnifegwu/go-microservices/internal/models"
	server "github.com/johnifegwu/go-microservices/internal/server/models"
	"github.com/labstack/echo/v4"
)

// GetAllServices godoc
// @Summary Get all services
// @Description Get a list of all services with optional pagination
// @Tags services
// @Accept  json
// @Produce  json
// @Param pageindex query string false "Page index for pagination"
// @Param pagesize query string false "Page size for pagination"
// @Success 200 {array} models.Service
// @Router /services [get]
func (s *EchoServer) GetAllServices(ctx echo.Context) error {
	pageindex := ctx.QueryParam("pageindex")
	pagesize := ctx.QueryParam("pagesize")
	services, err := s.DB.GetAllServices(ctx.Request().Context(), pageindex, pagesize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, services)
}

// GetServiceById godoc
// @Summary Get service by ID
// @Description Get a single service by its ID
// @Tags services
// @Accept  json
// @Produce  json
// @Param id path string true "Service ID"
// @Success 200 {object} models.Service
// @Router /services/{id} [get]
func (s *EchoServer) GetServiceById(ctx echo.Context) error {
	id := ctx.Param("id")
	service, err := s.DB.GetServiceById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, service)
}

// AddService godoc
// @Summary Add a new service
// @Description Add a new service to the database
// @Tags services
// @Accept  json
// @Produce  json
// @Param service body models.Service true "Service to add"
// @Success 201 {object} models.Service
// @Failure 409 {object} dberrors.ConflictError
// @Failure 415 {object} string "Unsupported Media Type"
// @Router /services [post]
func (s *EchoServer) AddService(ctx echo.Context) error {
	service := new(models.Service)

	if err := ctx.Bind(service); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	service, err := s.DB.AddService(ctx.Request().Context(), service)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusCreated, service)
}

// UpdateService godoc
// @Summary Update an existing service
// @Description Update a service's details by providing its ID
// @Tags services
// @Accept  json
// @Produce  json
// @Param service_id path string true "Service ID"
// @Param service body models.Service true "Updated service data"
// @Success 201 {object} models.Service
// @Failure 400 {object} string "Bad Request"
// @Failure 415 {object} string "Unsupported Media Type"
// @Router /services/{service_id} [put]
func (s *EchoServer) UpdateService(ctx echo.Context) error {
	service := new(models.Service)

	if err := ctx.Bind(service); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	ID, errUUID := uuid.Parse(ctx.Param("service_id"))
	if errUUID == nil {
		if ID != service.ServiceID {
			return ctx.JSON(http.StatusBadRequest, "id on path doesn't match id on body")
		}
	}

	service, err := s.DB.UpdateService(ctx.Request().Context(), service)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, service)
}

// DeleteService godoc
// @Summary Delete a service
// @Description Delete a service from the database by its ID
// @Tags services
// @Accept  json
// @Produce  json
// @Param id query string true "Service ID"
// @Success 200 {object} server.Response
// @Failure 500 {object} dberrors.ZeroRowsAffectedError
// @Router /services [delete]
func (s *EchoServer) DeleteService(ctx echo.Context) error {
	var serviceId = ctx.QueryParam("id")

	rowsaffected, err := s.DB.DeleteService(ctx.Request().Context(), serviceId)

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
