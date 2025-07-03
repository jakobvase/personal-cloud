package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Response represents the API response structure
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/api/test", testHandler).Methods("GET")
	r.HandleFunc("/api/health", healthHandler).Methods("GET")

	// Start the server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Available endpoints:\n")
	fmt.Printf("  GET  http://localhost%s/\n", port)
	fmt.Printf("  GET  http://localhost%s/api/test\n", port)
	fmt.Printf("  GET  http://localhost%s/api/health\n", port)
	
	log.Fatal(http.ListenAndServe(port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Go Server Home</title>
	</head>
	<body>
		<h1>Welcome to the Go Server!</h1>
		<p>This is the home page served as HTML.</p>
	</body>
	</html>
	`
	
	w.Write([]byte(html))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := Response{
		Message:   "This is a test API call! The server is working correctly.",
		Timestamp: time.Now(),
		Status:    "success",
	}
	
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "go-server",
	}
	
	json.NewEncoder(w).Encode(response)
} 