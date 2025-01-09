package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

type Movie struct {
	Title  string   `json:"title"`
	Year   int      `json:"year"`
	Genre  []string `json:"genre"`
	Rating float64  `json:"rating"`
}

type MovieDistance struct {
	MovieType Movie   `json: "moviee"`
	Distance  float64 `json: "distance"`
}

func loadMovies() ([]Movie, error) {
	file, err := os.Open("movies.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var movies []Movie
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func euclideanDistance(movie1, movie2 Movie) float64 {
	distance := 0.0
	distance += math.Pow(float64(movie1.Year-movie2.Year), 2)
	for i := 0; i < len(movie1.Genre); i++ {
		if i >= len(movie2.Genre) || movie1.Genre[i] != movie2.Genre[i] {
			distance += math.Pow(float64(1), 2)
		}
	}
	distance += math.Pow(movie1.Rating-movie2.Rating, 2)
	return math.Sqrt(distance)
}

func knnRecommend(movies []Movie, target Movie, k int) ([]Movie, error) {
	if k <= 0 {
		return nil, fmt.Errorf("k must be greater than 0")
	}

	distances := make([]*MovieDistance, len(movies))
	for i, movie := range movies {
		distances[i] = &MovieDistance{
			MovieType: movie,
			Distance:  euclideanDistance(movie, target),
		}
	}

	sort.Slice(distances, func(i, j int) bool {
		return distances[i].Distance < distances[j].Distance
	})

	recommendations := make([]Movie, k)
	for i := 0; i < k; i++ {
		if i < len(distances) {
			recommendations[i] = distances[i].MovieType
		}
	}

	return recommendations, nil
}

func main() {
	movies, err := loadMovies()
	if err != nil {
		log.Fatalf("Error loading movies: %v", err)
	}

	target := Movie{
		Title:  "",
		Year:   2019,
		Genre:  []string{"Sci-Fi", "Action"},
		Rating: 8.5,
	}

	k := 5
	recommendations, err := knnRecommend(movies, target, k)
	if err != nil {
		log.Fatalf("Error making recommendations: %v", err)
	}

	fmt.Println("Recommended Movies:")
	for _, movie := range recommendations {
		fmt.Printf("%s (%d), Genres: %s, Rating: %.1f\n", movie.Title, movie.Year, movie.Genre, movie.Rating)
	}
}
