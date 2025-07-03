package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"server/sessions"

	"github.com/gorilla/mux"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// In-memory user session store (for demonstration only)
var userSessions = make(map[string]string) // sessionID -> username

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
    
	// API routes
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


func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "go-server",
	}
	
	json.NewEncoder(w).Encode(response)
}

func htmlHead(title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>%s</title>
</head>
`, title)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	username, err := sessions.GetUser(r.Cookie("session"))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := htmlHead("Go Server Home") + `
	<body>
		<h1>Welcome to the Go Server!</h1>
		<p>This is the home page served as HTML.</p>
	</body>
	</html>
	`

	w.Write([]byte(html))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := htmlHead("File Store Login") + `
	<body>
		<h1>Login</h1>
		<form method="POST" action="/login">
			<label for="username">Username:</label><br>
			<input type="text" id="username" name="username"><br><br>
			<label for="password">Password:</label><br>
			<input type="password" id="password" name="password"><br><br>
			<input type="submit" value="Login">
		</form>
	</body>
	</html>
	`

	w.Write([]byte(html))
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	cookie, err := sessions.SetSession(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set session cookie
	http.SetCookie(w, cookie)
	// Redirect to home
	http.Redirect(w, r, "/", http.StatusFound)
} 