package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const secretKey = "your_super_secret_key"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func generateToken(user User) string {
	claims := &jwt.StandardClaims{
		Subject:   user.Email,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Fatalf("Error generating token: %v", err)
	}
	return signedToken
}

func verifyToken(r *http.Request) (string, error) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		return "", fmt.Errorf("no token found in query parameters")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", fmt.Errorf("invalid token")
}

func handler(w http.ResponseWriter, r *http.Request) {
	userEmail, err := verifyToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]string{"message": "Hello, " + userEmail}
	json.NewEncoder(w).Encode(response)
}

func main() {
	user := User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
	token := generateToken(user)
	fmt.Println("Generated token:", token)

	http.HandleFunc("/protected", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
