package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

type Movie struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Genre  string  `json:"genre"`
	Rating float64 `json:"rating"`
}

type FormatDetector interface {
	Detect(io.Reader) (string, error)
}

type CSVFormatDetector struct {
}

func (d *CSVFormatDetector) Detect(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ",") {
			return "csv", nil
		}
	}
	return "", fmt.Errorf("format not detected")
}

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

type FormatExtractor interface {
	Extract(io.Reader, string) ([]Movie, error)
}

type CSVFormatExtractor struct {
}

func (e *CSVFormatExtractor) Extract(r io.Reader, format string) ([]Movie, error) {
	if format != "csv" {
		return nil, fmt.Errorf("invalid format: %s", format)
	}

	decoder := csv.NewReader(r)
	var movies []Movie
	_, err := decoder.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := decoder.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) != 4 {
			log.Fatalf("Row has incorrect number of fields: %v", row)
			continue
		}
		id, _ := strconv.Atoi(row[0])
		rating, _ := strconv.ParseFloat(row[3], 32)
		movie := Movie{
			ID:     id,
			Title:  row[1],
			Genre:  row[2],
			Rating: rating,
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

type JSONFormatExtractor struct {
}

func (e *JSONFormatExtractor) Extract(r io.Reader, format string) ([]Movie, error) {
	if format != "json" {
		return nil, fmt.Errorf("invalid format: %s", format)
	}

	var movies []Movie
	if err := json.NewDecoder(r).Decode(&movies); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return movies, nil
}

type Transformer interface {
	Transform([]Movie) ([]Movie, error)
}

type FilterAndSortTransformer struct {
}

func (t *FilterAndSortTransformer) Transform(movies []Movie) ([]Movie, error) {
	filteredMovies := []Movie{}
	for _, movie := range movies {
		if movie.Rating >= 3 && movie.Rating <= 5 {
			filteredMovies = append(filteredMovies, movie)
		}
	}

	sort.Slice(filteredMovies, func(i, j int) bool {
		return filteredMovies[i].Title > filteredMovies[j].Title
	})

	return filteredMovies, nil
}

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

func loadToRedis(movies []Movie, redisClient *redis.Client) error {
	for _, movie := range movies {
		movieKey := fmt.Sprintf("movie:%d", movie.ID)
		if err := redisClient.HMSet(movieKey, map[string]interface{}{
			"title":  movie.Title,
			"genre":  movie.Genre,
			"rating": movie.Rating,
		}).Err(); err != nil {
			return fmt.Errorf("error loading movie %d into Redis: %v", movie.ID, err)
		}
	}
	return nil
}

func main() {
	filePath := "movies.json" // Replace with your data file path
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	extractor := NewDataExtractor()

	switch filePath {
	case "movies.csv":
		extractor.SetExtractor(&CSVFormatExtractor{})
	case "movies.json":
		extractor.SetExtractor(&JSONFormatExtractor{})
	default:
		log.Fatalf("Unknown file type: %s", filePath)
	}

	movies, err := extractor.Extract(filePath)
	if err != nil {
		log.Fatal(err)
	}

	transformer := &FilterAndSortTransformer{}
	transformedMovies, err := transformer.Transform(movies)
	if err != nil {
		log.Fatal(err)
	}

	if err := loadToRedis(transformedMovies, redisClient); err != nil {
		log.Fatal(err)
	}

	fmt.Println("ETL process completed successfully!")
}
