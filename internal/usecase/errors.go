package usecase

import "errors"

var (
	ErrForbidden  = errors.New("unauthorized")
	ErrNotFound   = errors.New("resource found")
	ErrRepoFailed = errors.New("repository failed")
	ErrValidation = errors.New("validation error")
)
