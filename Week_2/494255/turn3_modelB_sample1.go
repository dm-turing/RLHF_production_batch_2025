package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

func loadToRedisCSV(movies []Movie, redisClient *redis.Client) {
	// Set a key to store the movies as a CSV string
	redisKey := "movies_csv"

	// Create a CSV writer
	var csvBuffer bytes.Buffer
	csvWriter := csv.NewWriter(&csvBuffer)

	// Write the header
	if err := csvWriter.Write([]string{"ID", "Title", "Genre", "Rating"}); err != nil {
		log.Fatal(err)
	}

	// Write the movie data
	for _, movie := range movies {
		if err := csvWriter.Write([]string{fmt.Sprintf("%d", movie.ID), movie.Title, movie.Genre, fmt.Sprintf("%f", movie.Rating)}); err != nil {
			log.Fatal(err)
		}
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Fatal(err)
	}

	// Set the CSV data to Redis
	if err := redisClient.Set(redisKey, csvBuffer.String(), 0).Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Movie data loaded to Redis as CSV successfully!")
}

func main() {
	filePath := "movies.json" // Example file path

	// Load and transform the movie data as before
	extractor := NewDataExtractor()
	transformer := &FilterAndSortTransformer{}

	movies, err := extractor.Extract(filePath)
	if err != nil {
		log.Fatal(err)
	}

	transformedMovies, err := transformer.Transform(movies)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Default Redis port
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	defer redisClient.Close()

	// Ping Redis to check connection
	if _, err := redisClient.Ping().Result(); err != nil {
		log.Fatal(err)
	}

	// Load data to Redis as CSV
	loadToRedisCSV(transformedMovies, redisClient)
}
