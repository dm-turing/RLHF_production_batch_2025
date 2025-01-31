package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Quantity    int
	Price       float64
}

func main() {
	products := populateProducts()
	fmt.Println("Products in Map:")
	for _, product := range products {
		fmt.Println(product)
	}

	db, err := sql.Open("sqlite3", "./ecommerce.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		log.Fatal(err)
	}

	err = insertProducts(db, products)
	if err != nil {
		log.Fatal(err)
	}

	retrievedProducts, err := retrieveProducts(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nRetrieved Products from DB:")
	for _, product := range retrievedProducts {
		fmt.Println(product)
	}

	fmt.Println("\nComparison Results:")
	compareProducts(products, retrievedProducts)
}

func populateProducts() map[string]Product {
	productMap := make(map[string]Product)
	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("Product%d", i)
		productMap[name] = Product{
			ID:          i,
			Name:        name,
			Description: fmt.Sprintf("Description of Product%d", i),
			Quantity:    rand.Intn(100),
			Price:       float64(rand.Intn(10000)) / 100,
		}
	}
	return productMap
}

func createTable(db *sql.DB) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS products (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT,
		"description" TEXT,
		"quantity" INTEGER,
		"price" REAL
	);`
	_, err := db.Exec(createTableSQL)
	return err
}

func insertProducts(db *sql.DB, products map[string]Product) error {
	for _, product := range products {
		insertProductSQL := `INSERT INTO products(name, description, quantity, price) VALUES (?, ?, ?, ?)`
		statement, err := db.Prepare(insertProductSQL)
		if err != nil {
			return err
		}
		_, err = statement.Exec(product.Name, product.Description, product.Quantity, product.Price)
		if err != nil {
			return err
		}
	}
	return nil
}

func retrieveProducts(db *sql.DB) (map[string]Product, error) {
	retrievedProducts := make(map[string]Product)
	row, err := db.Query("SELECT id, name, description, quantity, price FROM products")
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var product Product
		err = row.Scan(&product.ID, &product.Name, &product.Description, &product.Quantity, &product.Price)
		if err != nil {
			return nil, err
		}
		retrievedProducts[product.Name] = product
	}
	return retrievedProducts, nil
}

func compareProducts(original, retrieved map[string]Product) {
	for id, product := range original {
		if retrievedProduct, exists := retrieved[product.Name]; exists {
			// Disregard ids
			product.ID, retrievedProduct.ID = 0, 0
			if retrievedProduct == product {
				fmt.Printf("Product ID %s matches.\n", id)
			} else {
				fmt.Printf("Product ID %s does not match. Original: %+v, Retrieved: %+v\n", id, product, retrievedProduct)
			}
		} else {
			fmt.Printf("Product ID %s not found in retrieved products.\n", id)
		}
	}
}
