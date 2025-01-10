package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type RequestData struct {
	State            string  `json:"state"`
	RequestRate      float64 `json:"request_rate"`
	NetworkBandwidth float64 `json:"network_bandwidth"`
	MimeType         string  `json:"mime_type"`
}

type Response struct {
	Message string   `json:"message"`
	Access  []string `json:"access"`
}

func corsMiddleware(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var allowedOrigins []string

		// Dummy logic to simulate ML predictions
		allowedOrigins = getDynamicCORS(r)

		w.Header().Set("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ","))
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

		next.ServeHTTP(w, r)
	}
}

func getDynamicCORS(r *http.Request) []string {
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		fmt.Println("Error decoding request:", err)
		return []string{}
	}

	fmt.Println("Request Data:", requestData, requestData.State == "idle")

	// Dummy logic to simulate a decision-making process similar to what an ML model might do
	if requestData.State == "active" && requestData.RequestRate < 10.0 && requestData.NetworkBandwidth > 5.0 {
		return []string{"http://example.com", "https://another-example.com"}
	} else if requestData.State == "idle" || requestData.RequestRate > 50.0 {
		fmt.Println([]string{"http://example.com"})
		return []string{"http://example.com"}
	} else {
		return []string{"https://another-example.com"}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// allowedOrigins := getDynamicCORS(r)
	allowedOrigins := []string{"http://example.com", "https://another-example.com"}

	if allowedOrigins == nil || len(allowedOrigins) == 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Request not allowed"))
		return
	}

	response := Response{
		Message: "Access granted!",
		Access:  allowedOrigins,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("Error encoding response:", err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/ml/cors", corsMiddleware(http.HandlerFunc(handler))).Methods("POST")

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}
