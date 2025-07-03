package sessions

import (
	"errors"
	"net/http"
	"time"
	"fmt"
)

var userSessions = make(map[string]string) // sessionID -> username

func SetSession(username string, password string) (*http.Cookie, error) {
	if username == "" || password == "" {
		return nil, errors.New("Username and password required")
	}
	// Generate a simple session ID (not secure, for demo only)
	sessionID := fmt.Sprintf("session-%d", time.Now().UnixNano())
	userSessions[sessionID] = username
	return &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires: time.Now().Add(24 * time.Hour),
		// Secure: true, // Uncomment if using HTTPS
	}, nil
}

func GetUser(cookie *http.Cookie, err error) (string, error) {
	// Check if session exists in userSessions
	username, ok := userSessions[cookie.Value]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return "", errors.New("session not found");
	}

	// Optionally, check if cookie is expired (handled by browser, but double-check for safety)
	if cookie.Expires.Before(time.Now()) {
		delete(userSessions, cookie.Value)
		http.Redirect(w, r, "/login", http.StatusFound)
		return "", errors.New("session expired")
	}

	return username, nil;
}