package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// ProductSales represents the structure for product sales data
type ProductSales struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	TotalSales  int     `json:"total_sales"`
	TotalAmount float64 `json:"total_amount"`
}

func main() {
	// Open the SQLite database, create it if it doesn't exist
	db, err := sql.Open("sqlite3", "./sales.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare the database: products table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		total_sales INTEGER,
		sales_amount REAL
	);`
	if _, err = db.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}

	// Insert sample data
	sampleDataSQL := `INSERT INTO products (name, total_sales, sales_amount) VALUES (?, ?, ?);`
	db.Exec("DELETE FROM products") // Clean previous entries for consistency
	for _, i := range []ProductSales{
		ProductSales{Name: "Product A", TotalSales: 120, TotalAmount: 2500.0},
		ProductSales{Name: "Product B", TotalSales: 60, TotalAmount: 1500.0},
		ProductSales{Name: "Product C", TotalSales: 250, TotalAmount: 5000.0},
	} {
		if _, err := db.Exec(sampleDataSQL, i.Name, i.TotalSales, i.TotalAmount); err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/query-sales", func(w http.ResponseWriter, r *http.Request) {
		handleQuerySales(w, r, db)
	})

	log.Println("Server is running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleQuerySales(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := `
		SELECT 
			id, 
			name, 
			SUM(total_sales) as total_sales, 
			SUM(sales_amount) as total_amount
		FROM products
		GROUP BY name
		HAVING SUM(total_sales) > 50
		ORDER BY total_amount DESC
		LIMIT 5;
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []ProductSales
	for rows.Next() {
		var ps ProductSales
		if err := rows.Scan(&ps.ID, &ps.Name, &ps.TotalSales, &ps.TotalAmount); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, ps)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
