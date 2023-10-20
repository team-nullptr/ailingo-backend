package apiutil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Json is an API response helper for returning json responses.
// If the given payload cannot be marshalled 500 Internal Server Error will be returned.
func Json(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		Err(w, http.StatusInternalServerError, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

// Err is an API response helper for returning error responses.
func Err(w http.ResponseWriter, status int, err error) {
	if err != nil {
		fmt.Println(err)
	}

	body, _ := json.Marshal(map[string]string{
		"error": http.StatusText(status),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}
