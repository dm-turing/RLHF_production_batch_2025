package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Rating float64 `json:"rating"`
}

type Recommendation struct {
	ProductID string  `json:"productId"`
	Score    float64 `json:"score"`
}

type Feedback struct {
	UserID    string `json:"userId"`
	ProductID string `json:"productId"`
	Rating    int    `json:"rating"`
}

func main() {
	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer driver.Close()

	http.HandleFunc("/recommendations/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Path[len("/recommendations/"):]
		recommendations := getRecommendations(driver, userID)
		json.NewEncoder(w).Encode(recommendations)
	})

	http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		productID := r.URL.Path[len("/products/"):]
		product, err := getProduct(driver, productID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Product not found"})
			return
		}
		json.NewEncoder(w).Encode(product)
	})

	http.HandleFunc("/feedback", func(w http.ResponseWriter, r *http.Request) {
		var feedback Feedback
		if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid feedback data"})
			return
		}
		if err := saveFeedback(driver, feedback); err != nil {
			w.WriteHeader(http.StatusInternalServerError)