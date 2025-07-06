package sessions

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var userSessions = make(map[string]string) // sessionID -> user id

func SetSession(userId string) (*http.Cookie, error) {
	// Generate a simple session ID (not secure, for demo only)
	sessionID := fmt.Sprintf("session-%d", time.Now().UnixNano())
	userSessions[sessionID] = userId
	return &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
		// Secure: true, // Uncomment if using HTTPS
	}, nil
}

func GetUserId(cookie *http.Cookie, err error) (string, error) {
	// Check if session exists in userSessions
	id, ok := userSessions[cookie.Value]
	if !ok {
		return "", errors.New("session not found")
	}

	// Optionally, check if cookie is expired (handled by browser, but double-check for safety)
	if cookie.Expires.Before(time.Now()) {
		delete(userSessions, cookie.Value)
		return "", errors.New("session expired")
	}

	return id, nil
}
