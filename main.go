package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
)

var (
	signingKey                = []byte("super-secret-key") // Should be in an environment variable for security
	refreshKey                = []byte("refresh-secret-key")
	accessTokenExpiration     = time.Hour * 1  // 1 hour
	refreshTokenExpiration    = time.Hour * 24 // 24 hours
	maxFailedAttempts         = 3
	maxFailedAttemptsDuration = time.Minute * 5
	r                         = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Replace with your Redis address
	})
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func generateAccessToken(user User) string {
	claims := &jwt.StandardClaims{
		Subject:   user.Email,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(accessTokenExpiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		log.Fatalf("Error generating access token: %v", err)
	}
	return signedToken
}

func generateRefreshToken(user User) string {
	claims := &jwt.StandardClaims{
		Subject:   user.Email,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(refreshTokenExpiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(refreshKey)
	if err != nil {
		log.Fatalf("Error generating refresh token: %v", err)
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
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", fmt.Errorf("invalid token")
}

func verifyRefreshToken(r *http.Request) (string, error) {
	tokenStr := r.URL.Query().Get("refresh")
	if tokenStr == "" {
		return "", fmt.Errorf("no refresh token found in query parameters")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return refreshKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", fmt.Errorf("invalid refresh token")
}

func rateLimitCheck(ip string) bool {
	val, err := r.Get(ip).Result()
	if err != nil {
		return true
	}
	// defer val.Close()

	count, err := strconv.Atoi(val)
	if err != nil {
		return true
	}

	if count >= maxFailedAttempts {
		return true
	}

	return false
}

func recordFailedAttempt(ip string) {
	res, err := r.Incr(ip).Result()
	if err != nil {
		log.Println(err)
		// log.Error(err)
	}
	err = r.Set(ip, fmt.Sprintf("%d", 1+res), time.Minute*5).Err()
	if err != nil {
		log.Println(err)
		// log.Error(err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	if rateLimitCheck(ip) {
		http.Error(w, "Too many failed login attempts. Try again later.", http.StatusTooManyRequests)
		return
	}

	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if email == "" || password == "" {
		http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
		recordFailedAttempt(ip)
		return
	}

	// Placeholder for user authentication logic
	user, err := authenticateUser(email, password)
	if err != nil {
		http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
		recordFailedAttempt(ip)
		return
	}

	accessToken := generateAccessToken(user)
	refreshToken := generateRefreshToken(user)

	response := map[string]string{"access_token": accessToken, "refresh_token": refreshToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func authenticateUser(email, password string) (User, error) {
	// Placeholder for your user authentication logic
	if email == "john.doe@example.com" && password == "secret" {
		return User{ID: 1, Name: "John Doe", Email: email}, nil
	}
	return User{}, fmt.Errorf("authentication failed")
}

func handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	userEmail, err := verifyRefreshToken(r)
	if err != nil {
		http.Error(w, "Invalid refresh token.", http.StatusUnauthorized)
		return
	}

	user, err := authenticateUser(userEmail, "") // No password needed for refresh
	if err != nil {
		http.Error(w, "Invalid refresh token.", http.StatusUnauthorized)
		return
	}

	accessToken := generateAccessToken(user)
	response := map[string]string{"access_token": accessToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/refresh", handleRefreshToken)
	http.HandleFunc("/protected", handleProtected)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleProtected(w http.ResponseWriter, r *http.Request) {
	userEmail, err := verifyToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", `Bearer realm="Example App"`)
		return
	}

	response := map[string]string{"message": "Hello, " + userEmail}
	json.NewEncoder(w).Encode(response)
}
