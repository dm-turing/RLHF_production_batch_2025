package main

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
	Image      string  `xml:"image,omitempty"`
	Video      string  `xml:"video,omitempty"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Urls    []URL    `xml:"url"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	baseURL := "https://example.com"
	sitemap := URLSet{
		Urls: generateRandomURLs(baseURL, 10), // Generate a sitemap with 10 random URLs
	}

	fileName := "sitemap.xml"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(sitemap); err != nil {
		fmt.Println("Error encoding XML to file:", err)
		return
	}

	fmt.Println("Sitemap successfully generated into", fileName)
}

func generateRandomURLs(baseURL string, count int) []URL {
	urls := make([]URL, count)

	changeFreqOptions := []string{"always", "hourly", "daily", "weekly", "monthly", "yearly", "never"}

	for i := 0; i < count; i++ {
		urls[i] = URL{
			Loc:        fmt.Sprintf("%s/page%d", baseURL, i),
			LastMod:    randomDate(),
			ChangeFreq: changeFreqOptions[rand.Intn(len(changeFreqOptions))],
			Priority:   rand.Float64(),
			Image:      fmt.Sprintf("%s/image%d.jpg", baseURL, i), // optional for randomness
			Video:      fmt.Sprintf("%s/video%d.mp4", baseURL, i), // optional for randomness
		}
	}

	return urls
}

func randomDate() string {
	latestTime := time.Now().Unix()
	oldestTime := latestTime - 5*365*24*60*60 // Generate random dates in the past 5 years
	randomTime := rand.Int63n(latestTime-oldestTime) + oldestTime
	return time.Unix(randomTime, 0).Format("2006-01-02")
}
