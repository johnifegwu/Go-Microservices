package server

import (
	"log"
	"net/http"
	"os"

	_ "github.com/johnifegwu/go-microservices/docs" // Import the generated docs
	database "github.com/johnifegwu/go-microservices/internal/infrastructure"
	"github.com/johnifegwu/go-microservices/internal/models"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger" // Import echo-swagger middleware
)

type Server interface {
	Start() error
	Readiness(ctx echo.Context) error
	Liveness(ctx echo.Context) error

	GetAllCustomers(ctx echo.Context) error
	GetCustomerById(ctx echo.Context) error
	AddCustomer(ctx echo.Context) error
	UpdateCustomer(ctx echo.Context) error
	DeleteCustomer(ctx echo.Context) error

	GetAllProducts(ctx echo.Context) error
	SearchProducts(ctx echo.Context) error
	GetProductById(ctx echo.Context) error
	GetAllProductsByVendor(ctx echo.Context) error
	AddProduct(ctx echo.Context) error
	UpdateProduct(ctx echo.Context) error
	DeleteProduct(ctx echo.Context) error

	GetAllServices(ctx echo.Context) error
	GetServiceById(ctx echo.Context) error
	AddService(ctx echo.Context) error
	UpdateService(ctx echo.Context) error
	DeleteService(ctx echo.Context) error

	GetAllVendors(ctx echo.Context) error
	GetVendorById(ctx echo.Context) error
	AddVendor(ctx echo.Context) error
	UpdateVendor(ctx echo.Context) error
	DeleteVendor(ctx echo.Context) error
}

// @title Echo Server API
// @version 1.0
// @description This is a sample server for managing customers, products, services, and vendors.
// @host localhost:8080
// @BasePath /

// EchoServer represents the server
type EchoServer struct {
	echo *echo.Echo
	DB   database.DatabaseClient
}

// @Summary Readiness probe
// @Description Check if the service is ready to accept requests
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Health
// @Failure 500 {object} models.Health
// @Router /readiness [get]
func (s *EchoServer) Readiness(ctx echo.Context) error {
	ready := s.DB.Ready()
	if ready {
		return ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
	}
	return ctx.JSON(http.StatusInternalServerError, models.Health{Status: "Failure"})
}

// @Summary Liveness probe
// @Description Check if the service is alive
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Health
// @Router /liveness [get]
func (s *EchoServer) Liveness(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
}

func NewEchoServer(db database.DatabaseClient) Server {
	server := &EchoServer{
		echo: echo.New(),
		DB:   db,
	}
	server.registerRoutes()
	return server
}

func (s *EchoServer) Start() error {
	// Retrieve environment variables
	port := ":" + os.Getenv("DEFAULTPORT")

	if err := s.echo.Start(port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server shutdown occurred: %s", err)
		return err
	}
	return nil
}

func (s *EchoServer) registerRoutes() {
	s.echo.GET("/readiness", s.Readiness)
	s.echo.GET("/liveness", s.Liveness)

	// Serve the Swagger documentation
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	cg := s.echo.Group("/customers")
	cg.GET("", s.GetAllCustomers)
	cg.GET("/:id", s.GetCustomerById)
	cg.POST("", s.AddCustomer)
	cg.PUT("", s.UpdateCustomer)
	cg.DELETE("", s.DeleteCustomer)

	pg := s.echo.Group("/products")
	pg.GET("", s.GetAllProducts)
	pg.GET("/search/:searchterm", s.SearchProducts)
	pg.GET("/vendor/:id", s.GetAllProductsByVendor)
	pg.GET("/:id", s.GetProductById)
	pg.POST("", s.AddProduct)
	pg.PUT("", s.UpdateProduct)
	pg.DELETE("", s.DeleteProduct)

	sg := s.echo.Group("/services")
	sg.GET("", s.GetAllServices)
	sg.GET("/:id", s.GetServiceById)
	sg.POST("", s.AddService)
	sg.PUT("", s.UpdateService)
	sg.DELETE("", s.DeleteService)

	vg := s.echo.Group("/vendors")
	vg.GET("", s.GetAllVendors)
	vg.GET("/:id", s.GetVendorById)
	vg.POST("", s.AddVendor)
	vg.PUT("", s.UpdateVendor)
	vg.DELETE("", s.DeleteVendor)
}
