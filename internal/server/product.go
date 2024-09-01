package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/johnifegwu/go-microservices/internal/dberrors"
	"github.com/johnifegwu/go-microservices/internal/models"
	server "github.com/johnifegwu/go-microservices/internal/server/models"
	"github.com/labstack/echo/v4"
)

// SearchProducts godoc
// @Summary Search products
// @Description Search for products by a search term with optional pagination
// @Tags products
// @Accept  json
// @Produce  json
// @Param searchterm path string true "Search term"
// @Param pageindex query string false "Page index for pagination"
// @Param pagesize query string false "Page size for pagination"
// @Success 200 {array} models.Product
// @Router /products/search/{searchterm} [get]
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

// GetAllProducts godoc
// @Summary Get all products
// @Description Get a list of all products with optional pagination
// @Tags products
// @Accept  json
// @Produce  json
// @Param pageindex query string false "Page index for pagination"
// @Param pagesize query string false "Page size for pagination"
// @Success 200 {array} models.Product
// @Router /products [get]
func (s *EchoServer) GetAllProducts(ctx echo.Context) error {
	pageindex := ctx.QueryParam("pageindex")
	pagesize := ctx.QueryParam("pagesize")
	products, err := s.DB.GetAllProducts(ctx.Request().Context(), pageindex, pagesize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, products)
}

// GetProductById godoc
// @Summary Get product by ID
// @Description Get a single product by its ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Router /products/{id} [get]
func (s *EchoServer) GetProductById(ctx echo.Context) error {
	id := ctx.Param("id")
	products, err := s.DB.GetProductById(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, products)
}

// GetAllProductsByVendor godoc
// @Summary Get all products by vendor
// @Description Get a list of all products by a specific vendor with optional pagination
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Vendor ID"
// @Param pageindex query string false "Page index for pagination"
// @Param pagesize query string false "Page size for pagination"
// @Success 200 {array} models.Product
// @Router /products/vendor/{id} [get]
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

// AddProduct godoc
// @Summary Add a new product
// @Description Add a new product to the database
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body models.Product true "Product to add"
// @Success 201 {object} models.Product
// @Failure 409 {object} dberrors.ConflictError
// @Failure 415 {object} string "Unsupported Media Type"
// @Router /products [post]
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

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update a product's details by providing its ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param product_id path string true "Product ID"
// @Param product body models.Product true "Updated product data"
// @Success 201 {object} models.Product
// @Failure 400 {object} string "Bad Request"
// @Failure 415 {object} string "Unsupported Media Type"
// @Router /products/{product_id} [put]
func (s *EchoServer) UpdateProduct(ctx echo.Context) error {
	product := new(models.Product)

	if err := ctx.Bind(product); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}

	ID, errUUID := uuid.Parse(ctx.Param("product_id"))
	if errUUID == nil {
		if ID != product.ProductID {
			return ctx.JSON(http.StatusBadRequest, "id on path doesn't match id on body")
		}
	}

	product, err := s.DB.UpdateProduct(ctx.Request().Context(), product)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product from the database by its ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id query string true "Product ID"
// @Success 200 {object} server.Response
// @Failure 500 {object} dberrors.ZeroRowsAffectedError
// @Router /products [delete]
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
