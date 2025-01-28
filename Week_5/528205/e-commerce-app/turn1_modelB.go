package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"e-commerce-app/internal/config"
	"e-commerce-app/internal/database"
	"e-commerce-app/internal/services"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Initialize user management service
	ums := services.NewUserManagementService(db)

	// Create a new router
	r := mux.NewRouter()

	// Define endpoints for User Management Service
	r.HandleFunc("/users", ums.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", ums.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", ums.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", ums.DeleteUser).Methods(http.MethodDelete)

	// Run the server
	log.Printf("Server is running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}  