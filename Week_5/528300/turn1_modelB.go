package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

// Product represents a product in the database.
type Product struct {
	Name        string  `db:"name"`
	Sales       int     `db:"sales"`
	SalesAmount float64 `db:"sales_amount"`
}

// db holds the global database connection.
var db *sqlx.DB

func init() {
	// Connect to SQLite database
	var err error
	db, err = sqlx.Open("sqlite3", "./products.db")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Create the products table if it doesn't exist
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY,
        name TEXT,
        sales INTEGER,
        sales_amount REAL
    )`); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Sanitize input query parameters
	groupBy := vars["groupBy"]
	if groupBy != "" {
		groupBy = fmt.Sprintf("GROUP BY %s", groupBy)
	}

	having := vars["having"]
	if having != "" {
		having = fmt.Sprintf("HAVING %s", having)
	}

	orderBy := vars["orderBy"]
	if orderBy != "" {
		orderBy = fmt.Sprintf("ORDER BY %s", orderBy)
	}

	limit := vars["limit"]
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		limit = fmt.Sprintf("LIMIT %d", limitInt)
	}

	// Build the query dynamically
	query := fmt.Sprintf(`
        SELECT name, sum(sales) as sales, sum(sales_amount) as sales_amount
        FROM products
        %s
        %s
        %s
        %s;
    `, groupBy, having, orderBy, limit)

	// Execute the query
	var products []Product
	err := db.Select(&products, query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"products": %+v}`, products)
}

func main() {
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/products/query/{groupBy}/{having}/{orderBy}/{limit}", queryHandler).Methods("GET")

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
