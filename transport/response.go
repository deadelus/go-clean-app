package transport

import (
	"encoding/json"
	"net/http"
)

// Response is the standard response format
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendJSON sends a JSON response
func SendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// SendError sends a formatted error response
func SendError(w http.ResponseWriter, status int, message string) {
	SendJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

// SendSuccess sends a formatted success response
func SendSuccess(w http.ResponseWriter, status int, data interface{}) {
	SendJSON(w, status, Response{
		Success: true,
		Data:    data,
	})
}
