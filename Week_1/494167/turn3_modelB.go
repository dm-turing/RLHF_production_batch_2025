package main

import (
   // ... (Same imports as before)
   "github.com/go-redis/redis/v8"
)

const (
   secretKey        = "your_super_secret_key"
   tokenExpiration  = time.Hour * 24
   tokenRotationTTL = time.Hour * 24 * 7 // Rotate tokens every week
   redisAddr        = "redis:6379"
   redisPassword    = ""
   redisDB          = 0
)

var (
   rd *redis.Client
)

func init() {
   rd = redis.NewClient(&redis.Options{
      Addr:     redisAddr,
      Password: redisPassword,
      DB:       redisDB,
   })
   if _, err := rd.Ping().Result(); err != nil {
      log.Fatalf("Error connecting to Redis: %v", err)
   }
}

func generateToken(user User) string {
   // ... (Same code as before)
}

func verifyToken(r *http.Request) (string, error) {
   // ... (Same code as before)
}

func handleTokenRotation(w http.ResponseWriter, r *http.Request) {
   userEmail, err := verifyToken(r)
   if err != nil {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
   }

   // Check if token rotation is required
   if isTokenRotationRequired(userEmail) {
      newToken, err := generateAndStoreNewToken(userEmail)
      if err != nil {
         http.Error(w, "Internal Server Error", http.StatusInternalServerError)
         return
      }

      // Return the new token to the client
      response := map[string]string{"message": "Token rotated", "token": newToken}
      json.NewEncoder(w).Encode(response)
      return
   }

   // If token rotation is not required, proceed with the original handler
   handler(w, r)
}

func isTokenRotationRequired(userEmail string) bool {
   // Check if the token has been used within the last rotation period
   lastRotationTime, err := rd.Get(userEmail + "-last-rotation").Result()
   if err == redis.Nil {
      // Token has never been rotated, so rotate it now
      return true
   }
   if err != nil {
      // Error retrieving rotation time, handle it appropriately
      log.Printf("Error retrieving rotation time for %s: %v", userEmail, err)
      return false
   }

   rotationTime, err := time.Parse(time.RFC3339, lastRotationTime)
   if err != nil {
      // Error parsing rotation time, handle it appropriately
      log.Printf("Error parsing rotation time for %s: %v", userEmail, err)
      return false
   }

   return time.Since(rotationTime) >= tokenRotationTTL
}

func generateAndStoreNewToken(userEmail string) (string, error) {
   user := User{Email: userEmail} // Retrieve user details from database or cache
   newToken := generateToken(user)

   // Store the new token in the Redis cache
   if err := rd.Set(userEmail+"-token", newToken, tokenExpiration).Err(); err != nil {
      return "", fmt.Errorf("error storing new token in Redis: %v", err)
   }

   // Update the last rotation time
   if err := rd.Set(userEmail+"-last-rotation", time.Now().Format(time.RFC3339), 0).Err(); err != nil {
      return "", fmt.Errorf("error updating last rotation time in Redis: %v", err)
   }

   return newToken, nil