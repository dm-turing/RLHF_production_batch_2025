package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

// Movie represents a movie recommendation.
type Movie struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Genre  string  `json:"genre"`
	Rating float64 `json:"rating"`
}

func main() {
	filePath := "movies.csv" // Replace with your CSV file path
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	header := scanner.Text()
	fmt.Println(header)
	if header != "id,title,genre,rating" {
		log.Fatalf("Invalid CSV header: %s", header)
	}

	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading CSV row: %v", err)
		}

		if len(row) != 4 {
			log.Fatalf("Row has incorrect number of fields: %v", row)
		}
		id, _ := strconv.Atoi(row[0])
		rating, _ := strconv.ParseFloat(row[3], 32)

		movie := Movie{
			ID:     id,
			Title:  row[1],
			Genre:  row[2],
			Rating: rating,
		}

		movieKey := fmt.Sprintf("movie:%d", movie.ID)
		fmt.Println(movieKey)
		err = redisClient.HMSet(movieKey, map[string]interface{}{
			"title":  movie.Title,
			"genre":  movie.Genre,
			"rating": movie.Rating,
		}).Err()
		if err != nil {
			log.Fatalf("Error loading movie into Redis: %v", err)
		}
	}

	fmt.Println("CSV data loaded into Redis successfully!")
}
