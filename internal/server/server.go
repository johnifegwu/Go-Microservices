package server

import (
	"log"
	"net/http"
	"os"

	database "github.com/johnifegwu/go-microservices/internal/infrastructure"
	"github.com/johnifegwu/go-microservices/internal/models"
	"github.com/labstack/echo/v4"
)

type Server interface {
	Start() error
	Readiness(ctx echo.Context) error
	Liveness(ctx echo.Context) error

	GetAllCustomers(ctx echo.Context) error

	GetAllProducts(ctx echo.Context) error

	GetProductById(ctx echo.Context) error

	GetAllProductsByVendor(ctx echo.Context) error

	GetAllServices(ctx echo.Context) error

	GetAllVendors(ctx echo.Context) error
}

type EchoServer struct {
	echo *echo.Echo
	DB   database.DatabaseClient
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

	cg := s.echo.Group("/customers")
	cg.GET("", s.GetAllCustomers)

	pg := s.echo.Group("/products")
	pg.GET("", s.GetAllProducts)

	pg2 := s.echo.Group("/productbyid")
	pg2.GET("", s.GetProductById)

	pg3 := s.echo.Group("/productsbyvendor")
	pg3.GET("", s.GetAllProductsByVendor)

	sg := s.echo.Group("/services")
	sg.GET("", s.GetAllServices)

	vg := s.echo.Group("/vendors")
	vg.GET("", s.GetAllVendors)

}

func (s *EchoServer) Readiness(ctx echo.Context) error {
	ready := s.DB.Ready()
	if ready {
		return ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
	}
	return ctx.JSON(http.StatusInternalServerError, models.Health{Status: "Failure"})
}

func (s *EchoServer) Liveness(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, models.Health{Status: "OK"})
}
