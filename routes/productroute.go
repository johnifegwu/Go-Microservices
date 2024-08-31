package productroute

import (
	"encoding/json"
	"net/http"
	"strconv"

	//import database
	database "github.com/johnifegwu/go-microservices/infrastructure"

	// import data models
	"github.com/johnifegwu/go-microservices/models"
)

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
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
