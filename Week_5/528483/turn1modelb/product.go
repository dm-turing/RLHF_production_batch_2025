// product.go
package model

import "math/rand"

type Product struct {
    ID          int
    Name        string
    Description string
    Quantity    int
    Price       float64
}

// NewRandomProduct generates a new random product.
func NewRandomProduct() *Product {
    return &Product{
        ID:          rand.Intn(1000),
        Name:        randomString(10),
        Description: randomString(20),
        Quantity:    rand.Intn(10) + 1,
        Price:       rand.Float64() * 100,
    }
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}