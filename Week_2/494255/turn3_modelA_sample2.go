package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)

func main() {
	filePath := "movies.json" // Replace with your JSON file path
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

	decoder := json.NewDecoder(file)
	var movies []Movie
	err = decoder.Decode(&movies)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	for _, movie := range movies {
		movieKey := fmt.Sprintf("movie:%d", movie.ID)
		err = redisClient.HMSet(movieKey, map[string]interface{}{
			"title":  movie.Title,
			"genre":  movie.Genre,
			"rating": movie.Rating,
		}).Err()
		if err != nil {
			log.Fatalf("Error loading movie into Redis: %v", err)
		}
	}

	fmt.Println("JSON data loaded into Redis successfully!")
}
