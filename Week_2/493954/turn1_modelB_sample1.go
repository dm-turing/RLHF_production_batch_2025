package main

import "fmt"

// Product represents an E-Commerce product
type Product struct {
	Name       string                 `json:"name"`
	Price      float32                `json:"price"`
	Attributes map[string]interface{} `json:"attributes"`
}

func main() {
	// Define a product with complex attributes
	product := Product{
		Name:  "Smartphone XYZ",
		Price: 799.99,
		Attributes: map[string]interface{}{
			"Display": map[string]string{
				"Size":       "6.5 inches",
				"Type":       "OLED",
				"Resolution": "2340 x 1080 pixels",
			},
			"Camera": []string{
				"Rear: 12MP, 1080p video",
				"Front: 8MP, 720p video",
			},
			"Memory":  "8GB RAM",
			"Storage": "128GB/256GB/512GB",
			"Battery": "5000mAh",
		},
	}

	fmt.Printf("Product: %+v\n", product)
}
