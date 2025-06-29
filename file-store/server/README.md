# Go Server

A simple Go HTTP server with test API endpoints.

## Setup

1. Navigate to the server directory:
   ```bash
   cd server
   ```

2. Initialize the Go module and download dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

## API Endpoints

The server will start on `http://localhost:8080` and provides the following endpoints:

### GET /
Welcome endpoint that returns a JSON response.

**Response:**
```json
{
  "message": "Welcome to the Go Server!",
  "timestamp": "2024-01-01T12:00:00Z",
  "status": "success"
}
```

### GET /api/test
Test API endpoint that confirms the server is working.

**Response:**
```json
{
  "message": "This is a test API call! The server is working correctly.",
  "timestamp": "2024-01-01T12:00:00Z",
  "status": "success"
}
```

### GET /api/health
Health check endpoint for monitoring.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "go-server"
}
```

## Testing the API

You can test the endpoints using curl:

```bash
# Test the home endpoint
curl http://localhost:8080/

# Test the API endpoint
curl http://localhost:8080/api/test

# Test the health endpoint
curl http://localhost:8080/api/health
```

Or simply open your browser and navigate to:
- http://localhost:8080/
- http://localhost:8080/api/test
- http://localhost:8080/api/health 