package main  
import (  
    "fmt"
    "log"
    "net/http"
    "time"
    "strings"
    "encoding/base64"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)
//URL struct represents the URL model
type URL struct {
    ID        uint `gorm:"primary_key"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time `sql:"index"`
    LongURL   string `gorm:"size:1024"`
    ShortURL  string `gorm:"unique_index"`
    UserID    uint
}

//User struct represents the User model
type User struct {
    ID        uint `gorm:"primary_key"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time `sql:"index"`
    Name      string `gorm:"size:255"`
    Email     string `gorm:"size:255"`
    URLs      []URL
}
//DB global database connection
var DB *gorm.DB

func init() {
    //open a sqlite database connection
    db, err := gorm.Open("sqlite3", "url_shortener.db")
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }
    DB = db
    //migrate the URL and User models
    DB.AutoMigrate(&URL{}, &User{})
}
func main() {
    router := gin.Default()
    //Register endpoint for generating the short URL
    router.POST("/shorten", shortenURL)
    //Register endpoint for retrieving the original URL
    router.GET("/:shortURL", getOriginalURL)
    //Register endpoint for retrieving all the short URLs created by a user
    router.GET("/user/:userID", getUserURLs)
    
    //start the web server
    router.Run(":8080")
}

func generateShortURL(longURL string) string {
    //base62 encoding 
    b := make([]byte, 6)
    _, err := base64.RawURLEncoding.Decode(b, []byte(longURL))
    if err != nil {
        log.Fatalf("Error decoding base64: %v", err)
    }
    return string(b)
}

func shortenURL(c *gin.Context) {
    //parse the request body to extract longURL
    var request struct {
        LongURL string `json:"longURL" binding:"required"`
    }
    if c.BindJSON(&request) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    //generate the shortURL
    shortURL := generateShortURL(request.LongURL)
    //get the userID from the request header
    userID, err := getUserIDFromHeader(c.Request)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    //create a new URL record
    url := URL{LongURL: request.LongURL, ShortURL: shortURL, UserID: userID}
    DB.Create(&url)
    //return the shortURL
    c.JSON(http.StatusCreated, gin.H{"shortURL": fmt.Sprintf("http://localhost:8080/%s", shortURL)})
}

func getOriginalURL(c *gin.Context) {
    //get the shortURL from the request path
    shortURL := c.Param("shortURL")
    //find the URL record in the database
    var url URL
    DB.Where("short_url = ?", shortURL).First(&url)
    if url.ID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
        return
    }
    //return the original URL