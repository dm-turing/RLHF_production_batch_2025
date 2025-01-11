package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
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

// CSVFormatDetector detects CSV format.
type CSVFormatDetector struct {
}

// CSVFormatExtractor extracts data in CSV format.
type CSVFormatExtractor struct {
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

func NewDataExtractor() *DataExtractor {
	return &DataExtractor{
		detector:  &CSVFormatDetector{},
		extractor: &CSVFormatExtractor{},
	}
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

// TransactionallyAndable interface allows for atomic execution of
// multiple operations on a shared dataset.
type TransactionalAndable interface {
	Execute() error
}

// MovieTransformer is used to transform movie data.
type MovieTransformer struct {
	Movies []Movie
}

func (mt *MovieTransformer) Execute() error {
	// Add your transformation logic here.
	// For example, here we capitalize the genre of each movie.
	for i := range mt.Movies {
		mt.Movies[i].Genre = strings.Title(mt.Movies[i].Genre)
	}
	return nil
}

// Added a FilteredByRating method to filter movies based on rating.
func (e *DataExtractor) FilteredByRating(filePath string, minRating, maxRating float64) ([]Movie, error) {
	movies, err := e.Extract(filePath)
	if err != nil {
		return nil, err
	}

	var filteredMovies []Movie
	for _, movie := range movies {
		if movie.Rating >= minRating && movie.Rating <= maxRating {
			filteredMovies = append(filteredMovies, movie)
		}
	}
	return filteredMovies, nil
}

// ByRatingDesc is a custom sort type to sort movies in reverse lexicographical order of ratings.
type ByRatingDesc []Movie

func (m ByRatingDesc) Len() int           { return len(m) }
func (m ByRatingDesc) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ByRatingDesc) Less(i, j int) bool { return m[i].Rating > m[j].Rating }

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

func main() {
	filePath := "movies.csv" // Assume this file contains movie data
	extractor := NewDataExtractor()

	// After extraction, we can perform transformations and filtering.
	minRating, maxRating := 3.0, 5.0
	filteredMovies, err := extractor.FilteredByRating(filePath, minRating, maxRating)
	if err != nil {
		log.Fatal(err)
	}

	// Optional: Transforming the filtered movies
	mt := &MovieTransformer{Movies: filteredMovies}
	if err := mt.Execute(); err != nil {
		log.Fatalf("transformation error: %v", err)
	}

	// Sort the filtered movies in reverse lexicographical order of ratings.
	sort.Slice(filteredMovies, func(i, j int) bool {
		return filteredMovies[i].Rating > filteredMovies[j].Rating
	})

	fmt.Println("Movies with ratings between 3.0 and 5.0 (in reverse order):")
	for _, movie := range filteredMovies {
		fmt.Println(movie)
	}
}
