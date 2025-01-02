package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockConfig struct {
	Environment string `json:"environment"`
	DatabaseURL string `json:"databaseURL"`
	FeatureFlag bool   `json:"featureFlag"`
}

func TestConfigServerEnvRetrieval(t *testing.T) {
	client := &http.Client{}
	resp, err := client.Get("http://localhost:8080/config")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var config MockConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		t.Fatal(err)
	}

	// Assert the environment variable is being used
	assert.Equal(t, os.Getenv("ENVIRONMENT"), config.Environment)
}