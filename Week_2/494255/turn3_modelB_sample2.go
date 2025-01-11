package main  
import (  
        "encoding/json"
        "fmt"
        "log"
        "os"
        "time"

        "github.com/go-redis/redis"
)

func loadToRedisJSON(movies []Movie, redisClient *redis.Client) {
        redisKey := "movies_json"

        // Marshal the movie data to JSON
        jsonData, err := json.Marshal(movies)
        if err != nil {
                log.Fatal(err)
        }

        // Set the JSON data to Redis
        if err := redisClient.Set(redisKey, jsonData, 0).Err(); err != nil {
                log.Fatal(err)
        }