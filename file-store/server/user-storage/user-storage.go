package userStorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

// File represents a file in the user storage service
type File struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Created  UnixMillis `json:"created"`
	Updated  UnixMillis `json:"updated"`
	MimeType string     `json:"mimeType"`
}

// Currently use the same URL for files and authorization. May need to be split later.
const dataAuthorizerUrl = "http://localhost:8000"

// Authorize exchanges an OAuth 2.0 authorization code for an access token.
// Should probably live in its own authorization package.
func Authorize(code string, redirectURI string) (string, error) {
	clientID := "YOUR_CLIENT_ID"
	clientSecret := "YOUR_CLIENT_SECRET"
	tokenURL := dataAuthorizerUrl + "/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respData struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	if respData.AccessToken == "" {
		return "", fmt.Errorf("no access token in response")
	}

	return respData.AccessToken, nil
}

func ListFiles(accessToken string) ([]File, error) {
	url := dataAuthorizerUrl + "/files"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var files []File
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return files, nil
}

func GetFileContents(accessToken string, fileId string) (io.Reader, error) {

}

func CreateFile(accessToken string, name string, contents io.Reader) (File, error) {
	// Prepare multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the file field
	fileWriter, err := writer.CreateFormFile("file", name)
	if err != nil {
		return File{}, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(fileWriter, contents); err != nil {
		return File{}, fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Add the name field (if needed by the API)
	_ = writer.WriteField("name", name)

	// Close the writer to finalize the form
	if err := writer.Close(); err != nil {
		return File{}, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Prepare the HTTP request
	url := dataAuthorizerUrl + "/files"
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return File{}, fmt.Errorf("failed to create request: %w", err)
	}
	// Set headers
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return File{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return File{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the response
	var file File
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return File{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return file, nil
}

func DeleteFile(accessToken string, fileId string) error {

}
