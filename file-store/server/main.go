package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/sessions"
	"server/user_storage"
	"server/users"
	"time"

	"github.com/gorilla/mux"
)

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
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")

	// API routes
	r.HandleFunc("/api/health", healthHandler).Methods("GET")

	// Start the server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)

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
	id, err := sessions.GetUserId(r.Cookie("session"))

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	var _, err2 = users.GetUser(id)

	if err2 != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := htmlHead("File storage example") + `
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
			<label for="password">Password - but this is a test application, use a test password:</label><br>
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

	// First try to login the user
	user, err := users.LoginUser(username, password)

	// If login fails, try to create the user
	if err != nil {
		user, err = users.CreateUser(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}

	cookie, err := sessions.SetSession(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set session cookie
	http.SetCookie(w, cookie)
	// Redirect to home
	http.Redirect(w, r, "/", http.StatusFound)
}
