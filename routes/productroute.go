package productroute

import (
	"encoding/json"
	"net/http"
	"strconv"

	//import database
	"github.com/google/uuid"
	database "github.com/johnifegwu/go-microservices/infrastructure"
	"gorm.io/gorm"

	// import data models
	"github.com/johnifegwu/go-microservices/models"
)

func GetProductByIdHandler(w http.ResponseWriter, r *http.Request) {

	uuidStr := r.URL.Query().Get("id")

	// Check if the id is provided
	// Parse the string into a uuid.UUID
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		http.Error(w, "Product ID not provided or invalid", http.StatusBadRequest)
		return
	}

	// Query the product by product_id
	var product models.Product
	result := database.DB.Model(models.Product{ProductID: parsedUUID}).First(&product)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to query product", http.StatusInternalServerError)
		}
		return
	}

	// Convert the product to JSON and write to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode product to JSON", http.StatusInternalServerError)
	}
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters from query string
	pageStr := r.URL.Query().Get("pageIndex")
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
	var products []models.Product
	result := database.DB.Limit(pageSize).Offset(offset).Order("product_id").Find(&products)
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
