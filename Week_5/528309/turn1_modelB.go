package main  
import (  
    "encoding/xml"  
    "fmt"  
    "math/rand"  
    "net/url"  
    "os"  
    "strings"  
    "time"  
)  

const (  
    baseURL       = "https://example.com"  
    numURLs        = 1000  
    maxURLDepth   = 3  
    changefreqs   = "always hourly daily weekly monthly yearly never"  
    priorities     = "0.0 0.1 0.2 0.3 0.4 0.5 0.6 0.7 0.8 0.9 1.0"  
    mediaTypes    = "image/jpeg image/png video/mp4 video/webm"  
    extensionTypes = ".jpg .png .mp4 .webm"  
)  

type urlset struct {  
    XMLName xml.Name `xml:"urlset"`  
    XMLNS   string   `xml:"xmlns,attr"`  
    URLs    []url    `xml:"url"`  
}  

type url struct {  
    Loc        string    `xml:"loc"`  
    Lastmod    time.Time `xml:"lastmod"`  
    Changefreq string    `xml:"changefreq"`  
    Priority   string    `xml:"priority"`  
    Image      *image    `xml:"image,omitempty"`  
    Video      *video    `xml:"video,omitempty"`  
}  

type image struct {  
    Loc     string `xml:"loc"`  
    Title   string `xml:"title"`  
    Caption string `xml:"caption"`  
    GeoLocation   *geoLocation `xml:"geo_location,omitempty"`
}  

type video struct {  
    Loc         string `xml:"loc"`  
    Title       string `xml:"title"`  
    Description string `xml:"description"`  
    ContentLoc  string `xml:"content_loc"`  
    ThumbnailLoc string `xml:"thumbnail_loc"`  
}  

type geoLocation struct {  
    Lat string `xml:"lat"`  
    Lon string `xml:"lon"`  
}  

func randomString(n int) string {  
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"  
    b := make([]byte, n)  
    for i := range b {  
        b[i] = letters[rand.Intn(len(letters))]  
    }  
    return string(b)  
}  

func generateRandomURL(baseURL string, depth int) string {  
    if depth == 0 {  
        return baseURL  
    }  
    path := "/" + randomString(rand.Intn(10)+1)  
    if rand.Float32() < 0.5 {  
        path += "/" + randomString(rand.Intn(10)+1)  
    }  
    return generateRandomURL(baseURL+path, depth-1)  
}  

func generateSitemap() *urlset {  
    rand.Seed(time.Now().UnixNano())  
    us := &urlset{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9"}  
    for i := 0; i < numURLs; i++ {  
        u := &url{  
            Loc:        generateRandomURL(baseURL, rand.Intn(maxURLDepth)+1),  
            Lastmod:    time.Now().Add(-time.Duration(rand.Intn(365*24*60*60)) * time.Second),  
            Changefreq: strings.Fields(changefreqs)[rand.Intn(len(strings.Fields(changefreqs)))],  
            Priority:   strings.Fields(priorities)[rand.Intn(len(strings.Fields(priorities)))],  