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
	r.HandleFunc("/storage/authorize", authorizeStorageHandler).Methods("GET")
	r.HandleFunc("/storage/oauth-callback", getStorageTokenHandler).Methods("GET")

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

func authorizeStorageHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.GetUserId(r.Cookie("session"))
	user, err2 := users.GetUser(userId)

	if err != nil || err2 != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	// Generate a state parameter for CSRF protection
	// TODO probably unsafe.
	state := fmt.Sprintf("%d", time.Now().UnixNano())
	// Store it on the user
	user.SetStorageAuthorizationPkce(state)

	redirectURI := "http://localhost:8080/storage/oauth-callback"

	// Generate the authorization URL
	authURL := user_storage.GenerateAuthUrl(redirectURI, state)

	// Redirect the user to the authorization URL
	http.Redirect(w, r, authURL, http.StatusFound)
}

func getStorageTokenHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.GetUserId(r.Cookie("session"))
	user, err2 := users.GetUser(userId)

	if err != nil || err2 != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	// Get the authorization code from the query parameters
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "Authorization code is required", http.StatusBadRequest)
		return
	}

	// Verify the state parameter to prevent CSRF attacks
	if state != user.StorageAuthorizationCodePkce {
		http.Error(w, "Authorization failed. Please try again.", http.StatusBadRequest)
	}

	// Exchange the authorization code for an access token
	redirectURI := "http://localhost:8080/storage/oauth-callback"
	accessToken, err := user_storage.GetAuthToken(code, redirectURI)
	if err != nil {
		http.Error(w, "Failed to get access token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the token
	user.SetStorageToken(accessToken)

	// Redirect back to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}
