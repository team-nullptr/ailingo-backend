package usecase

import "errors"

var (
	// ErrForbidden represents an error meaning that user didn't have
	// enough permissions to perform the operation.
	ErrForbidden = errors.New("unauthorized")

	// ErrNotFound means that resource was not found.
	ErrNotFound = errors.New("resource not found")

	// ErrRepoFailed means that repository failed to complete a required operation.
	ErrRepoFailed = errors.New("repository failed")

	// ErrValidation means that user submitted data did not meet validation requirements.
	ErrValidation = errors.New("validation error")
)
