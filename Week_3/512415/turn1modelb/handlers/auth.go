package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	jwtSecret = "your-secret-key" // Replace this with a strong secret key
)

var users = map[string]string{
	"user": "password",
}

// Token represents a JWT token
type Token struct {
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}

// CreateToken generates a new JWT token for the given user ID
func CreateToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	return token.SignedString([]byte(jwtSecret))
}

// ParseToken parses and validates the JWT token from the request
func ParseToken(r *http.Request) (*Token, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("Authorization token missing")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &Token{
			UserID:    claims["user_id"].(string),
			ExpiresAt: claims["exp"].(int64),
		}, nil
	}

	return nil, fmt.Errorf("Invalid token")
}

// Authenticate authenticates the user using JWT token
func Authenticate(w http.ResponseWriter, r *http.Request) bool {
	token, err := ParseToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return false
	}

	// Validate the token expiration
	if time.Now().Unix() > token.ExpiresAt {
		http.Error(w, "Token expired", http.StatusUnauthorized)
		return false
	}

	// For simplicity, we'll assume user IDs are unique
	_, ok := users[token.UserID]
	if !ok {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return false
	}

	return true
}
