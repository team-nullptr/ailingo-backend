package usecase

import (
	"errors"
	"fmt"
)

const StudySetResource = "study_set"
const DefinitionResource = "definition"

// ErrNotFound means that resource was not found.
type ErrNotFound struct {
	Resource string
}

func (err *ErrNotFound) Error() string {
	return fmt.Sprintf("resource not found: [%s]", err.Resource)
}

var (
	// ErrForbidden represents an error meaning that user didn't have
	// enough permissions to perform the operation.
	ErrForbidden = errors.New("unauthorized")

	// ErrRepoFailed means that repository failed to complete a required operation.
	ErrRepoFailed = errors.New("repository failed")

	// ErrValidation means that user submitted data did not meet validation requirements.
	ErrValidation = errors.New("validation error")
)
