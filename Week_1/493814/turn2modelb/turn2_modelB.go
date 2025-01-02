package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
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

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		// Marshal configuration back to YAML and send as response
		data, err := yaml.Marshal(config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	port := config.Server.Port
	fmt.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
