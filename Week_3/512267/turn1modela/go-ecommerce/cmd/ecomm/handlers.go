package ecomm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go-ecommerce/lib/db"
	"go-ecommerce/lib/redis"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	ConversionRate float32 `json:"conversion_rate"`
	SessionDuration float32 `json:"session_duration"`
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]User, 0)

	// Read users from Redis cache
	cacheKey := "users"
	b, err := redis.Get(r.Context(), cacheKey)
	if err != nil || b == nil {
		log.Printf("Cache miss for key: %s", cacheKey)
		users, err = db.GetAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Store users in Redis cache
		err = redis.Set(r.Context(), cacheKey, users)
		if err != nil {
			log.Printf("Error caching users: %v", err)
		}
	} else {
		err = json.Unmarshal(b, &users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	// Check cache first
	userCacheKey := fmt.Sprintf("user:%s", userId)
	b, err := redis.Get(r.Context(), userCacheKey)
	if err != nil || b == nil {
		log.Printf("Cache miss for key: %s", userCacheKey)
		user, err := db.GetUserById(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Store user in Redis cache
		b, err = json.Marshal(user)
		if err != nil {
			log.Printf("Error caching user: %v", err)
		}
		err = redis.Set(r.Context(), userCacheKey, b)
		if err != nil {
			log.Printf("Error caching user: %v", err)
		}
	} else {
		err = json.Unmarshal(b, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate random engagement metrics
	user.ConversionRate = float32(rand.Float64())*100
	user.SessionDuration = float32(rand.Float64())*60

	err := db.AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Clear cache for user key
	userCacheKey := fmt.Sprintf("user:%d", user.ID)
	_, err = redis.Del(r.Context(), userCacheKey)
	if err != nil {
		log.Printf("Error deleting user cache: %v", err)
	}

	// Clear list of users in cache
	_, err = redis.Del(r.Context(), "users")
	if err != nil {
		log.Printf("Error deleting users cache: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
