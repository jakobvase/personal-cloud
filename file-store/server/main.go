package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	// Check for session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Check if session exists in userSessions
	username, ok := userSessions[cookie.Value]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Optionally, check if cookie is expired (handled by browser, but double-check for safety)
	if cookie.Expires.Before(time.Now()) {
		delete(userSessions, cookie.Value)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

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
	// For demonstration, accept any username/password
	if username == "" || password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}
	// Generate a simple session ID (not secure, for demo only)
	sessionID := fmt.Sprintf("session-%d", time.Now().UnixNano())
	userSessions[sessionID] = username
	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires: time.Now().Add(24 * time.Hour),
		// Secure: true, // Uncomment if using HTTPS
	})
	// Redirect to home
	http.Redirect(w, r, "/", http.StatusFound)
} 