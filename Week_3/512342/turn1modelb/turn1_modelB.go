package main  
import (  
    "fmt"
    "log"
    "net/http"
)

// ErrorLogger is a middleware function that takes an http.Handler
// and returns an http.Handler. It's responsible for logging errors
// based on user input.
type ErrorLogger func(http.Handler) http.Handler

// ComposeErrorLogs takes a variable number of ErrorLogger middleware
// functions and returns a new ErrorLogger that chains them together.
func ComposeErrorLogs(middleware ...ErrorLogger) ErrorLogger {
    return func(h http.Handler) http.Handler {
        for _, m := range middleware {
            h = m(h)
        }
        return h
    }
}

// LogUserInput is an ErrorLogger middleware that logs the user's input
// if it's not empty.
func LogUserInput(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        userInput := r.FormValue("userInput")
        if userInput != "" {
            log.Printf("User input: %s", userInput)
        }

        next.ServeHTTP(w, r)
    })
}

// LogErrorMessage is an ErrorLogger middleware that logs the error message
// if it's not empty.
func LogErrorMessage(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        errorMessage := r.FormValue("errorMessage")
        if errorMessage != "" {
            log.Printf("Error message: %s", errorMessage)
        }

        next.ServeHTTP(w, r)
    })
}

// main is the entry point of the application.
func main() {
    // Create an HTTP handler for the root path ("/").
    rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, world!")
    })

    // Compose the error loggers using the ComposeErrorLogs function.
    errorLogger := ComposeErrorLogs(LogUserInput, LogErrorMessage)

    // Wrap the root handler with the error logger.
    http.Handle("/", errorLogger(rootHandler))

    // Start the HTTP server on port 8080.
    log.Println("Server starting on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

