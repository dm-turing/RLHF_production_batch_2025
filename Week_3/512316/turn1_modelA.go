package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/casbin/casbin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const (
	secretKey      = "this_is_a_secret_key"
	adminRole      = "admin"
	userRole       = "user"
	sqliteFilePath = "resource.db"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Project struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
	Status string `json:"status"`
}

var (
	users    map[string]string
	db       *sql.DB
	enforcer *casbin.Enforcer
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", sqliteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE, password TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS projects (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, user_id INTEGER, status TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	users = make(map[string]string)

	enforcer = casbin.NewEnforcer("./example.conf", "./policy.csv")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	users[user.Username] = user.Password

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if password, ok := users[user.Username]; ok && password == user.Password {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = user.Username
		claims["role"] = userRole
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			http.Error(w, "Error signing token", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: tokenString,
		})
		fmt.Println("token = ", tokenString) // Only for demonstration purposes.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged in successfully"))
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "No token found", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		tok, err := jwt.ParseWithClaims(token.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tok.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createProject(w http.ResponseWriter, r *http.Request) {
	var project Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user := r.Context().Value("user").(jwt.Claims).(jwt.MapClaims)

	// Check permissions
	sub := user["username"] // The user
	obj := project.Name     // The resource
	act := "create_project" // The action

	// Check if the user has permission
	if !enforcer.Enforce(sub, obj, act) {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	_, err = db.Exec("INSERT INTO projects (name, user_id, status) VALUES (?, ?, ?)", project.Name, project.UserID, project.Status)
	if err != nil {
		http.Error(w, "Error creating project", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Project created successfully"))
}

func main() {
	// defer db.Close()

	// Setting up a signal handler
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		log.Println("Shutting down the database instance gracefully...")
		db.Close() // Close the database connection
		os.Exit(0)
	}()
	r := mux.NewRouter()
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/login", loginUser).Methods("POST")
	r.Handle("/projects", authorize(http.HandlerFunc(createProject))).Methods("POST")
	fmt.Println("Starting localhost on 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
	select {}
}
