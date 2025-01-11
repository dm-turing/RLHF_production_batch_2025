package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

// CSVFormatDetector detects CSV format.
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

// FormatExtractor extracts data from a file based on its format.
type FormatExtractor interface {
	Extract(io.Reader, string) ([]Movie, error)
}

// CSVFormatExtractor extracts data in CSV format.
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

// DataExtractor extracts movie recommendations from a file.
type DataExtractor struct {
	detector  FormatDetector
	extractor FormatExtractor
}

func NewDataExtractor() *DataExtractor {
	return &DataExtractor{
		detector:  &CSVFormatDetector{},
		extractor: &CSVFormatExtractor{},
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

func main() {
	filePath := "movies.json" // Example file path
	extractor := NewDataExtractor()

	// You can specify different detectors and extractors if needed
	extractor.SetDetector(&JSONFormatDetector{})
	extractor.SetExtractor(&JSONFormatExtractor{})

	movies, err := extractor.Extract(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, movie := range movies {
		fmt.Println(movie)
	}
}
