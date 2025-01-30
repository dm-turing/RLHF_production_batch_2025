package main  
 
import (  
    "database/sql"  
    "fmt"  
    "math/rand"  
    "net/url"  
    "strings"  
    "time"  
 
    _ "github.com/mattn/go-sqlite3"  
)  
 
// Database-related constants  
const dbFile = "urls.db"  
const createTableSQL = `CREATE TABLE IF NOT EXISTS urls (id INTEGER PRIMARY KEY AUTOINCREMENT, full_url TEXT, short_url TEXT)`  
 
// Function to generate a random string  
func randomString(length int) string {  
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"  
    var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))  
    result := make([]byte, length)  
    for i := range result {  
        result[i] = charset[seededRand.Intn(len(charset))]  
    }  
    return string(result)  
}  
 
// Function to generate a random URL  
func generateRandomURL() string {  
    baseURL := "https://www.example.com"  
    path := "/" + randomString(rand.Intn(10)+5) // Random path of length 5 to 15  
 
    params := url.Values{}  
    params.Add("id", randomString(5))  
    params.Add("page", fmt.Sprintf("%d", rand.Intn(100)+1))  
    params.Add("sort", randomString(6))  
    params.Add("filter", randomString(4))  
    params.Add("limit", fmt.Sprintf("%d", rand.Intn(50)+1))  
 
    return baseURL + path + "?" + params.Encode()  
}  
 
// Function to shorten a URL if it's longer than 25 characters  
func shortenURL(longURL string) string {  
    if len(longURL) > 25 {  
        return "https://short.url/" + randomString(7) // Short URL with random 7 character string  
    }  
    return longURL  
}  
 
// Function to store URLs in the SQLite database  
func storeURLs(numURLs int) {  
    db, err := sql.Open("sqlite3", dbFile)  
    if err != nil {  
        fmt.Println(err)  
        return  
    }  
    defer db.Close()  
 
    _, err = db.Exec(createTableSQL)  
    if err != nil {  
        fmt.Println(err)  
        return  
    }  
 
    tx, err := db.Begin()  
    if err != nil {  
        fmt.Println(err)  
        return  
    }  
 
    stmt, err := tx.Prepare("INSERT INTO urls (full_url, short_url) VALUES (?, ?)")  
    if err != nil {  
        fmt.Println(err)  
        return  
    }  
    defer stmt.Close()  
 
    for i := 0; i < numURLs; i++ {  
        fullURL := generateRandomURL()  
        shortURL := shortenURL(fullURL)  
        _, err := stmt.Exec(fullURL, shortURL)  
        if err != nil {  
            fmt.Println(err)  
            return  
        }  
    }  
 
    err = tx.Commit()  
    if err != nil {  
        fmt.Println(err)  
        return  
    }  
    fmt.Println("URLs stored successfully.")  
}  
 
func main() {  
    // Generate and store 10000 URLs