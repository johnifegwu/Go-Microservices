package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	//import database
	database "github.com/johnifegwu/go-microservices/infrastructure"

	// import data routes
	productroute "github.com/johnifegwu/go-microservices/routes"
)

func main() {

	database.InitDb()

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	//Products Router
	router.Get("/products", productroute.GetProductsHandler)

	//Orders Router
	router.Get("/productbyid", productroute.GetProductByIdHandler)

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
