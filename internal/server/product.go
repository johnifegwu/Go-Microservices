package server

import (
	"net/http"

	"github.com/johnifegwu/go-microservices/internal/dberrors"
	"github.com/johnifegwu/go-microservices/internal/models"
	server "github.com/johnifegwu/go-microservices/internal/server/models"
	"github.com/labstack/echo/v4"
)

func (s *EchoServer) SearchProducts(ctx echo.Context) error {
	searchterm := ctx.Param("searchterm")
	pageindex := ctx.QueryParam("pageindex")
	pagesize := ctx.QueryParam("pagesize")

	products, err := s.DB.SearchProducts(ctx.Request().Context(), searchterm, pageindex, pagesize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, products)
}

func (s *EchoServer) GetAllProducts(ctx echo.Context) error {
	pageindex := ctx.QueryParam("pageindex")
	pagesize := ctx.QueryParam("pagesize")
	products, err := s.DB.GetAllProducts(ctx.Request().Context(), pageindex, pagesize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, products)
}

func (s *EchoServer) GetProductById(ctx echo.Context) error {
	id := ctx.Param("id")
	products, err := s.DB.GetProductById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, products)
}

func (s *EchoServer) GetAllProductsByVendor(ctx echo.Context) error {
	vendorid := ctx.Param("id")
	pageindex := ctx.QueryParam("pageindex")
	pagesize := ctx.QueryParam("pagesize")
	products, err := s.DB.GetAllProductsByVendor(ctx.Request().Context(), vendorid, pageindex, pagesize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, products)
}

func (s *EchoServer) AddProduct(ctx echo.Context) error {
	product := new(models.Product)

	if err := ctx.Bind(product); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	product, err := s.DB.AddProduct(ctx.Request().Context(), product)

	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, product)
}

func (s *EchoServer) UpdateProduct(ctx echo.Context) error {
	product := new(models.Product)

	if err := ctx.Bind(product); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	product, err := s.DB.UpdateProduct(ctx.Request().Context(), product)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, product)
}

func (s *EchoServer) DeleteProduct(ctx echo.Context) error {
	var productId = ctx.QueryParam("id")

	rowsaffected, err := s.DB.DeleteProduct(ctx.Request().Context(), productId)

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
