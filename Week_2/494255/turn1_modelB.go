package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Data record struct
type Movie struct {
	Title       string `json:"title" csv:"title"`
	Genre       string `json:"genre" csv:"genre"`
	ReleaseDate string `json:"release_date" csv:"release_date"`
}

// Extractor interface
type Extractor interface {
	Extract() (<-chan Movie, error)
}

type csvExtractor struct {
	file *os.File
}

// NewCSVExtractor creates a new CSV extractor
func NewCSVExtractor(file *os.File) Extractor {
	return &csvExtractor{file}
}

// Extract movies from CSV file
func (e *csvExtractor) Extract() (<-chan Movie, error) {
	movieChan := make(chan Movie)
	go func() {
		defer close(movieChan)
		defer e.file.Close()

		r := csv.NewReader(e.file)
		r.FieldsPerRecord = 3

		// Skip header row
		_, err := r.Read()
		if err != nil {
			log.Fatalf("Error reading CSV header: %v", err)
		}

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading CSV record: %v", err)
				continue
			}

			m := Movie{
				Title:       record[0],
				Genre:       record[1],
				ReleaseDate: record[2],
			}
			movieChan <- m
		}
	}()
	return movieChan, nil
}

type jsonExtractor struct {
	file *os.File
}

// NewJSONExtractor creates a new JSON extractor
func NewJSONExtractor(file *os.File) Extractor {
	return &jsonExtractor{file}
}

// Extract movies from JSON file
func (e *jsonExtractor) Extract() (<-chan Movie, error) {
	movieChan := make(chan Movie)
	go func() {
		defer close(movieChan)
		defer e.file.Close()

		scanner := bufio.NewScanner(e.file)
		for scanner.Scan() {
			var m Movie
			err := json.Unmarshal(scanner.Bytes(), &m)
			if err != nil {
				log.Printf("Error unmarshalling JSON record: %v", err)
				continue
			}
			movieChan <- m
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading JSON file: %v", err)
		}
	}()
	return movieChan, nil
}

// Dynamic data extractor
func ExtractMovies(inputFile string) (<-chan Movie, error) {
	// Detect file format based on extension
	ext := strings.ToLower(strings.Split(inputFile, ".")[1])
	var extractor Extractor
	var err error

	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	switch ext {
	case "csv":
		extractor = NewCSVExtractor(file)
	case "json":
		extractor = NewJSONExtractor(file)
	default:
		file.Close()
		return nil, fmt.Errorf("unsupported file format: %s", ext)