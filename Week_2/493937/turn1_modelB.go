package main

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Movie struct {
	ID     int
	Title  string
	Genres []string
}
type User struct {
	ID      int
	Ratings map[int]int
}
type KNNRecommender struct {
	K      int
	Movies []Movie
	Users  []User
}

func (r *KNNRecommender) loadData(moviesFile string, ratingsFile string) {
	// Load movie data
	r.loadMovies(moviesFile)
	// Load ratings data and update users' ratings
	r.loadRatings(ratingsFile)
}
func (r *KNNRecommender) loadMovies(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records[1:] {
		id, _ := strconv.Atoi(record[0])
		title := record[1]
		genres := strings.Split(record[2], "|")
		r.Movies = append(r.Movies, Movie{ID: id, Title: title, Genres: genres})
	}
}
func (r *KNNRecommender) loadRatings(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records[1:] {
		userID, _ := strconv.Atoi(record[0])
		movieID, _ := strconv.Atoi(record[1])
		rating, _ := strconv.Atoi(record[2])
		r.updateUserRating(userID, movieID, rating)
	}
}
func (r *KNNRecommender) updateUserRating(userID int, movieID int, rating int) {
	// Find the user or create a new one
	var user *User
	for i, u := range r.Users {
		if u.ID == userID {
			user = &r.Users[i]
			break
		}
	}
	if user == nil {
		user = &User{ID: userID, Ratings: make(map[int]int)}
		r.Users = append(r.Users, *user)
	}
	// Update the user's rating for the movie
	user.Ratings[movieID] = rating
}
func (r *KNNRecommender) calculateDistance(user1 User, user2 User) float64 {
	sum := 0.0
	count := 0
	// Find common movies rated by both users
	for movieID1, rating1 := range user1.Ratings {
		if rating2, ok := user2.Ratings[movieID1]; ok {
			sum += math.Pow(float64(rating1-rating2), 2)
			count++
		}
	}
	if count == 0 {
		// Return a large distance if there are no common movies
		return math.MaxFloat64
	}
	return math.Sqrt(sum) / float64(count)
}
