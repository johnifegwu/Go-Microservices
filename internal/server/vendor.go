package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/johnifegwu/go-microservices/internal/dberrors"
	"github.com/johnifegwu/go-microservices/internal/models"
	server "github.com/johnifegwu/go-microservices/internal/server/models"
	"github.com/labstack/echo/v4"
)

// GetAllVendors godoc
// @Summary Get all vendors
// @Description Retrieve a list of all vendors with optional pagination
// @Tags vendors
// @Accept  json
// @Produce  json
// @Param pageindex query string false "Page index for pagination"
// @Param pagesize query string false "Page size for pagination"
// @Success 200 {array} models.Vendor
// @Router /vendors [get]
func (s *EchoServer) GetAllVendors(ctx echo.Context) error {
	pageindex := ctx.QueryParam("pageindex")
	pagesize := ctx.QueryParam("pagesize")
	vendors, err := s.DB.GetAllVendors(ctx.Request().Context(), pageindex, pagesize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, vendors)
}

// GetVendorById godoc
// @Summary Get vendor by ID
// @Description Retrieve a vendor by its ID
// @Tags vendors
// @Accept  json
// @Produce  json
// @Param id path string true "Vendor ID"
// @Success 200 {object} models.Vendor
// @Router /vendors/{id} [get]
func (s *EchoServer) GetVendorById(ctx echo.Context) error {
	id := ctx.Param("id")
	vendor, err := s.DB.GetVendorById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, vendor)
}

// AddVendor godoc
// @Summary Add a new vendor
// @Description Create a new vendor in the database
// @Tags vendors
// @Accept  json
// @Produce  json
// @Param vendor body models.Vendor true "Vendor to add"
// @Success 201 {object} models.Vendor
// @Failure 409 {object} dberrors.ConflictError
// @Failure 415 {object} string "Unsupported Media Type"
// @Router /vendors [post]
func (s *EchoServer) AddVendor(ctx echo.Context) error {
	vendor := new(models.Vendor)

	if err := ctx.Bind(vendor); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	vendor, err := s.DB.AddVendor(ctx.Request().Context(), vendor)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusCreated, vendor)
}

// UpdateVendor godoc
// @Summary Update an existing vendor
// @Description Update a vendor's details by providing its ID
// @Tags vendors
// @Accept  json
// @Produce  json
// @Param vendor_id path string true "Vendor ID"
// @Param vendor body models.Vendor true "Updated vendor data"
// @Success 201 {object} models.Vendor
// @Failure 400 {object} string "Bad Request"
// @Failure 415 {object} string "Unsupported Media Type"
// @Router /vendors/{vendor_id} [put]
func (s *EchoServer) UpdateVendor(ctx echo.Context) error {
	vendor := new(models.Vendor)

	if err := ctx.Bind(vendor); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	ID, errUUID := uuid.Parse(ctx.Param("vendor_id"))
	if errUUID == nil {
		if ID != vendor.VendorID {
			return ctx.JSON(http.StatusBadRequest, "ID on path doesn't match ID in body")
		}
	}

	vendor, err := s.DB.UpdateVendor(ctx.Request().Context(), vendor)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, vendor)
}

// DeleteVendor godoc
// @Summary Delete a vendor
// @Description Delete a vendor from the database by its ID
// @Tags vendors
// @Accept  json
// @Produce  json
// @Param id query string true "Vendor ID"
// @Success 200 {object} server.Response
// @Failure 500 {object} dberrors.ZeroRowsAffectedError
// @Router /vendors [delete]
func (s *EchoServer) DeleteVendor(ctx echo.Context) error {
	var vendorId = ctx.QueryParam("id")

	rowsaffected, err := s.DB.DeleteVendor(ctx.Request().Context(), vendorId)

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
