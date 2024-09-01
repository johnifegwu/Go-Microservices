package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/johnifegwu/go-microservices/internal/dberrors"
	"github.com/johnifegwu/go-microservices/internal/models"
	server "github.com/johnifegwu/go-microservices/internal/server/models"
	"github.com/labstack/echo/v4"
)

func (s *EchoServer) GetAllVendors(ctx echo.Context) error {
	pageindex := ctx.QueryParam("pageindex")
	pagesize := ctx.QueryParam("pagesize")
	vendors, err := s.DB.GetAllVendors(ctx.Request().Context(), pageindex, pagesize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, vendors)
}

func (s *EchoServer) GetVendorById(ctx echo.Context) error {
	id := ctx.Param("id")
	vendor, err := s.DB.GetVendorById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, vendor)
}

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

func (s *EchoServer) UpdateVendor(ctx echo.Context) error {
	vendor := new(models.Vendor)

	if err := ctx.Bind(vendor); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	ID, errUUID := uuid.Parse(ctx.Param("vendor_id"))
	if errUUID == nil {
		if ID != vendor.VendorID {
			return ctx.JSON(http.StatusBadRequest, "id on path doesn't match id on body")
		}
	}

	service, err := s.DB.UpdateVendor(ctx.Request().Context(), vendor)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, service)
}

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
