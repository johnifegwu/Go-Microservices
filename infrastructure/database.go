package database

import (
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/johnifegwu/go-microservices/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

//go:embed data.sql
var fdat embed.FS

func InitDb() {
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
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Configure the connection pool
	SqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get the underlying database object:", err)
	}
	SqlDB.SetMaxOpenConns(1000)
	SqlDB.SetMaxIdleConns(10)
	SqlDB.SetConnMaxLifetime(time.Minute * 5)

	// Migrate the schema and seed data
	MigrateAndSeedData()
}

func MigrateAndSeedData() {

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
