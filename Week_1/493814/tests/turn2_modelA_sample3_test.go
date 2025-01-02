package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigServerValidation(t *testing.T) {
	client := &http.Client{}
	resp, err := client.Get("http://localhost:8080/config")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var config MockConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		t.Fatal(err)
	}

	// Assert the environment variable is being used
	assert.Equal(t, os.Getenv("ENVIRONMENT"), config.Environment)
}
