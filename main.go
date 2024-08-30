package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// import data models
)

var db *gorm.DB

func init() {
	// Retrieve environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Create the connection string using environment variables
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get the underlying database object:", err)
	}
	sqlDB.SetMaxOpenConns(1000)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	// Migrate the schema and seed data
	migrateAndSeedData()
}

func migrateAndSeedData() {
	// Create the wisdom schema
	var dbErr = db.Exec("CREATE SCHEMA IF NOT EXISTS wisdom;").Error
	if dbErr != nil {
		log.Fatal("Failed to create wisdom schema:", dbErr)
	}

	// Enable the pgcrypto extension
	dbErr = db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error
	if dbErr != nil {
		log.Fatal("Failed to enable pgcrypto extension:", dbErr)
	}

	// Automatically migrate the schema
	err := db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatalf("Failed to migrate schema: %s", err)
	}

	// Check if there is any product data
	var product Product
	result := db.First(&product)
	if result.Error == nil {
		return // Data already exists
	}

	// Seed data if needed (example logic)
	// Read the Data SQL file and execute the commands
	dataSqlFile := "data.sql"
	dataSqlBytes, err := os.ReadFile(dataSqlFile)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %s", err)
	}

	dataSqlQuery := string(dataSqlBytes)
	err = db.Exec(dataSqlQuery).Error
	if err != nil {
		log.Fatalf("Failed to execute SQL script: %s", err)
	}
}

func main() {

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	//Products Router
	router.Get("/products", getProductsHandler)

	//Orders Router
	router.Get("/orders", getOrdersHandler)

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	fmt.Println("Server runing on port 3000 ")

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to listen to server", err)
	}
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters from query string
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Default values for page and pageSize
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Calculate the offset for SQL query
	offset := (page - 1) * pageSize

	// Query the products table with LIMIT and OFFSET for pagination
	var products []Product
	result := db.Limit(pageSize).Offset(offset).Order("product_id").Find(&products)
	if result.Error != nil {
		http.Error(w, "Failed to query products", http.StatusInternalServerError)
		return
	}

	// Convert the products to JSON and write to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products to JSON", http.StatusInternalServerError)
	}
}

func getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[{\"orderId\": 100, \"productId\":\"gdhfg7473\", \"productName\":\"Samsung S22 Ultra\"}]"))
}
