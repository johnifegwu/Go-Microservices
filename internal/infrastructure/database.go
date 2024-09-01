package database

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/johnifegwu/go-microservices/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseClient interface {
	Ready() bool

	SearchProducts(ctx context.Context, searchterm string, pageindex string, pagesize string) ([]models.Product, error)

	GetAllProducts(ctx context.Context, pageIndex string, pageSize string) ([]models.Product, error)

	GetProductById(ctx context.Context, productId string) (models.Product, error)

	GetAllProductsByVendor(ctx context.Context, vendorID string, pageIndex string, pageSize string) ([]models.Product, error)

	GetAllCustomers(ctx context.Context, email, pageindex, pagesize string) ([]models.Customer, error)

	AddProduct(ctx context.Context, product *models.Product) (*models.Product, error)

	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)

	DeleteProduct(ctx context.Context, productId string) (int64, error)

	AddCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error)

	UpdateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error)

	DeleteCustomer(ctx context.Context, customerId string) (int64, error)

	GetAllServices(ctx context.Context, pageIndex string, pageSize string) ([]models.Service, error)

	GetAllVendors(ctx context.Context, pageIndex string, pageSize string) ([]models.Vendor, error)
}

type Client struct {
	DB *gorm.DB
}

//go:embed data.sql
var fdat embed.FS

func NewDatabaseClient() (DatabaseClient, error) {

	// Retrieve environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Create the connection string using environment variables
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// NamingStrategy: schema.NamingStrategy{
		// 	TablePrefix: "wisdom.",
		// },
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		QueryFields: true,
	})

	if err != nil {
		return nil, err
	}

	// Configure the connection pool
	SqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get the underlying database object:", err)
	}
	SqlDB.SetMaxOpenConns(1000)
	SqlDB.SetMaxIdleConns(10)
	SqlDB.SetConnMaxLifetime(time.Minute * 5)

	client := Client{
		DB: db,
	}

	// Migrate the schema and seed data
	MigrateAndSeedData(client.DB)

	return client, nil
}

func (c Client) Ready() bool {
	var ready string
	tx := c.DB.Raw("SELECT 1 as ready").Scan(&ready)
	if tx.Error != nil {
		return false
	}
	if ready == "1" {
		return true
	}
	return false
}

func MigrateAndSeedData(DB *gorm.DB) {
	// Create the wisdom schema
	var dbErr = DB.Exec("CREATE SCHEMA IF NOT EXISTS wisdom;").Error
	if dbErr != nil {
		log.Fatal("Failed to create wisdom schema:", dbErr)
	}

	// Enable the pgcrypto extension
	dbErr = DB.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error
	if dbErr != nil {
		log.Fatal("Failed to enable pgcrypto extension:", dbErr)
	}

	// Automatically migrate the schema
	err := DB.AutoMigrate(&models.Customer{}, &models.Vendor{}, &models.Product{}, &models.Service{})

	if err != nil {
		log.Fatalf("Failed to migrate schema: %s", err)
	}

	// Check if there is any product data
	var product models.Product
	result := DB.First(&product)
	if result.Error == nil && result.RowsAffected > 0 {
		return // Data already exists
	}

	// Seed data if needed (example logic)
	// Read the Data SQL file and execute the commands

	dataSqlBytes, err := fdat.ReadFile("data.sql")

	if err != nil {
		log.Fatalf("Failed to read SQL file: %s", err)
	}

	dataSqlQuery := string(dataSqlBytes)
	err = DB.Exec(dataSqlQuery).Error
	if err != nil {
		log.Fatalf("Failed to execute SQL script: %s", err)
	}
}
