package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	router := mux.NewRouter()
	// Create a new logger
	logger := logrus.New()
	// Set the log level to warn for production environment, trace otherwise
	if os.Getenv("GO_ENV") == "production" {
		logger.SetLevel(logrus.WarnLevel)
	} else {
		logger.SetLevel(logrus.TraceLevel)
	}
	// Attach the logger to the request-response logging middleware
	logrusMiddleware := NewLogrusMiddleware(logger)
	// Assign the middleware to the Gorilla mux router
	router.Use(logrusMiddleware.Middleware)

	// Define routes
	router.HandleFunc("/products", GetProductsHandler).Methods("GET")
	router.HandleFunc("/products/{id}", GetProductHandler).Methods("GET")
	router.HandleFunc("/orders", CreateOrderHandler).Methods("POST")
	//... other routes

	log.Fatal(http.ListenAndServe(":8080", router))
} 

// Product represents a product in the E-Commerce website
type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Price float32 `json:"price"`
}

// example handlers for demonstration purposes
func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	// Implement the product fetching logic
	products := []Product{
		{ID: "1", Name: "Phone", Price: 699.99},
		{ID: "2", Name: "Laptop", Price: 1099.99},
	}

	respondWithJSON(w, http.StatusOK, products)
}

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
    // Implement the product fetching logic
    productID := mux.Vars(r)["id"]
    product := Product{ID: productID, Name: "Sample Product", Price: 99.99}
    respondWithJSON(w, http.StatusOK, product)
}
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
    // Implement the order creation logic
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Order created successfully"})
}
// responseWithJSON writes the given payload with the specified HTTP status code and a 'Content-Type: application/json' header.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// LogrusMiddleware is the custom middleware that log the requests and responses.
type LogrusMiddleware struct {
	logger *logrus.Logger
}

// NewLogrusMiddleware returns a new LogrusMiddleware
func NewLogrusMiddleware(logger *logrus.Logger) *LogrusMiddleware {
	return &LogrusMiddleware{logger: logger}
}

func (lm *LogrusMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Process request
		next.ServeHTTP(w, r)

		// Calculate request time
		elapsed := time.Since(start)

		// Extract request information
		reqMethod := r.Method
		reqURI := r.RequestURI
		reqStatus := w.(statuswriter.StatusWriter).Status()
		reqSize := w.(statuswriter.StatusWriter).Size()