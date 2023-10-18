package apiutil

import (
	"encoding/json"
	"net/http"
)

// Json returns the given payload. If given payload cannot be marshalled returns 500 error message.
func Json(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		Err(w, http.StatusInternalServerError, "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

// Err returns error response with the given message. If the message is not provided default status text will be used.
func Err(w http.ResponseWriter, status int, message string) {
	if message == "" {
		message = http.StatusText(status)
	}

	body, _ := json.Marshal(map[string]string{
		"error": message,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}
