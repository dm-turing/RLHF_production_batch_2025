

package main  
import (  
    "fmt"
    "math/rand"
    "time"
)

func main() {  
    // Step 1: Generate a map of type map[string]int with 10000 random values.
    // Initialize a map and a slice
    m := make(map[string]int)  
    keys := make([]string, 0)  
    rand.Seed(time.Now().UnixNano())  

    // Generate 10000 random key-value pairs and add keys to the slice
    for i := 0; i < 10000; i++ {
        key := fmt.Sprintf("key%06d", i)
        value := rand.Intn(10000)
        m[key] = value
        keys = append(keys, key)
    }

    // Step 2: Iterate through the keys in the slice and print to the console.
    start := time.Now()
    for _, key := range keys {  
        fmt.Println(key)  
    }  
    elapsed := time.Since(start)
    fmt.Printf("Iterating through the slice: %s\n", elapsed)

    // Step 3: Iterate through the keys in the map and print to the console.
    start = time.Now()
    for key := range m {  
        fmt.Println(key)  
    }  
    elapsed = time.Since(start)
    fmt.Printf("Iterating through the map: %s\n", elapsed)
    
} 

