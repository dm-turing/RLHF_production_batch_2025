package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/go-redis/redis"
)

// Movie represents a movie recommendation.
type Movie struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Genre  string  `json:"genre"`
	Rating float64 `json:"rating"`
}

// FormatDetector detects the format of a given file.
type FormatDetector interface {
	Detect(io.Reader) (string, error)
}

// FormatExtractor extracts data from a file based on its format.
type FormatExtractor interface {
	Extract(io.Reader, string) ([]Movie, error)
}

// DataExtractor extracts movie recommendations from a file.
type DataExtractor struct {
	detector  FormatDetector
	extractor FormatExtractor
}

func NewDataExtractor() *DataExtractor {
	return &DataExtractor{
		detector:  &JSONFormatDetector{},
		extractor: &JSONFormatExtractor{},
	}
}

// JSONFormatDetector detects JSON format.
type JSONFormatDetector struct {
}

func (d *JSONFormatDetector) Detect(r io.Reader) (string, error) {
	buf := make([]byte, 1024)
	_, err := r.Read(buf)
	if err != nil {
		return "", err
	}

	// r.Seek(0, 0) // Reset the reader to the beginning

	if strings.Contains(string(buf), "{") {
		return "json", nil
	}
	return "", fmt.Errorf("format not detected")
}

// JSONFormatExtractor extracts data in JSON format.
type JSONFormatExtractor struct {
}

func (e *JSONFormatExtractor) Extract(r io.Reader, format string) ([]Movie, error) {
	if format != "json" {
		return nil, fmt.Errorf("invalid format: %s", format)
	}

	var movies []Movie
	if err := json.NewDecoder(r).Decode(&movies); err != nil {
		return nil, err
	}

	return movies, nil
}

// Transformer transforms the extracted movie data.
type Transformer interface {
	Transform([]Movie) ([]Movie, error)
}

// FilterAndSortTransformer filters movies by rating and sorts them in reverse lexicographical order.
type FilterAndSortTransformer struct {
}

func (t *FilterAndSortTransformer) Transform(movies []Movie) ([]Movie, error) {
	filteredMovies := []Movie{}
	for _, movie := range movies {
		if movie.Rating >= 3 && movie.Rating <= 5 {
			filteredMovies = append(filteredMovies, movie)
		}
	}

	// Sort movies by title in reverse lexicographical order
	sort.Slice(filteredMovies, func(i, j int) bool {
		return filteredMovies[i].Title > filteredMovies[j].Title
	})

	return filteredMovies, nil
}

func (e *DataExtractor) SetDetector(d FormatDetector) {
	e.detector = d
}

func (e *DataExtractor) SetExtractor(ext FormatExtractor) {
	e.extractor = ext
}

func (e *DataExtractor) Extract(filePath string) ([]Movie, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	format, err := e.detector.Detect(file)
	if err != nil {
		return nil, err
	}

	file.Seek(0, 0) // Reset the reader to the beginning

	return e.extractor.Extract(file, format)
}

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
