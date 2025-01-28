package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Product struct {
	Name        string  `json:"name"`
	SalesCount  int     `json:"sales_count"`
	SalesAmount float64 `json:"sales_amount"`
}

func setupDatabase() {
	db, err := sql.Open("sqlite3", "./product.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create products table
	createTableSQL := `CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		sales_count INTEGER,
		sales_amount REAL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Insert sample data
	insertProductsSQL := `INSERT INTO products (name, sales_count, sales_amount) VALUES
		('Product A', 100, 2000.00),
		('Product B', 150, 3000.50),
		('Product C', 120, 2500.75),
		('Product D', 80, 1600.30),
		('Product E', 230, 4600.90);`

	_, err = db.Exec(insertProductsSQL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database setup complete.")
}

func queryProducts(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./product.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Retrieve query parameters
	groupBy := r.URL.Query().Get("group_by")
	having := r.URL.Query().Get("having")
	orderBy := r.URL.Query().Get("order_by")
	limitStr := r.URL.Query().Get("limit")

	limit := 5 // Default limit
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit value", http.StatusBadRequest)
			return
		}
	}

	// Construct the SQL query
	query := fmt.Sprintf(`SELECT name, SUM(sales_count) as sales_count, SUM(sales_amount) as sales_amount
						  FROM products
						  GROUP BY %s
						  HAVING %s
						  ORDER BY %s
						  LIMIT %d;`, groupBy, having, orderBy, limit)

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.Name, &product.SalesCount, &product.SalesAmount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func main() {
	setupDatabase()
	http.HandleFunc("/products", queryProducts)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
