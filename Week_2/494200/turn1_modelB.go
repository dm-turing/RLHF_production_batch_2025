package main  
import (  
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
)

type Product struct {  
    ID    string `json:"id"`
    Name  string `json:"name"`
    Price float64 `json:"price"`
}

var products []Product  

// Logging Middleware  
func LoggingMiddleware(next http.Handler) http.Handler {  
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {  
        logger := logrus.WithFields(logrus.Fields{  
            "method": r.Method,
            "path":   r.URL.Path,
        })

        // Start timer
        start := logrus.Now()

        // Handle request
        next.ServeHTTP(w, r)

        // Log request details
        logger.WithFields(logrus.Fields{  
            "status": w.Header().Get("Status"),
            "latency": logrus.Since(start),
        }).Info("Request processed")
    })  
}

func GetProducts(w http.ResponseWriter, r *http.Request) {  
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {  
    params := mux.Vars(r)
    for _, item := range products {  
        if item.ID == params["id"] {  
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(&Product{})
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {  
    var newProduct Product
    json.NewDecoder(r.Body).Decode(&newProduct)
    products = append(products, newProduct)
    json.NewEncoder(w).Encode(newProduct)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {  
    params := mux.Vars(r)
    for index, item := range products {  
        if item.ID == params["id"] {  
            products = append(products[:index], products[index+1:]...)
            var updatedProduct Product
            json.NewDecoder(r.Body).Decode(&updatedProduct)
            updatedProduct.ID = params["id"]
            products = append(products, updatedProduct)
            json.NewEncoder(w).Encode(updatedProduct)
            return
        }
    }
    json.NewEncoder(w).Encode(products)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {  
    params := mux.Vars(r)
    for index, item := range products {  
        if item.ID == params["id"] {  
            products = append(products[:index], products[index+1:]...)
            break
        }
    }
    json.NewEncoder(w).Encode(products)
}

func main() {  
    router := mux.NewRouter()

    // Initialize products data
    products = append(products, Product{ID: "1", Name: "Product1", Price: 100})
    products = append(products, Product{ID: "2", Name: "Product2", Price: 200})

    // Define routes
    router.HandleFunc("/products", GetProducts).Methods("GET")
    router.HandleFunc("/products/{id}", GetProduct).Methods("GET")
    router.HandleFunc("/products", CreateProduct).Methods("POST")
    router.HandleFunc("/products/{id}", UpdateProduct).Methods("PUT")
    router.HandleFunc("/products/{id}", DeleteProduct).Methods("DELETE")