package userStorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type File struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	MimeType  string    `json:"mimeType"`
}


const userStorageServiceUrl = "http://localhost:8000"

func Authorize(username string) (string, error) {

}

func ListFiles(accessToken string) ([]File, error) {

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
	url := userStorageServiceUrl + "/files"
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return File{}, fmt.Errorf("failed to create request: %w", err)
	}
	// Set headers
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{Timeout: 10 * time.Second}
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