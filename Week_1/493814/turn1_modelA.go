package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Config represents a configuration block
type Config struct {
	Environment string `json:"environment"`
	DatabaseURL string `json:"databaseURL"`
	FeatureFlag bool   `json:"featureFlag"`
}

var configStore = map[string]Config{}

func loadConfigs() {
	if err := ioutil.ReadFile("configs/dev.yaml"); err != nil {
		fmt.Println("Error loading dev config:", err)
		return
	}
	configStore["dev"] = parseConfig(`{
		"environment": "dev",
		"databaseURL": "postgres://dev:dev@localhost:5432/devdb",
		"featureFlag": false
	}`)

	if err := ioutil.ReadFile("configs/staging.yaml"); err != nil {
		fmt.Println("Error loading staging config:", err)
		return
	}
	configStore["staging"] = parseConfig(`{
		"environment": "staging",
		"databaseURL": "postgres://staging:staging@staging-db.example.com:5432/stagingdb",
		"featureFlag": true
	}`)

	if err := ioutil.ReadFile("configs/prod.yaml"); err != nil {
		fmt.Println("Error loading prod config:", err)
		return
	}
	configStore["prod"] = parseConfig(`{
		"environment": "prod",
		"databaseURL": "postgres://prod:prod@prod-db.example.com:5432/proddb",
		"featureFlag": true
	}`)
}

func parseConfig(raw string) Config {
	var config Config
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		panic(err)
	}
	return config
}

func getConfigByEnv(c *gin.Context) {
	envVar := os.Getenv("ENVIRONMENT")
	if envVar == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Environment variable 'ENVIRONMENT' not set"})
		return
	}

	config, exists := configStore[envVar]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration for environment '" + envVar + "' not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

func main() {
	loadConfigs()

	router := gin.Default()
	router.GET("/config", getConfigByEnv)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
