package main

import (
	"fmt"
	"log"

	"github.com/qe-team/marmot"
)

func main() {
	// Load a Marmot model (e.g., English language model)
	model, err := marmot.LoadModel("en")
	if err != nil {
		log.Fatal(err)
	}

	// Example feedback text
	feedbackText := "This product is fantastic! I love it."

	// Analyze sentiment of the feedback text
	sentiment, err := model.Sentiment(feedbackText)
	if err != nil {
		log.Fatal(err)
	}

	// Print the sentiment score (range: -1 to 1, where -1 is negative and 1 is positive)
	fmt.Println("Sentiment score:", sentiment)
}
