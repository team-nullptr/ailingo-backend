package apiutil

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	Status int
	// Message is an end user facing message explaining what went wrong.
	Message string
	Cause   error
}

func (e *ApiError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", http.StatusText(e.Status), e.Cause.Error())
	} else {
		return fmt.Sprintf("%s: %s", http.StatusText(e.Status), e.Message)
	}
}

func (e *ApiError) Unwrap() error {
	return e.Cause
}
