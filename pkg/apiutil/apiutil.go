package apiutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// Json is an API response helper for returning json responses.
// If the given payload cannot be marshalled 500 Internal Server Error is returned.
func Json(l *slog.Logger, w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		Err(l, w, &ApiError{
			Status: http.StatusInternalServerError,
			Cause:  fmt.Errorf("failed to marshal response body: %w", err),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func Empty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

// Err is an API response helper for returning error responses.
func Err(l *slog.Logger, w http.ResponseWriter, err error) {
	var status int
	var errorMessage string

	var apiError *ApiError
	if errors.As(err, &apiError) {
		if apiError.Cause != nil {
			l.Warn(err.Error())
		}
		status = apiError.Status
		errorMessage = apiError.Message
		if errorMessage == "" {
			errorMessage = http.StatusText(status)
		}
	} else {
		l.Error(err.Error())
		status = http.StatusInternalServerError
		errorMessage = http.StatusText(status)
	}

	body, _ := json.Marshal(map[string]string{
		"error": errorMessage,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}
