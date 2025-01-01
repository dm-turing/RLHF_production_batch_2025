package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Feedback represents a single feedback entry
type Feedback struct {
	ID     int       `json:"id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
	Rating int       `json:"rating"`
}

var (
	feedbacks     map[int]Feedback
	nextID        int = 1
	totalRating   int = 0
	feedbackCount int = 0
	mutex         sync.RWMutex
)

func init() {
	feedbacks = make(map[int]Feedback)
}

// AddFeedback adds a new feedback entry
func AddFeedback(w http.ResponseWriter, r *http.Request) {
	var f Feedback
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	f.ID = nextID
	f.Date = time.Now()
	nextID++
	totalRating += f.Rating
	feedbackCount++
	feedbacks[f.ID] = f
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}

// GetAllFeedbacks returns all feedback entries
func GetAllFeedbacks(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()
	allFeedbacks := make([]Feedback, 0, len(feedbacks))
	for _, f := range feedbacks {
		allFeedbacks = append(allFeedbacks, f)
	}
	json.NewEncoder(w).Encode(allFeedbacks)
}

// DeleteFeedback deletes a feedback entry by ID
func DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid feedback ID", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	f, ok := feedbacks[id]
	if !ok {
		http.Error(w, "Feedback not found", http.StatusNotFound)
		return
	}
	delete(feedbacks, id)
	totalRating -= f.Rating
	feedbackCount--
	w.WriteHeader(http.StatusNoContent)
}

// AverageRating calculates the average rating of feedbacks
func AverageRating() float64 {
	mutex.RLock()
	defer mutex.RUnlock()
	if feedbackCount == 0 {
		return 0.0
	}
	return float64(totalRating) / float64(feedbackCount)
}

func main() {
	http.HandleFunc("/feedback", AddFeedback)
	http.HandleFunc("/feedbacks", GetAllFeedbacks)
	http.HandleFunc("/delete", DeleteFeedback)
	fmt.Println("Feedback server running on :8080")
	http.ListenAndServe(":8080", nil)
}
