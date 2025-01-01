package main

import "fmt"

type Rating int

const (
	Excellent Rating = iota + 1
	Good
	Average
	Poor
	Terrible
)

type Feedback struct {
	Comment string
	Rating  Rating
}
type FeedbackCollection struct {
	feedbacks    []Feedback
	ratingCounts [Terrible + 1]int
}

func (r Rating) String() string {
	switch r {
	case Excellent:
		return "Excellent"
	case Good:
		return "Good"
	case Average:
		return "Average"
	case Poor:
		return "Poor"
	case Terrible:
		return "Terrible"
	default:
		return "Invalid Rating"
	}
}

// Initialize a new FeedbackCollection
func NewFeedbackCollection() *FeedbackCollection {
	return &FeedbackCollection{
		feedbacks:    make([]Feedback, 0),
		ratingCounts: [Terrible + 1]int{},
	}
}

// AddFeedback adds a new feedback entry to the collection
func (fc *FeedbackCollection) AddFeedback(feedback Feedback) {
	fc.feedbacks = append(fc.feedbacks, feedback)
	fc.ratingCounts[feedback.Rating]++
}

// PrintFeedbacks prints all the collected feedback
func (fc *FeedbackCollection) PrintFeedbacks() {
	for _, f := range fc.feedbacks {
		fmt.Printf("Comment: %s, Rating: %s\n", f.Comment, f.Rating.String())
	}
}

// PrintRatingStatistics prints the rating statistics
func (fc *FeedbackCollection) PrintRatingStatistics() {
	totalFeedback := len(fc.feedbacks)
	if totalFeedback == 0 {
		fmt.Println("No feedback collected yet.")
		return
	}

	for rating, count := range fc.ratingCounts {
		fmt.Printf("%s: %d (%.2f%%)\n", rating, count, float64(count)*100/float64(totalFeedback))
	}
}

func main() {
	fc := NewFeedbackCollection()

	f1 := Feedback{"Excellent service!", Excellent}
	fc.AddFeedback(f1)

	f2 := Feedback{"The food is average.", Average}
	fc.AddFeedback(f2)

	// Printing the feedbacks
	fc.PrintFeedbacks()

	// Printing the rating statistics
	fmt.Println("\nRating Statistics:")
	fc.PrintRatingStatistics()
}
