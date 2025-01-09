package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Product struct {
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Rating float64 `json:"rating"`
}

func fetchProducts(w http.ResponseWriter, r *http.Request) {
	// Extract price range parameters
	minPrice, err := strconv.ParseFloat(r.URL.Query().Get("minPrice"), 64)
	if err != nil {
		http.Error(w, "Invalid minPrice parameter", http.StatusBadRequest)
		return
	}
	maxPrice, err := strconv.ParseFloat(r.URL.Query().Get("maxPrice"), 64)
	if err != nil {
		http.Error(w, "Invalid maxPrice parameter", http.StatusBadRequest)
		return
	}

	// Connect to Neo4j
	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer driver.Close()

	// Run the query
	session, err := driver.Session(neo4j.AccessModeRead)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	results, err := session.Run(
		`
		MATCH (p:Product)
		WHERE $minPrice <= p.price AND p.price <= $maxPrice AND p.rating >= 4
		RETURN p.name AS name, p.price AS price, p.rating AS rating
		`,
		map[string]interface{}{
			"minPrice": minPrice,
			"maxPrice": maxPrice,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// defer results.Close()
	var products []Product
	for results.Next() {
		var product Product
		name, ok := results.Record().Get("name")
		if !ok {
			http.Error(w, "error retrieving name", http.StatusInternalServerError)
			return
		}
		price, ok := results.Record().Get("price")
		if !ok {
			http.Error(w, "error retrieving price", http.StatusInternalServerError)
			return
		}
		rating, ok := results.Record().Get("rating")
		if !ok {
			http.Error(w, "error retrieving rating", http.StatusInternalServerError)
			return
		}
		/*if err := results.Record().Scan(&product.Name, &product.Price, &product.Rating); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}*/
		product = Product{Name: name.(string), Price: price.(float64), Rating: rating.(float64)}
		products = append(products, product)
	}

	// Respond with JSON
	json.NewEncoder(w).Encode(products)
}

func main() {
	http.HandleFunc("/api/products", fetchProducts)
	fmt.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
