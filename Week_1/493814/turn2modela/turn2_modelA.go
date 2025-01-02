package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string `json:"environment"`
	DatabaseURL string `json:"databaseURL"`
	FeatureFlag bool   `json:"featureFlag"`
}

var config Config

func init() {
	// Load configuration from YAML file based on environment variable
	viper.SetConfigName("config")              // name of config file (without extension)
	viper.SetConfigType("yaml")                // set config file type
	viper.AddConfigPath(".")                   // add current directory to search path
	viper.AddConfigPath("$HOME/.config/myapp") // add home directory to search path
	viper.AddConfigPath("/etc/myapp")          // add system-wide config path

	env := os.Getenv("ENV") // Get the environment variable
	if env != "" {
		viper.SetConfigName(fmt.Sprintf("config.%s", env)) // Look for environment-specific config file
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; use default config.
			fmt.Println("Config file not found, using default config.")
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Fatal error config file: %s \n", err)
		}
	}

	// Unmarshal YAML configuration into Config struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into config struct, %v", err)
	}
}

func getConfigByEnv(c *gin.Context) {
	envVar := os.Getenv("ENVIRONMENT")
	if envVar == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Environment variable 'ENVIRONMENT' not set"})
		return
	}
	c.JSON(http.StatusOK, config)
}

func main() {
	router := gin.Default()
	router.GET("/config", getConfigByEnv)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}
