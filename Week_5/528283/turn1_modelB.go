package main
  
import (  
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strconv"
  
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type ProductSales struct {
    Product string  `json:"product"`
    Sales   int     `json:"sales"`
    Amount  float64 `json:"amount"`
}

var db *sqlx.DB

func init() {
    // Initialize the SQLite database
    var err error
    db, err = sqlx.Connect("sqlite3", "./products.db")
    if err != nil {
        log.Fatalf("Error connecting to SQLite database: %v", err)
    }

    // Create the table if it doesn't exist
    createTableSQL := `CREATE TABLE IF NOT EXISTS products_sales (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        product TEXT NOT NULL,
                        sales INTEGER NOT NULL,
                        amount REAL NOT NULL
                    )`

    _, err = db.Exec(createTableSQL)
    if err != nil {
        log.Fatalf("Error creating table: %v", err)
    }
}

func main() {
    // Close the database connection at the end
    defer db.Close()
    
    router := gin.Default()
    router.GET("/query-products", handleProductQuery)
    
    // Populate sample data for testing
    populateSampleData()
    
    log.Fatal(router.Run(":8080"))
}
  
// Populate sample data
func populateSampleData() {
    sampleData := []struct {
        Product string  `json:"product"`
        Sales   int     `json:"sales"`
        Amount  float64 `json:"amount"`
    }{
        {Product: "Pencil", Sales: 100, Amount: 200},
        {Product: "Paper", Sales: 50, Amount: 100},
        {Product: "Pen", Sales: 80, Amount: 160},
        {Product: "Eraser", Sales: 30, Amount: 60},
        {Product: "Ruler", Sales: 60, Amount: 120},
    }

    for _, data := range sampleData {
        insertSQL := `INSERT INTO products_sales (product, sales, amount) VALUES (?, ?, ?)`
        _, err := db.Exec(insertSQL, data.Product, data.Sales, data.Amount)
        if err != nil {
            log.Printf("Error populating data: %v", err)
        }
    }
}

// Expose the REST endpoint for executing the SQL query
func handleProductQuery(c *gin.Context) {
    query := `SELECT product, SUM(sales) AS sales, SUM(amount) AS amount 
             FROM products_sales 
             GROUP BY product 
             HAVING sales > ? 
             ORDER BY amount DESC 
             LIMIT ?`
    
    // Validate request parameters
    salesLimitStr := c.DefaultQuery("sales_limit", "0")
    salesLimit, err := strconv.Atoi(salesLimitStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sales_limit parameter"})
        return