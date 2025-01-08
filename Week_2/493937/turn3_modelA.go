package main

import (
	"fmt"
	"log"
)

type Movie struct {
	Title  string   `json:"title"`
	Year   int      `json:"year"`
	Genre  []string `json:"genre"`
	Rating float64  `json:"rating"`
}

type Recommender interface {
	Train([]Movie, []Movie) error
	Recommend(Movie) ([]Movie, error)
}

func (r *Recommender) handleError(err error) error {
	switch err {
	case fmt.Errorf("could not open file: %w", _):
		return fmt.Errorf("could not load movie data")
	case fmt.Errorf("error parsing JSON: %w", _):
		return fmt.Errorf("error parsing movie data")
	case fmt.Errorf("k must be greater than 0"):
		return fmt.Errorf("invalid number of neighbors")
	case fmt.Errorf("no movies available"):
		return fmt.Errorf("no movies in the database")
	default:
		return fmt.Errorf("an unexpected error occurred")
	}
}

func (r *Recommender) Train(movies []Movie, ratings []Movie) error {
	err := r.model.Train(movies, ratings)
	if err != nil {
		return r.handleError(err)
	}
	return nil
}

func (r *Recommender) Recommend(target Movie) ([]Movie, error) {
	recommendations, err := r.model.Recommend(target)
	if err != nil {
		return nil, r.handleError(err)
	}
	return recommendations, nil
}

func main() {
	r := Recommender()
	log.Printf("Model training started.")
	movies, ratings := make([]Movie, 0), make([]Movie, 0)
	err := r.model.Train(movies, ratings)
	if err != nil {
		log.Printf("Model training failed: %v", err)
		return r.handleError(err)
	}
	log.Printf("Model training completed.")
}
