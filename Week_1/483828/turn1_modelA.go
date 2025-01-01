package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Feedback represents a single feedback entry
type Feedback struct {
	ID     int       `json:"id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
	Rating int       `json:"rating"`
}

var feedbacks []Feedback
var nextID int = 1

// AddFeedback adds a new feedback entry
func AddFeedback(w http.ResponseWriter, r *http.Request) {
	var f Feedback
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f.ID = nextID
	f.Date = time.Now()
	nextID++
	feedbacks = append(feedbacks, f)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}

// GetAllFeedbacks returns all feedback entries
func GetAllFeedbacks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(feedbacks)
}

// AverageRating calculates the average rating of feedbacks
func AverageRating() float64 {
	if len(feedbacks) == 0 {
		return 0.0
	}
	sum := 0
	for _, f := range feedbacks {
		sum += f.Rating
	}
	return float64(sum) / float64(len(feedbacks))
}

func main() {
	http.HandleFunc("/feedback", AddFeedback)
	http.HandleFunc("/feedbacks", GetAllFeedbacks)
	fmt.Println("Feedback server running on :8080")
	http.ListenAndServe(":8080", nil)
}
